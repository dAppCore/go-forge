package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestWebhookService_Good_List(t *testing.T) {
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

func TestWebhookService_Good_Get(t *testing.T) {
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

func TestWebhookService_Good_Create(t *testing.T) {
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

func TestWebhookService_Good_TestHook(t *testing.T) {
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

func TestWebhookService_Good_ListOrgHooks(t *testing.T) {
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

func TestWebhookService_Bad_NotFound(t *testing.T) {
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
