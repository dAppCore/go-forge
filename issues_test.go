package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forge.lthn.ai/core/go-forge/types"
)

func TestIssueService_Good_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Issue{
			{ID: 1, Title: "bug report"},
			{ID: 2, Title: "feature request"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Issues.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].Title != "bug report" {
		t.Errorf("got title=%q, want %q", result.Items[0].Title, "bug report")
	}
}

func TestIssueService_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Issue{ID: 1, Title: "bug report", Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issue, err := f.Issues.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "index": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if issue.Title != "bug report" {
		t.Errorf("got title=%q", issue.Title)
	}
}

func TestIssueService_Good_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body types.CreateIssueOption
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(types.Issue{ID: 1, Title: body.Title, Index: 1})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	issue, err := f.Issues.Create(context.Background(), Params{"owner": "core", "repo": "go-forge"}, &types.CreateIssueOption{
		Title: "new issue",
		Body:  "description here",
	})
	if err != nil {
		t.Fatal(err)
	}
	if issue.Title != "new issue" {
		t.Errorf("got title=%q", issue.Title)
	}
}

func TestIssueService_Good_Pin(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/issues/42/pin" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Issues.Pin(context.Background(), "core", "go-forge", 42)
	if err != nil {
		t.Fatal(err)
	}
}
