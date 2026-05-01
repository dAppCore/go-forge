package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestReleaseService_ListReleases_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Release{{ID: 1, TagName: "v1.0.0"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	releases, err := f.Releases.ListReleases(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(releases) != 1 || releases[0].TagName != "v1.0.0" {
		t.Fatalf("got %#v", releases)
	}
}

func TestReleaseService_CreateRelease_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/releases" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateReleaseOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.TagName != "v1.0.0" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.Release{ID: 1, TagName: body.TagName})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	release, err := f.Releases.CreateRelease(context.Background(), "core", "go-forge", &types.CreateReleaseOption{
		TagName: "v1.0.0",
		Title:   "Release 1.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v1.0.0" {
		t.Fatalf("got tag=%q", release.TagName)
	}
}
