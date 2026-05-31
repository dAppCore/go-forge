package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"dappco.re/go/forge/types"
)

// IterIssues / ListIssues mirror the commits compat shape: a paged branch
// taken when a ListIssueOption carries Page/PageSize/Limit, and an unpaged
// branch delegating to ListIter / ListAll. These tests pin both branches and
// the issueListQuery / issueListPageOptions filter helpers.

func TestIssueService_IterIssues_Paged_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("state"); got != "open" {
			t.Errorf("got state=%q, want open", got)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Issue{{Index: 1}, {Index: 2}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []int64
	for issue, err := range f.Issues.IterIssues(context.Background(), "core", "go-forge", &types.ListIssueOption{State: "open", Limit: 2}) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, issue.Index)
	}
	if len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestIssueService_IterIssues_Paged_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Issues.IterIssues(context.Background(), "core", "go-forge", &types.ListIssueOption{Page: 1}) {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from the paged branch")
	}
}

func TestIssueService_IterIssues_Paged_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode([]types.Issue{{Index: 1}, {Index: 2}, {Index: 3}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Issues.IterIssues(context.Background(), "core", "go-forge", types.ListIssueOption{PageSize: 50}) {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}

func TestIssueService_ListIssues_Paged_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want 2", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Issue{{Index: 7}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListIssues(context.Background(), "core", "go-forge", &types.ListIssueOption{Page: 2, PageSize: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Index != 7 {
		t.Fatalf("got %#v", issues)
	}
}

func TestIssueService_ListIssues_Paged_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"message": "gateway"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if _, err := f.Issues.ListIssues(context.Background(), "core", "go-forge", &types.ListIssueOption{Limit: 5}); err == nil {
		t.Fatal("expected error from paged ListIssues")
	}
}

func TestIssueListQuery_AllShapes_Good(t *testing.T) {
	opt := IssueListOptions{State: "open"}
	optPtr := &IssueListOptions{Sort: "created"}
	compat := types.ListIssueOption{Labels: "bug", CreatedBy: "snider"}
	compatPtr := &types.ListIssueOption{Type: "pulls"}

	q := issueListQuery(opt, optPtr, compat, compatPtr)
	for _, want := range []struct{ k, v string }{
		{"state", "open"}, {"sort", "created"}, {"labels", "bug"}, {"created_by", "snider"}, {"type", "pulls"},
	} {
		if q[want.k] != want.v {
			t.Errorf("key %q = %q, want %q", want.k, q[want.k], want.v)
		}
	}
}

func TestIssueListQuery_CompatAllFields_Good(t *testing.T) {
	// Pin every field of the compat ListIssueOption translation path.
	since := time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)
	before := time.Date(2026, 6, 7, 8, 9, 10, 0, time.UTC)
	compat := types.ListIssueOption{
		State:       "closed",
		Sort:        "updated",
		Labels:      "bug,help",
		Query:       "panic",
		Type:        "issues",
		Milestones:  "beta",
		Since:       &since,
		Before:      &before,
		CreatedBy:   "snider",
		AssignedBy:  "cladius",
		MentionedBy: "virgil",
	}
	q := issueListQuery(compat)
	wantQuery(t, "issueListQueryFromCompat", q, map[string]string{
		"state":        "closed",
		"sort":         "updated",
		"labels":       "bug,help",
		"q":            "panic",
		"type":         "issues",
		"milestones":   "beta",
		"since":        since.Format(time.RFC3339),
		"before":       before.Format(time.RFC3339),
		"created_by":   "snider",
		"assigned_by":  "cladius",
		"mentioned_by": "virgil",
	})
}

func TestIssueListQuery_NilPointers_Bad(t *testing.T) {
	var nilOpt *IssueListOptions
	var nilCompat *types.ListIssueOption
	if q := issueListQuery(nilOpt, nilCompat); q != nil {
		t.Errorf("nil pointers should yield nil query, got %v", q)
	}
}

func TestIssueListPageOptions_Bad(t *testing.T) {
	if _, ok := issueListPageOptions(IssueListOptions{State: "open"}); ok {
		t.Error("non-compat filter should not produce page options")
	}
	if _, ok := issueListPageOptions(types.ListIssueOption{State: "open"}); ok {
		t.Error("compat filter with no paging should not produce page options")
	}
}

func TestIssueListPageOptions_PointerPage_Good(t *testing.T) {
	opts, ok := issueListPageOptions(&types.ListIssueOption{Page: 4})
	if !ok {
		t.Fatal("page on pointer filter should produce page options")
	}
	if opts.Page != 4 {
		t.Fatalf("got %#v", opts)
	}
}

func TestIssueListPageOptions_NilPointer_Bad(t *testing.T) {
	var nilCompat *types.ListIssueOption
	if _, ok := issueListPageOptions(nilCompat); ok {
		t.Error("nil pointer filter should not produce page options")
	}
}
