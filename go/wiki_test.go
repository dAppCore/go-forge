package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestWikiService_ListPages_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/pages" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.WikiPageMetaData{
			{Title: "Home", SubURL: "Home"},
			{Title: "Setup", SubURL: "Setup"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pages, err := f.Wiki.ListPages(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(pages) != 2 {
		t.Fatalf("got %d pages, want 2", len(pages))
	}
	if pages[0].Title != "Home" {
		t.Errorf("got title=%q, want %q", pages[0].Title, "Home")
	}
	if pages[1].Title != "Setup" {
		t.Errorf("got title=%q, want %q", pages[1].Title, "Setup")
	}
}

func TestWikiService_IterPages_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/pages" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.WikiPageMetaData{
			{Title: "Home", SubURL: "Home"},
			{Title: "Setup", SubURL: "Setup"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var titles []string
	for page, err := range f.Wiki.IterPages(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		titles = append(titles, page.Title)
	}
	if len(titles) != 2 {
		t.Fatalf("got %d pages, want 2", len(titles))
	}
	if titles[0] != "Home" || titles[1] != "Setup" {
		t.Fatalf("unexpected titles: %+v", titles)
	}
}

func TestWikiService_GetPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/page/Home" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.WikiPage{
			Title:         "Home",
			ContentBase64: "IyBXZWxjb21l",
			SubURL:        "Home",
			CommitCount:   3,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Wiki.GetPage(context.Background(), "core", "go-forge", "Home")
	if err != nil {
		t.Fatal(err)
	}
	if page.Title != "Home" {
		t.Errorf("got title=%q, want %q", page.Title, "Home")
	}
	if page.ContentBase64 != "IyBXZWxjb21l" {
		t.Errorf("got content=%q, want %q", page.ContentBase64, "IyBXZWxjb21l")
	}
	if page.CommitCount != 3 {
		t.Errorf("got commit_count=%d, want 3", page.CommitCount)
	}
}

func TestWikiService_GetPageRevisions_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/revisions/Home" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "2" {
			t.Errorf("got page=%q, want %q", got, "2")
		}
		json.NewEncoder(w).Encode(types.WikiCommitList{
			Count: 2,
			WikiCommits: []*types.WikiCommit{
				{ID: "abc123", Message: "Initial import"},
				{ID: "def456", Message: "Updated home page"},
			},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	revisions, err := f.Wiki.GetPageRevisions(context.Background(), "core", "go-forge", "Home", 2)
	if err != nil {
		t.Fatal(err)
	}
	if revisions.Count != 2 {
		t.Fatalf("got count=%d, want 2", revisions.Count)
	}
	if len(revisions.WikiCommits) != 2 {
		t.Fatalf("got %d revisions, want 2", len(revisions.WikiCommits))
	}
	if revisions.WikiCommits[0].ID != "abc123" || revisions.WikiCommits[1].Message != "Updated home page" {
		t.Fatalf("got %#v", revisions.WikiCommits)
	}
}

func TestWikiService_CreatePage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/new" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateWikiPageOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Title != "Install" {
			t.Errorf("got title=%q, want %q", opts.Title, "Install")
		}
		if opts.ContentBase64 != "IyBJbnN0YWxs" {
			t.Errorf("got content=%q, want %q", opts.ContentBase64, "IyBJbnN0YWxs")
		}
		json.NewEncoder(w).Encode(types.WikiPage{
			Title:         "Install",
			ContentBase64: "IyBJbnN0YWxs",
			SubURL:        "Install",
			CommitCount:   1,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Wiki.CreatePage(context.Background(), "core", "go-forge", &types.CreateWikiPageOptions{
		Title:         "Install",
		ContentBase64: "IyBJbnN0YWxs",
		Message:       "create install page",
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.Title != "Install" {
		t.Errorf("got title=%q, want %q", page.Title, "Install")
	}
	if page.CommitCount != 1 {
		t.Errorf("got commit_count=%d, want 1", page.CommitCount)
	}
}

func TestWikiService_EditPage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/page/Home" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateWikiPageOptions
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		json.NewEncoder(w).Encode(types.WikiPage{
			Title:         "Home",
			ContentBase64: "dXBkYXRlZA==",
			CommitCount:   4,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Wiki.EditPage(context.Background(), "core", "go-forge", "Home", &types.CreateWikiPageOptions{
		ContentBase64: "dXBkYXRlZA==",
		Message:       "update home page",
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.ContentBase64 != "dXBkYXRlZA==" {
		t.Errorf("got content=%q, want %q", page.ContentBase64, "dXBkYXRlZA==")
	}
}

func TestWikiService_DeletePage_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/wiki/page/Old" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Wiki.DeletePage(context.Background(), "core", "go-forge", "Old")
	if err != nil {
		t.Fatal(err)
	}
}

func TestWikiService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "page not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Wiki.GetPage(context.Background(), "core", "go-forge", "nonexistent")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
