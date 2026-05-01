package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestWebhookService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Hook{
			{ID: 1, Type: "forgejo", Active: true, URL: "https://example.com/hook1"},
			{ID: 2, Type: "forgejo", Active: false, URL: "https://example.com/hook2"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Webhooks.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].URL != "https://example.com/hook1" {
		t.Errorf("got url=%q, want %q", result.Items[0].URL, "https://example.com/hook1")
	}
}

func TestWebhookService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     1,
			Type:   "forgejo",
			Active: true,
			URL:    "https://example.com/hook1",
			Events: []string{"push", "pull_request"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "id": "1"})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 1 {
		t.Errorf("got id=%d, want 1", hook.ID)
	}
	if hook.URL != "https://example.com/hook1" {
		t.Errorf("got url=%q, want %q", hook.URL, "https://example.com/hook1")
	}
	if len(hook.Events) != 2 {
		t.Errorf("got %d events, want 2", len(hook.Events))
	}
}

func TestWebhookService_Create_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var opts types.CreateHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Type != "forgejo" {
			t.Errorf("got type=%q, want %q", opts.Type, "forgejo")
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     3,
			Type:   opts.Type,
			Active: opts.Active,
			Events: opts.Events,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.Create(context.Background(), Params{"owner": "core", "repo": "go-forge"}, &types.CreateHookOption{
		Type:   "forgejo",
		Active: true,
		Events: []string{"push"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 3 {
		t.Errorf("got id=%d, want 3", hook.ID)
	}
	if hook.Type != "forgejo" {
		t.Errorf("got type=%q, want %q", hook.Type, "forgejo")
	}
}

func TestWebhookService_TestHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/1/tests" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Webhooks.TestHook(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebhookService_ListGitHooks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/git" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.GitHook{
			{Name: "pre-receive", Content: "#!/bin/sh\nexit 0", IsActive: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hooks, err := f.Webhooks.ListGitHooks(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 1 {
		t.Fatalf("got %d hooks, want 1", len(hooks))
	}
	if hooks[0].Name != "pre-receive" {
		t.Errorf("got name=%q, want %q", hooks[0].Name, "pre-receive")
	}
}

func TestWebhookService_IterGitHooks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/git" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.GitHook{
			{Name: "pre-receive", Content: "#!/bin/sh\nexit 0", IsActive: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for hook, err := range f.Webhooks.IterGitHooks(context.Background(), "core", "go-forge") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, hook.Name)
	}
	if len(got) != 1 {
		t.Fatalf("got %d hooks, want 1", len(got))
	}
	if got[0] != "pre-receive" {
		t.Errorf("got name=%q, want %q", got[0], "pre-receive")
	}
}

func TestWebhookService_GetGitHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/git/pre-receive" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.GitHook{
			Name:     "pre-receive",
			Content:  "#!/bin/sh\nexit 0",
			IsActive: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.GetGitHook(context.Background(), "core", "go-forge", "pre-receive")
	if err != nil {
		t.Fatal(err)
	}
	if hook.Name != "pre-receive" {
		t.Errorf("got name=%q, want %q", hook.Name, "pre-receive")
	}
	if !hook.IsActive {
		t.Error("expected is_active=true")
	}
}

func TestWebhookService_EditGitHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/git/pre-receive" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditGitHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Content != "#!/bin/sh\nexit 0" {
			t.Fatalf("unexpected edit payload: %+v", opts)
		}
		json.NewEncoder(w).Encode(types.GitHook{
			Name:     "pre-receive",
			Content:  opts.Content,
			IsActive: true,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.EditGitHook(context.Background(), "core", "go-forge", "pre-receive", &types.EditGitHookOption{
		Content: "#!/bin/sh\nexit 0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.Content != "#!/bin/sh\nexit 0" {
		t.Errorf("got content=%q, want %q", hook.Content, "#!/bin/sh\nexit 0")
	}
}

func TestWebhookService_DeleteGitHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks/git/pre-receive" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Webhooks.DeleteGitHook(context.Background(), "core", "go-forge", "pre-receive"); err != nil {
		t.Fatal(err)
	}
}

func TestWebhookService_ListUserHooks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Hook{
			{ID: 20, Type: "forgejo", Active: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hooks, err := f.Webhooks.ListUserHooks(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 1 {
		t.Fatalf("got %d hooks, want 1", len(hooks))
	}
	if hooks[0].ID != 20 {
		t.Errorf("got id=%d, want 20", hooks[0].ID)
	}
}

func TestWebhookService_GetUserHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/hooks/20" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     20,
			Type:   "forgejo",
			Active: true,
			URL:    "https://example.com/user-hook",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.GetUserHook(context.Background(), 20)
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 20 {
		t.Errorf("got id=%d, want 20", hook.ID)
	}
	if hook.URL != "https://example.com/user-hook" {
		t.Errorf("got url=%q, want %q", hook.URL, "https://example.com/user-hook")
	}
}

func TestWebhookService_CreateUserHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Type != "forgejo" {
			t.Errorf("got type=%q, want %q", opts.Type, "forgejo")
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     21,
			Type:   opts.Type,
			Active: opts.Active,
			Events: opts.Events,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.CreateUserHook(context.Background(), &types.CreateHookOption{
		Type:   "forgejo",
		Active: true,
		Events: []string{"push"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 21 {
		t.Errorf("got id=%d, want 21", hook.ID)
	}
	if hook.Type != "forgejo" {
		t.Errorf("got type=%q, want %q", hook.Type, "forgejo")
	}
}

func TestWebhookService_EditUserHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/hooks/20" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Active != false {
			t.Fatalf("unexpected edit payload: %+v", opts)
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     20,
			Type:   "forgejo",
			Active: false,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.EditUserHook(context.Background(), 20, &types.EditHookOption{
		Active: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 20 {
		t.Errorf("got id=%d, want 20", hook.ID)
	}
	if hook.Active {
		t.Error("expected active=false")
	}
}

func TestWebhookService_DeleteUserHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user/hooks/20" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Webhooks.DeleteUserHook(context.Background(), 20); err != nil {
		t.Fatal(err)
	}
}

func TestWebhookService_ListOrgHooks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Hook{
			{ID: 10, Type: "forgejo", Active: true},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hooks, err := f.Webhooks.ListOrgHooks(context.Background(), "myorg")
	if err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 1 {
		t.Errorf("got %d hooks, want 1", len(hooks))
	}
	if hooks[0].ID != 10 {
		t.Errorf("got id=%d, want 10", hooks[0].ID)
	}
}

func TestWebhookService_GetOrgHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/hooks/10" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     10,
			Type:   "forgejo",
			Active: true,
			URL:    "https://example.com/org-hook",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.GetOrgHook(context.Background(), "myorg", 10)
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 10 {
		t.Errorf("got id=%d, want 10", hook.ID)
	}
	if hook.URL != "https://example.com/org-hook" {
		t.Errorf("got url=%q, want %q", hook.URL, "https://example.com/org-hook")
	}
}

func TestWebhookService_CreateOrgHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Type != "forgejo" {
			t.Errorf("got type=%q, want %q", opts.Type, "forgejo")
		}
		json.NewEncoder(w).Encode(types.Hook{
			ID:     11,
			Type:   opts.Type,
			Active: opts.Active,
			Events: opts.Events,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.CreateOrgHook(context.Background(), "myorg", &types.CreateHookOption{
		Type:   "forgejo",
		Active: true,
		Events: []string{"push"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 11 {
		t.Errorf("got id=%d, want 11", hook.ID)
	}
	if hook.Type != "forgejo" {
		t.Errorf("got type=%q, want %q", hook.Type, "forgejo")
	}
}

func TestWebhookService_EditOrgHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/hooks/10" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditHookOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Active != false {
			t.Fatalf("unexpected edit payload: %+v", opts)
		}
		active := false
		json.NewEncoder(w).Encode(types.Hook{
			ID:     10,
			Type:   "forgejo",
			Active: active,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.EditOrgHook(context.Background(), "myorg", 10, &types.EditHookOption{
		Active: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.ID != 10 {
		t.Errorf("got id=%d, want 10", hook.ID)
	}
	if hook.Active {
		t.Error("expected active=false")
	}
}

func TestWebhookService_DeleteOrgHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/hooks/10" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Webhooks.DeleteOrgHook(context.Background(), "myorg", 10); err != nil {
		t.Fatal(err)
	}
}

func TestWebhookService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "hook not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Webhooks.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "id": "999"})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
