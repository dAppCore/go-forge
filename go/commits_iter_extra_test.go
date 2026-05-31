package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

// IterCommits has two distinct branches: a paged branch (taken when a
// compat ListCommitsOption carries Page/PageSize/Limit) that materialises a
// single page and yields its items, and the unpaged branch that delegates to
// ListIter. These tests exercise the paged branch's happy path, its error
// yield, and its early-termination, plus the unpaged branch.

func TestCommitService_IterCommits_Paged_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want 1", got)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Commit{{SHA: "a"}, {SHA: "b"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for commit, err := range f.Commits.IterCommits(context.Background(), "core", "go-forge", &types.ListCommitsOption{Limit: 2}) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, commit.SHA)
	}
	if len(got) != 2 || got[0] != "a" || got[1] != "b" {
		t.Fatalf("got %v", got)
	}
}

func TestCommitService_IterCommits_Paged_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var sawErr bool
	for _, err := range f.Commits.IterCommits(context.Background(), "core", "go-forge", &types.ListCommitsOption{Page: 1}) {
		if err != nil {
			sawErr = true
			break
		}
	}
	if !sawErr {
		t.Fatal("expected an error to be yielded from the paged branch")
	}
}

func TestCommitService_IterCommits_Paged_Ugly(t *testing.T) {
	// Early break after the first item must stop the iteration; the paged
	// branch honours the yield return value.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "3")
		json.NewEncoder(w).Encode([]types.Commit{{SHA: "x"}, {SHA: "y"}, {SHA: "z"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var count int
	for range f.Commits.IterCommits(context.Background(), "core", "go-forge", &types.ListCommitsOption{PageSize: 50}) {
		count++
		break
	}
	if count != 1 {
		t.Fatalf("early break should stop after one item, got %d", count)
	}
}

func TestCommitService_IterCommits_Unpaged_Good(t *testing.T) {
	// No paging fields on the filter selects the ListIter delegation branch.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Commit{{SHA: "solo"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for commit, err := range f.Commits.IterCommits(context.Background(), "core", "go-forge", CommitListOptions{Sha: "main"}) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, commit.SHA)
	}
	if len(got) != 1 || got[0] != "solo" {
		t.Fatalf("got %v", got)
	}
}

func TestCommitService_ListCommits_Paged_Bad(t *testing.T) {
	// The paged branch of ListCommits propagates the page error instead of
	// silently returning an empty slice.
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
		json.NewEncoder(w).Encode(map[string]string{"message": "gateway"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Commits.ListCommits(context.Background(), "core", "go-forge", &types.ListCommitsOption{Page: 1})
	if err == nil {
		t.Fatal("expected error from paged ListCommits")
	}
}

func TestCommitService_ListCommits_Unpaged_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Commit{{SHA: "nopage"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	commits, err := f.Commits.ListCommits(context.Background(), "core", "go-forge", CommitListOptions{Path: "go/"})
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 1 || commits[0].SHA != "nopage" {
		t.Fatalf("got %#v", commits)
	}
}

// commitCompatListQuery merges several filter shapes. Cover the typed-value,
// typed-pointer and nil-pointer cases explicitly.
func TestCommitCompatListQuery_AllShapes_Good(t *testing.T) {
	yes := true
	opts := CommitListOptions{Sha: "v"}
	ptr := &CommitListOptions{Path: "p"}
	no := false
	compat := types.ListCommitsOption{Sha: "cs", Path: "cp", Not: "main", Stat: &yes, Verification: &no}
	compatPtr := &types.ListCommitsOption{Files: &yes}

	q := commitCompatListQuery(opts, ptr, compat, compatPtr)
	for _, want := range []struct{ k, v string }{
		{"sha", "cs"}, {"path", "cp"}, {"not", "main"}, {"stat", "true"}, {"verification", "false"}, {"files", "true"},
	} {
		if q[want.k] != want.v {
			t.Errorf("key %q = %q, want %q", want.k, q[want.k], want.v)
		}
	}
}

func TestCommitCompatListQuery_NilPointers_Bad(t *testing.T) {
	// Nil typed-pointers must be skipped, and an all-empty filter set yields nil.
	var nilOpts *CommitListOptions
	var nilCompat *types.ListCommitsOption
	if q := commitCompatListQuery(nilOpts, nilCompat); q != nil {
		t.Errorf("nil pointers should yield nil query, got %v", q)
	}
}

func TestCommitCompatPageOptions_Bad(t *testing.T) {
	// No paging fields anywhere → ok=false.
	if _, ok := commitCompatPageOptions(CommitListOptions{Sha: "x"}); ok {
		t.Error("non-paging filter should not produce page options")
	}
	if _, ok := commitCompatPageOptions(); ok {
		t.Error("no filters should not produce page options")
	}
}

func TestCommitCompatPageOptions_PointerLimit_Good(t *testing.T) {
	opts, ok := commitCompatPageOptions(&types.ListCommitsOption{Limit: 10})
	if !ok {
		t.Fatal("limit on pointer filter should produce page options")
	}
	if opts.Limit != 10 || opts.Page != 1 {
		t.Fatalf("got %#v", opts)
	}
}

func TestCommitCompatPageOptions_ValuePage_Good(t *testing.T) {
	// Typed-value (non-pointer) compat filter carrying a page is the other
	// branch of the type switch.
	opts, ok := commitCompatPageOptions(types.ListCommitsOption{Page: 3, PageSize: 20})
	if !ok {
		t.Fatal("page on value filter should produce page options")
	}
	if opts.Page != 3 || opts.PageSize != 20 {
		t.Fatalf("got %#v", opts)
	}
}
