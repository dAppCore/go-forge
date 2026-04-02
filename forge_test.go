package forge

import (
	"bytes"
	"context"
	"fmt"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestForge_NewForge_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	if f.Repos == nil {
		t.Fatal("Repos service is nil")
	}
	if f.Issues == nil {
		t.Fatal("Issues service is nil")
	}
	if f.ActivityPub == nil {
		t.Fatal("ActivityPub service is nil")
	}
}

func TestForge_Client_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	c := f.Client()
	if c == nil {
		t.Fatal("Client() returned nil")
	}
	if c.baseURL != "https://forge.lthn.ai" {
		t.Errorf("got baseURL=%q", c.baseURL)
	}
	if got := c.BaseURL(); got != "https://forge.lthn.ai" {
		t.Errorf("got BaseURL()=%q", got)
	}
}

func TestForge_BaseURL_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	if got := f.BaseURL(); got != "https://forge.lthn.ai" {
		t.Fatalf("got base URL %q", got)
	}
}

func TestForge_RateLimit_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	if got := f.RateLimit(); got != (RateLimit{}) {
		t.Fatalf("got rate limit %#v", got)
	}
}

func TestForge_UserAgent_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok", WithUserAgent("go-forge/1.0"))
	if got := f.UserAgent(); got != "go-forge/1.0" {
		t.Fatalf("got user agent %q", got)
	}
}

func TestForge_HTTPClient_Good(t *testing.T) {
	custom := &http.Client{}
	f := NewForge("https://forge.lthn.ai", "tok", WithHTTPClient(custom))
	if got := f.HTTPClient(); got != custom {
		t.Fatal("expected HTTPClient() to return the configured HTTP client")
	}
}

func TestForge_HasToken_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok")
	if !f.HasToken() {
		t.Fatal("expected HasToken to report configured token")
	}
}

func TestForge_HasToken_Bad(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "")
	if f.HasToken() {
		t.Fatal("expected HasToken to report missing token")
	}
}

func TestForge_String_Good(t *testing.T) {
	f := NewForge("https://forge.lthn.ai", "tok", WithUserAgent("go-forge/1.0"))
	got := fmt.Sprint(f)
	want := `forge.Forge{forge.Client{baseURL="https://forge.lthn.ai", token=set, userAgent="go-forge/1.0"}}`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := f.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
}

func TestRepoService_ListOrgRepos_Good(t *testing.T) {
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

func TestRepoService_Get_Good(t *testing.T) {
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

func TestRepoService_Update_Good(t *testing.T) {
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

func TestRepoService_Delete_Good(t *testing.T) {
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

func TestRepoService_Get_Bad(t *testing.T) {
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

func TestRepoService_Fork_Good(t *testing.T) {
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

func TestRepoService_GetArchive_Good(t *testing.T) {
	want := []byte("zip-bytes")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/archive/master.zip" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(want)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	got, err := f.Repos.GetArchive(context.Background(), "core", "go-forge", "master.zip")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRepoService_GetRawFile_Good(t *testing.T) {
	want := []byte("# go-forge\n\nA Go client for Forgejo.")
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/raw/README.md" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(want)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	got, err := f.Repos.GetRawFile(context.Background(), "core", "go-forge", "README.md")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestRepoService_ListTags_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/tags" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Tag{{Name: "v1.0.0"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	tags, err := f.Repos.ListTags(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(tags) != 1 || tags[0].Name != "v1.0.0" {
		t.Fatalf("unexpected result: %+v", tags)
	}
}
