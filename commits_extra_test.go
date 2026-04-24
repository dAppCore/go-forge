package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestCommitService_GetCombinedStatusByRef_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits/main/status" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.CombinedStatus{
			SHA:        "main",
			TotalCount: 3,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	status, err := f.Commits.GetCombinedStatusByRef(context.Background(), "core", "go-forge", "main")
	if err != nil {
		t.Fatal(err)
	}
	if status.SHA != "main" || status.TotalCount != 3 {
		t.Fatalf("got %#v", status)
	}
}

func TestCommitService_ListCommits_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("sha"); got != "main" {
			t.Errorf("got sha=%q, want %q", got, "main")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Errorf("got limit=%q, want %q", got, "25")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Commit{{SHA: "abc123"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	commits, err := f.Commits.ListCommits(context.Background(), "core", "go-forge", &types.ListCommitsOption{
		Sha:      "main",
		Page:     2,
		PageSize: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(commits) != 1 || commits[0].SHA != "abc123" {
		t.Fatalf("got %#v", commits)
	}
}

func TestCommitService_GetCommit_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/git/commits/abc123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Commit{SHA: "abc123"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	commit, err := f.Commits.GetCommit(context.Background(), "core", "go-forge", "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if commit.SHA != "abc123" {
		t.Fatalf("got %#v", commit)
	}
}
