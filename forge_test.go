package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forge.lthn.ai/core/go-forge/types"
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

func TestRepoService_Good_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Repository{{Name: "go-forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.List(context.Background(), Params{"org": "core"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 1 || result.Items[0].Name != "go-forge" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestRepoService_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
