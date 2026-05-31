package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

// IterPullRequests / ListPullRequests follow the compat-filter pattern with a
// paged branch (compat ListPullRequestsOption with Page/PageSize/Limit) and an
// unpaged branch that pages internally via listAll / listIter. These tests pin
// both branches plus the addPullFilters / addCompatPullFilter /
// compatPullListPageOptions helpers.

func TestPullService_IterPullRequests_Paged_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("state"); got != "open" {
			t.Errorf("got state=%q, want open", got)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.PullRequest{{Index: 1}, {Index: 2}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []int64
	for pr, err := range f.Pulls.IterPullRequests(context.Background(), "core", "go-forge", &types.ListPullRequestsOption{State: "open", Limit: 2}) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, pr.Index)
	}
	if len(got) != 2 {
		t.Fatalf("got %v", got)
	}
}

func TestPullService_IterPullRequests_Paged_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Pulls.IterPullRequests(context.Background(), "core", "go-forge", &types.ListPullRequestsOption{Page: 1}) {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error yielded from the paged branch")
	}
}

func TestPullService_IterPullRequests_Paged_Ugly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode([]types.PullRequest{{Index: 1}, {Index: 2}, {Index: 3}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Pulls.IterPullRequests(context.Background(), "core", "go-forge", types.ListPullRequestsOption{PageSize: 50}) {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}

func TestPullService_IterPullRequests_Unpaged_Good(t *testing.T) {
	// No paging fields → the listIter internal-paging branch. A single page
	// with no X-Total-Count and fewer items than the page size terminates.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		json.NewEncoder(w).Encode([]types.PullRequest{{Index: 9}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []int64
	for pr, err := range f.Pulls.IterPullRequests(context.Background(), "core", "go-forge", PullListOptions{State: "closed"}) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, pr.Index)
	}
	if len(got) != 1 || got[0] != 9 {
		t.Fatalf("got %v", got)
	}
}

func TestPullService_ListPullRequests_Unpaged_MultiPage_Good(t *testing.T) {
	// The unpaged listAll branch follows HasMore across pages. Serve a full
	// first page (50 items, total 60) then a short second page.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("X-Total-Count", "60")
		var items []types.PullRequest
		if page == "2" {
			items = make([]types.PullRequest, 10)
		} else {
			items = make([]types.PullRequest, 50)
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	prs, err := f.Pulls.ListPullRequests(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(prs) != 60 {
		t.Fatalf("got %d PRs across pages, want 60", len(prs))
	}
}

func TestPullService_IterPullRequests_Unpaged_MultiPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("X-Total-Count", "60")
		var items []types.PullRequest
		if page == "2" {
			items = make([]types.PullRequest, 10)
		} else {
			items = make([]types.PullRequest, 50)
		}
		json.NewEncoder(w).Encode(items)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for _, err := range f.Pulls.IterPullRequests(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		count++
	}
	if count != 60 {
		t.Fatalf("got %d PRs across pages, want 60", count)
	}
}

func TestPullService_ListPullRequests_Paged_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"message": "gateway"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if _, err := f.Pulls.ListPullRequests(context.Background(), "core", "go-forge", &types.ListPullRequestsOption{Limit: 5}); err == nil {
		t.Fatal("expected error from paged ListPullRequests")
	}
}

func TestPullService_ListPullRequests_Unpaged_Bad(t *testing.T) {
	// listAll surfaces a page error too.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if _, err := f.Pulls.ListPullRequests(context.Background(), "core", "go-forge"); err == nil {
		t.Fatal("expected error from unpaged ListPullRequests")
	}
}

func TestAddPullFilters_AllShapes_Good(t *testing.T) {
	q := newQueryBuilder()
	opt := PullListOptions{State: "open", Labels: []int64{1, 0, 2}}
	optPtr := &PullListOptions{Sort: "created", Milestone: 5}
	compat := types.ListPullRequestsOption{Poster: "snider"}
	compatPtr := &types.ListPullRequestsOption{State: "closed", Sort: "oldest", Labels: []int64{3, 0}, Milestone: 7}

	addPullFilters(q, opt, optPtr, compat, compatPtr)
	enc := q.Encode()
	// Single-valued keys use Set, so the last writer wins (milestone 5 -> 7,
	// state open -> closed); multi-valued labels accumulate via Add.
	for _, want := range []string{"sort=oldest", "milestone=7", "state=closed", "poster=snider", "labels=1", "labels=2", "labels=3"} {
		if !contains(enc, want) {
			t.Errorf("encoded query %q missing %q", enc, want)
		}
	}
	if contains(enc, "milestone=5") {
		t.Errorf("milestone should have been overwritten to 7, got %q", enc)
	}
	// Zero-valued labels must be dropped.
	if contains(enc, "labels=0") {
		t.Errorf("zero label should be dropped, got %q", enc)
	}
}

func TestAddPullFilters_NilPointers_Bad(t *testing.T) {
	q := newQueryBuilder()
	var nilOpt *PullListOptions
	var nilCompat *types.ListPullRequestsOption
	addPullFilters(q, nilOpt, nilCompat)
	if got := q.Encode(); got != "" {
		t.Errorf("nil pointers should add nothing, got %q", got)
	}
}

func TestCompatPullListPageOptions_Bad(t *testing.T) {
	if _, ok := compatPullListPageOptions(PullListOptions{State: "open"}); ok {
		t.Error("non-compat filter should not produce page options")
	}
	if _, ok := compatPullListPageOptions(types.ListPullRequestsOption{State: "open"}); ok {
		t.Error("compat filter with no paging should not produce page options")
	}
	var nilCompat *types.ListPullRequestsOption
	if _, ok := compatPullListPageOptions(nilCompat); ok {
		t.Error("nil pointer should not produce page options")
	}
}

func TestCompatPullListPageOptions_ValuePage_Good(t *testing.T) {
	opts, ok := compatPullListPageOptions(types.ListPullRequestsOption{Page: 2, PageSize: 30})
	if !ok {
		t.Fatal("page on value filter should produce page options")
	}
	if opts.Page != 2 || opts.PageSize != 30 {
		t.Fatalf("got %#v", opts)
	}
}

func contains(haystack, needle string) bool {
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}
