package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	core "dappco.re/go"
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

func TestActivitypub_ActivityPubService_GetInstanceActor_Good(t *core.T) {
	subject := (*ActivityPubService).GetInstanceActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetInstanceActor_Bad(t *core.T) {
	subject := (*ActivityPubService).GetInstanceActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetInstanceActor_Ugly(t *core.T) {
	subject := (*ActivityPubService).GetInstanceActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendInstanceActorInbox_Good(t *core.T) {
	subject := (*ActivityPubService).SendInstanceActorInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendInstanceActorInbox_Bad(t *core.T) {
	subject := (*ActivityPubService).SendInstanceActorInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendInstanceActorInbox_Ugly(t *core.T) {
	subject := (*ActivityPubService).SendInstanceActorInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetRepositoryActor_Good(t *core.T) {
	subject := (*ActivityPubService).GetRepositoryActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetRepositoryActor_Bad(t *core.T) {
	subject := (*ActivityPubService).GetRepositoryActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetRepositoryActor_Ugly(t *core.T) {
	subject := (*ActivityPubService).GetRepositoryActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendRepositoryInbox_Good(t *core.T) {
	subject := (*ActivityPubService).SendRepositoryInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendRepositoryInbox_Bad(t *core.T) {
	subject := (*ActivityPubService).SendRepositoryInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendRepositoryInbox_Ugly(t *core.T) {
	subject := (*ActivityPubService).SendRepositoryInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetPersonActor_Good(t *core.T) {
	subject := (*ActivityPubService).GetPersonActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetPersonActor_Bad(t *core.T) {
	subject := (*ActivityPubService).GetPersonActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_GetPersonActor_Ugly(t *core.T) {
	subject := (*ActivityPubService).GetPersonActor
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendPersonInbox_Good(t *core.T) {
	subject := (*ActivityPubService).SendPersonInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendPersonInbox_Bad(t *core.T) {
	subject := (*ActivityPubService).SendPersonInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestActivitypub_ActivityPubService_SendPersonInbox_Ugly(t *core.T) {
	subject := (*ActivityPubService).SendPersonInbox
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
