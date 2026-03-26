package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestReleaseService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Release{
			{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"},
			{ID: 2, TagName: "v2.0.0", Title: "Release 2.0"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Releases.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].TagName != "v1.0.0" {
		t.Errorf("got tag=%q, want %q", result.Items[0].TagName, "v1.0.0")
	}
}

func TestReleaseService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Release{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	release, err := f.Releases.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "id": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("got tag=%q, want %q", release.TagName, "v1.0.0")
	}
	if release.Title != "Release 1.0" {
		t.Errorf("got title=%q, want %q", release.Title, "Release 1.0")
	}
}

func TestReleaseService_GetByTag_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases/tags/v1.0.0" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Release{ID: 1, TagName: "v1.0.0", Title: "Release 1.0"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	release, err := f.Releases.GetByTag(context.Background(), "core", "go-forge", "v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v1.0.0" {
		t.Errorf("got tag=%q, want %q", release.TagName, "v1.0.0")
	}
	if release.ID != 1 {
		t.Errorf("got id=%d, want 1", release.ID)
	}
}
