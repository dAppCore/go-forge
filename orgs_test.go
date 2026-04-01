package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestOrgService_List_Good(t *testing.T) {
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

func TestOrgService_Get_Good(t *testing.T) {
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

func TestOrgService_ListMembers_Good(t *testing.T) {
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

func TestOrgService_ListBlockedUsers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/blocks" {
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
	blocked, err := f.Orgs.ListBlockedUsers(context.Background(), "core")
	if err != nil {
		t.Fatal(err)
	}
	if len(blocked) != 2 {
		t.Fatalf("got %d blocked users, want 2", len(blocked))
	}
	if blocked[0].UserName != "alice" {
		t.Errorf("got username=%q, want %q", blocked[0].UserName, "alice")
	}
}

func TestOrgService_Block_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/blocks/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Orgs.Block(context.Background(), "core", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestOrgService_Unblock_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/blocks/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Orgs.Unblock(context.Background(), "core", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestOrgService_IsBlocked_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/blocks/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	blocked, err := f.Orgs.IsBlocked(context.Background(), "core", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if !blocked {
		t.Fatal("got blocked=false, want true")
	}
}
