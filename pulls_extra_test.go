package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestPullService_ListPullRequests_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PullRequest{{ID: 1, Title: "add feature"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	prs, err := f.Pulls.ListPullRequests(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(prs) != 1 || prs[0].Title != "add feature" {
		t.Fatalf("got %#v", prs)
	}
}

func TestPullService_CreatePullRequest_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreatePullRequestOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Title != "add feature" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.PullRequest{ID: 1, Title: body.Title, Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pr, err := f.Pulls.CreatePullRequest(context.Background(), "core", "go-forge", &types.CreatePullRequestOption{
		Title: "add feature",
		Base:  "main",
		Head:  "feature",
	})
	if err != nil {
		t.Fatal(err)
	}
	if pr.Title != "add feature" {
		t.Fatalf("got title=%q", pr.Title)
	}
}

func TestPullService_ListPullRequests_CompatPagination_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/pulls" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("state"); got != "open" {
			t.Errorf("got state=%q, want %q", got, "open")
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		if got := r.URL.Query().Get("limit"); got != "25" {
			t.Errorf("got limit=%q, want %q", got, "25")
		}
		w.Header().Set("X-Total-Count", "40")
		json.NewEncoder(w).Encode([]types.PullRequest{{ID: 2, Title: "paged pull"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	prs, err := f.Pulls.ListPullRequests(context.Background(), "core", "go-forge", &types.ListPullRequestsOption{
		State:    "open",
		Page:     2,
		PageSize: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(prs) != 1 || prs[0].Title != "paged pull" {
		t.Fatalf("got %#v", prs)
	}
}
