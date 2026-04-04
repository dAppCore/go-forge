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
