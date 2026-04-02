package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestWebhookService_ListHooks_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Hook{{ID: 1, Type: "forgejo"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hooks, err := f.Webhooks.ListHooks(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(hooks) != 1 || hooks[0].ID != 1 {
		t.Fatalf("got %#v", hooks)
	}
}

func TestWebhookService_CreateHook_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/hooks" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateHookOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.Type != "forgejo" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.Hook{ID: 1, Type: body.Type})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	hook, err := f.Webhooks.CreateHook(context.Background(), "core", "go-forge", &types.CreateHookOption{
		Type: "forgejo",
		Config: &types.CreateHookOptionConfig{
			"content_type": "json",
			"url":          "https://example.com/hook",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	if hook.Type != "forgejo" {
		t.Fatalf("got type=%q", hook.Type)
	}
}
