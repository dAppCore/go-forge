package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestCommitService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("got page=%q, want %q", r.URL.Query().Get("page"), "1")
		}
		if r.URL.Query().Get("limit") != "50" {
			t.Errorf("got limit=%q, want %q", r.URL.Query().Get("limit"), "50")
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Commit{
			{
				SHA: "abc123",
				Commit: &types.RepoCommit{
					Message: "first commit",
				},
			},
			{
				SHA: "def456",
				Commit: &types.RepoCommit{
					Message: "second commit",
				},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Commits.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Fatalf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].SHA != "abc123" {
		t.Errorf("got sha=%q, want %q", result.Items[0].SHA, "abc123")
	}
	if result.Items[1].Commit == nil {
		t.Fatal("expected commit payload, got nil")
	}
	if result.Items[1].Commit.Message != "second commit" {
		t.Errorf("got message=%q, want %q", result.Items[1].Commit.Message, "second commit")
	}
}

func TestCommitService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/git/commits/abc123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Commit{
			SHA:     "abc123",
			HTMLURL: "https://forge.example/core/go-forge/commit/abc123",
			Commit: &types.RepoCommit{
				Message: "initial import",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	commit, err := f.Commits.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "sha": "abc123"})
	if err != nil {
		t.Fatal(err)
	}
	if commit.SHA != "abc123" {
		t.Errorf("got sha=%q, want %q", commit.SHA, "abc123")
	}
	if commit.Commit == nil {
		t.Fatal("expected commit payload, got nil")
	}
	if commit.Commit.Message != "initial import" {
		t.Errorf("got message=%q, want %q", commit.Commit.Message, "initial import")
	}
}

func TestCommitService_GetDiffOrPatch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/git/commits/abc123.diff" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "text/plain")
		io.WriteString(w, "diff --git a/README.md b/README.md")
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	data, err := f.Commits.GetDiffOrPatch(context.Background(), "core", "go-forge", "abc123", "diff")
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "diff --git a/README.md b/README.md" {
		t.Fatalf("got body=%q", string(data))
	}
}

func TestCommitService_GetPullRequest_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits/abc123/pull" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.PullRequest{
			ID:    17,
			Index: 9,
			Title: "Add commit-linked pull request",
			Head: &types.PRBranchInfo{
				Ref: "feature/commit-link",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Commits.GetPullRequest(context.Background(), "core", "go-forge", "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if pr.ID != 17 {
		t.Errorf("got id=%d, want 17", pr.ID)
	}
	if pr.Index != 9 {
		t.Errorf("got index=%d, want 9", pr.Index)
	}
	if pr.Head == nil || pr.Head.Ref != "feature/commit-link" {
		t.Fatalf("unexpected head branch info: %+v", pr.Head)
	}
}

func TestCommitService_ListStatuses_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits/abc123/statuses" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.CommitStatus{
			{ID: 1, Context: "ci/build", Description: "Build passed"},
			{ID: 2, Context: "ci/test", Description: "Tests passed"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	statuses, err := f.Commits.ListStatuses(context.Background(), "core", "go-forge", "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if len(statuses) != 2 {
		t.Fatalf("got %d statuses, want 2", len(statuses))
	}
	if statuses[0].Context != "ci/build" {
		t.Errorf("got context=%q, want %q", statuses[0].Context, "ci/build")
	}
	if statuses[1].Context != "ci/test" {
		t.Errorf("got context=%q, want %q", statuses[1].Context, "ci/test")
	}
}

func TestCommitService_IterStatuses_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/commits/abc123/statuses" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.CommitStatus{
			{ID: 1, Context: "ci/build"},
			{ID: 2, Context: "ci/test"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for status, err := range f.Commits.IterStatuses(context.Background(), "core", "go-forge", "abc123") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, status.Context)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if len(got) != 2 || got[0] != "ci/build" || got[1] != "ci/test" {
		t.Fatalf("got %#v", got)
	}
}

func TestCommitService_CreateStatus_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/statuses/abc123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateStatusOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Context != "ci/build" {
			t.Errorf("got context=%q, want %q", opts.Context, "ci/build")
		}
		if opts.Description != "Build passed" {
			t.Errorf("got description=%q, want %q", opts.Description, "Build passed")
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.CommitStatus{
			ID:          1,
			Context:     "ci/build",
			Description: "Build passed",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	status, err := f.Commits.CreateStatus(context.Background(), "core", "go-forge", "abc123", &types.CreateStatusOption{
		Context:     "ci/build",
		Description: "Build passed",
	})
	if err != nil {
		t.Fatal(err)
	}
	if status.Context != "ci/build" {
		t.Errorf("got context=%q, want %q", status.Context, "ci/build")
	}
	if status.ID != 1 {
		t.Errorf("got id=%d, want 1", status.ID)
	}
}

func TestCommitService_GetNote_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/git/notes/abc123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Note{
			Message: "reviewed and approved",
			Commit: &types.Commit{
				SHA: "abc123",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	note, err := f.Commits.GetNote(context.Background(), "core", "go-forge", "abc123")
	if err != nil {
		t.Fatal(err)
	}
	if note.Message != "reviewed and approved" {
		t.Errorf("got message=%q, want %q", note.Message, "reviewed and approved")
	}
	if note.Commit.SHA != "abc123" {
		t.Errorf("got commit sha=%q, want %q", note.Commit.SHA, "abc123")
	}
}

func TestCommitService_SetNote_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/git/notes/abc123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.NoteOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Message != "reviewed and approved" {
			t.Errorf("got message=%q, want %q", opts.Message, "reviewed and approved")
		}
		json.NewEncoder(w).Encode(types.Note{
			Message: "reviewed and approved",
			Commit: &types.Commit{
				SHA: "abc123",
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	note, err := f.Commits.SetNote(context.Background(), "core", "go-forge", "abc123", "reviewed and approved")
	if err != nil {
		t.Fatal(err)
	}
	if note.Message != "reviewed and approved" {
		t.Errorf("got message=%q, want %q", note.Message, "reviewed and approved")
	}
	if note.Commit.SHA != "abc123" {
		t.Errorf("got commit sha=%q, want %q", note.Commit.SHA, "abc123")
	}
}

func TestCommitService_DeleteNote_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/git/notes/abc123" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Commits.DeleteNote(context.Background(), "core", "go-forge", "abc123"); err != nil {
		t.Fatal(err)
	}
}

func TestCommitService_GetCombinedStatus_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/statuses/main" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.CombinedStatus{
			SHA:        "abc123",
			TotalCount: 2,
			Statuses: []*types.CommitStatus{
				{ID: 1, Context: "ci/build"},
				{ID: 2, Context: "ci/test"},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	cs, err := f.Commits.GetCombinedStatus(context.Background(), "core", "go-forge", "main")
	if err != nil {
		t.Fatal(err)
	}
	if cs.SHA != "abc123" {
		t.Errorf("got sha=%q, want %q", cs.SHA, "abc123")
	}
	if cs.TotalCount != 2 {
		t.Errorf("got total_count=%d, want 2", cs.TotalCount)
	}
	if len(cs.Statuses) != 2 {
		t.Fatalf("got %d statuses, want 2", len(cs.Statuses))
	}
}

func TestCommitService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Commits.GetNote(context.Background(), "core", "go-forge", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
