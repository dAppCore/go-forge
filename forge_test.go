package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestForge_Good_NewForge(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	if f.Repos == nil {
		t.Fatal("Repos service is nil")
	}
	if f.Issues == nil {
		t.Fatal("Issues service is nil")
	}
}

func TestForge_Good_Client(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	c := f.Client()
	if c == nil {
		t.Fatal("Client() returned nil")
	}
	if c.baseURL != "https://forge.lthn.ai" {
		t.Errorf("got baseURL=%q", c.baseURL)
	}
}

func TestRepoService_Good_ListOrgRepos(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{{Name: "go-forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repos, err := f.Repos.ListOrgRepos(context.Background(), "core")
	if err != nil {
		t.Fatal(err)
	}
	if len(repos) != 1 || repos[0].Name != "go-forge" {
		t.Errorf("unexpected result: %+v", repos)
	}
}

func TestRepoService_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/repos/core/go-forge" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.Repository{Name: "go-forge", FullName: "core/go-forge"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repo, err := f.Repos.Get(context.Background(), Params{"owner": "core", "repo": "go-forge"})
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "go-forge" {
		t.Errorf("got name=%q", repo.Name)
	}
}

func TestRepoService_Good_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.EditRepoOption
		json.NewDecoder(r.Body).Decode(&body)
		json.NewEncoder(w).Encode(types.Repository{Name: body.Name, FullName: "core/" + body.Name})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repo, err := f.Repos.Update(context.Background(), Params{"owner": "core", "repo": "go-forge"}, &types.EditRepoOption{
		Name: "go-forge-renamed",
	})
	if err != nil {
		t.Fatal(err)
	}
	if repo.Name != "go-forge-renamed" {
		t.Errorf("got name=%q", repo.Name)
	}
}

func TestRepoService_Good_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.Delete(context.Background(), Params{"owner": "core", "repo": "go-forge"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_Bad_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if _, err := f.Repos.Get(context.Background(), Params{"owner": "core", "repo": "go-forge"}); !IsNotFound(err) {
		t.Fatalf("expected not found, got %v", err)
	}
}

func TestRepoService_Good_Fork(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(types.Repository{Name: "go-forge", Fork: true})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	repo, err := f.Repos.Fork(context.Background(), "core", "go-forge", "my-org")
	if err != nil {
		t.Fatal(err)
	}
	if !repo.Fork {
		t.Error("expected fork=true")
	}
}
