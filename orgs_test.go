package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forge.lthn.ai/core/go-forge/types"
)

func TestOrgService_Good_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Organization{
			{ID: 1, Name: "core"},
			{ID: 2, Name: "labs"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Orgs.ListAll(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(result) != 2 {
		t.Errorf("got %d items, want 2", len(result))
	}
	if result[0].Name != "core" {
		t.Errorf("got name=%q, want %q", result[0].Name, "core")
	}
}

func TestOrgService_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Organization{ID: 1, Name: "core"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	org, err := f.Orgs.Get(context.Background(), Params{"org": "core"})
	if err != nil {
		t.Fatal(err)
	}
	if org.Name != "core" {
		t.Errorf("got name=%q, want %q", org.Name, "core")
	}
}

func TestOrgService_Good_ListMembers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/members" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.User{
			{ID: 1, UserName: "alice"},
			{ID: 2, UserName: "bob"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	members, err := f.Orgs.ListMembers(context.Background(), "core")
	if err != nil {
		t.Fatal(err)
	}
	if len(members) != 2 {
		t.Errorf("got %d members, want 2", len(members))
	}
	if members[0].UserName != "alice" {
		t.Errorf("got username=%q, want %q", members[0].UserName, "alice")
	}
}
