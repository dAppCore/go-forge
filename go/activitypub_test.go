package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestActivityPubService_GetInstanceActor_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/activitypub/actor" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(types.ActivityPub{Context: "https://www.w3.org/ns/activitystreams"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	actor, err := f.ActivityPub.GetInstanceActor(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if actor.Context != "https://www.w3.org/ns/activitystreams" {
		t.Fatalf("got context=%q", actor.Context)
	}
}

func TestActivityPubService_SendRepositoryInbox_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/activitypub/repository-id/42/inbox" {
			t.Errorf("wrong path: %s", r.URL.Path)
			http.NotFound(w, r)
			return
		}
		var body types.ForgeLike
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.ActivityPub.SendRepositoryInbox(context.Background(), 42, &types.ForgeLike{}); err != nil {
		t.Fatal(err)
	}
}
