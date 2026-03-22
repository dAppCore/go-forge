package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestUserService_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.User{ID: 1, UserName: "alice"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Users.Get(context.Background(), Params{"username": "alice"})
	if err != nil {
		t.Fatal(err)
	}
	if user.UserName != "alice" {
		t.Errorf("got username=%q, want %q", user.UserName, "alice")
	}
}

func TestUserService_Good_GetCurrent(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/user" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.User{ID: 1, UserName: "me"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	user, err := f.Users.GetCurrent(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if user.UserName != "me" {
		t.Errorf("got username=%q, want %q", user.UserName, "me")
	}
}

func TestUserService_Good_ListFollowers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/users/alice/followers" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.User{
			{ID: 2, UserName: "bob"},
			{ID: 3, UserName: "charlie"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	followers, err := f.Users.ListFollowers(context.Background(), "alice")
	if err != nil {
		t.Fatal(err)
	}
	if len(followers) != 2 {
		t.Errorf("got %d followers, want 2", len(followers))
	}
	if followers[0].UserName != "bob" {
		t.Errorf("got username=%q, want %q", followers[0].UserName, "bob")
	}
}
