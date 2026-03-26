package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestPullService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.PullRequest{
			{ID: 1, Title: "add feature"},
			{ID: 2, Title: "fix bug"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Pulls.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].Title != "add feature" {
		t.Errorf("got title=%q, want %q", result.Items[0].Title, "add feature")
	}
}

func TestPullService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.PullRequest{ID: 1, Title: "add feature", Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Pulls.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Title != "add feature" {
		t.Errorf("got title=%q", pr.Title)
	}
}

func TestPullService_Create_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.CreatePullRequestOption
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.PullRequest{ID: 1, Title: body.Title, Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Pulls.Create(context.Background(), Params{"owner": "core", "repo": "go-forge"}, &types.CreatePullRequestOption{
		Title: "new pull request",
		Head:  "feature-branch",
		Base:  "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Title != "new pull request" {
		t.Errorf("got title=%q", pr.Title)
	}
}

func TestPullService_Merge_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls/7/merge" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		if body["Do"] != "merge" {
			t.Errorf("got Do=%q, want %q", body["Do"], "merge")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Pulls.Merge(context.Background(), "core", "go-forge", 7, "merge")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPullService_Merge_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]string{"message": "already merged"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Pulls.Merge(context.Background(), "core", "go-forge", 7, "merge"); !IsConflict(err) {
		t.Fatalf("expected conflict, got %v", err)
	}
}
