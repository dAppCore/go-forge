package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestIssueService_ListIssues_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Issue{{ID: 1, Title: "bug"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListIssues(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Title != "bug" {
		t.Fatalf("got %#v", issues)
	}
}

func TestIssueService_CreateIssue_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateIssueOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Title != "new issue" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.Issue{ID: 1, Index: 1, Title: body.Title})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issue, err := f.Issues.CreateIssue(context.Background(), "core", "go-forge", &types.CreateIssueOption{Title: "new issue"})
	if err != nil {
		t.Fatal(err)
	}
	if issue.Title != "new issue" {
		t.Fatalf("got title=%q", issue.Title)
	}
}

func TestIssueService_ListRepoIssues_CompatPagination_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues" {
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
		json.NewEncoder(w).Encode([]types.Issue{{ID: 2, Title: "paged issue"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issues, err := f.Issues.ListRepoIssues(context.Background(), "core", "go-forge", &types.ListIssueOption{
		State:    "open",
		Page:     2,
		PageSize: 25,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(issues) != 1 || issues[0].Title != "paged issue" {
		t.Fatalf("got %#v", issues)
	}
}
