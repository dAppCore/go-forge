package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestRepoService_ListTopics_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.TopicName{TopicNames: []string{"go", "forge"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	topics, err := f.Repos.ListTopics(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(topics, []string{"go", "forge"}) {
		t.Fatalf("got %#v", topics)
	}
}

func TestRepoService_UpdateTopics_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/topics" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.RepoTopicOptions
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode body: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !reflect.DeepEqual(body.Topics, []string{"go", "forge"}) {
			t.Fatalf("got %#v", body.Topics)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.UpdateTopics(context.Background(), "core", "go-forge", []string{"go", "forge"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_AddTopic_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.EscapedPath() != "/api/v1/repos/core/go-forge/topics/release%20candidate" {
			t.Errorf("wrong path: %s", r.URL.EscapedPath())
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.AddTopic(context.Background(), "core", "go-forge", "release candidate"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteTopic_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.EscapedPath() != "/api/v1/repos/core/go-forge/topics/release%20candidate" {
			t.Errorf("wrong path: %s", r.URL.EscapedPath())
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteTopic(context.Background(), "core", "go-forge", "release candidate"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_GetNewPinAllowed_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/new_pin_allowed" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.NewIssuePinsAllowed{
			Issues:       true,
			PullRequests: false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.GetNewPinAllowed(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Issues || result.PullRequests {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_GetSubscription_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscription" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.WatchInfo{
			Subscribed: true,
			Ignored:    false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.GetSubscription(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Subscribed || result.Ignored {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_ListStargazers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/stargazers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}, {UserName: "bob"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListStargazers(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 2 || users[0].UserName != "alice" || users[1].UserName != "bob" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_ListSubscribers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscribers" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "charlie"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListSubscribers(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].UserName != "charlie" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_ListCollaborators_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode([]types.User{{UserName: "alice"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	users, err := f.Repos.ListCollaborators(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 1 || users[0].UserName != "alice" {
		t.Fatalf("got %#v", users)
	}
}

func TestRepoService_AddCollaborator_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.AddCollaboratorOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		if body.Permission != "write" {
			t.Fatalf("got permission=%q, want %q", body.Permission, "write")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.AddCollaborator(context.Background(), "core", "go-forge", "alice", &types.AddCollaboratorOption{Permission: "write"}); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_DeleteCollaborator_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.DeleteCollaborator(context.Background(), "core", "go-forge", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_CheckCollaborator_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	ok, err := f.Repos.CheckCollaborator(context.Background(), "core", "go-forge", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected collaborator check to return true")
	}
}

func TestRepoService_GetCollaboratorPermission_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/collaborators/alice/permission" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.RepoCollaboratorPermission{
			Permission: "write",
			RoleName:   "collaborator",
			User:       &types.User{UserName: "alice"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	perm, err := f.Repos.GetCollaboratorPermission(context.Background(), "core", "go-forge", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if perm.Permission != "write" || perm.User == nil || perm.User.UserName != "alice" {
		t.Fatalf("got %#v", perm)
	}
}

func TestRepoService_Watch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscription" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.WatchInfo{
			Subscribed: true,
			Ignored:    false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Repos.Watch(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Subscribed || result.Ignored {
		t.Fatalf("got %#v", result)
	}
}

func TestRepoService_Unwatch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/subscription" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Repos.Unwatch(context.Background(), "core", "go-forge"); err != nil {
		t.Fatal(err)
	}
}

func TestRepoService_PathParamsAreEscaped_Good(t *testing.T) {
	owner := "acme org"
	repo := "my/repo"
	org := "team alpha"

	t.Run("ListOrgRepos", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.EscapedPath() != "/api/v1/orgs/team%20alpha/repos" {
				t.Errorf("got path %q, want %q", r.URL.EscapedPath(), "/api/v1/orgs/team%20alpha/repos")
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode([]types.Repository{})
		}))
		defer srv.Close()

		f := NewForge(srv.URL, "tok")
		_, err := f.Repos.ListOrgRepos(context.Background(), org)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Fork", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			want := "/api/v1/repos/acme%20org/my%2Frepo/forks"
			if r.URL.EscapedPath() != want {
				t.Errorf("got path %q, want %q", r.URL.EscapedPath(), want)
				http.NotFound(w, r)
				return
			}
			json.NewEncoder(w).Encode(types.Repository{Name: repo})
		}))
		defer srv.Close()

		f := NewForge(srv.URL, "tok")
		if _, err := f.Repos.Fork(context.Background(), owner, repo, ""); err != nil {
			t.Fatal(err)
		}
	})
}
