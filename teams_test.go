package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestTeamService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/teams/42" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Team{ID: 42, Name: "developers"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	team, err := f.Teams.Get(context.Background(), Params{"id": "42"})
	if err != nil {
		t.Fatal(err)
	}
	if team.Name != "developers" {
		t.Errorf("got name=%q, want %q", team.Name, "developers")
	}
}

func TestTeamService_CreateOrgTeam_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/teams" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateTeamOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Name != "platform" {
			t.Errorf("got name=%q, want %q", opts.Name, "platform")
		}
		json.NewEncoder(w).Encode(types.Team{ID: 7, Name: opts.Name})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	team, err := f.Teams.CreateOrgTeam(context.Background(), "core", &types.CreateTeamOption{
		Name: "platform",
	})
	if err != nil {
		t.Fatal(err)
	}
	if team.ID != 7 || team.Name != "platform" {
		t.Fatalf("got %#v", team)
	}
}

func TestTeamService_ListMembers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/teams/42/members" {
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
	members, err := f.Teams.ListMembers(context.Background(), 42)
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

func TestTeamService_AddMember_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/teams/42/members/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Teams.AddMember(context.Background(), 42, "alice")
	if err != nil {
		t.Fatal(err)
	}
}

func TestTeamService_GetMember_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/teams/42/members/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.User{ID: 1, UserName: "alice"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	member, err := f.Teams.GetMember(context.Background(), 42, "alice")
	if err != nil {
		t.Fatal(err)
	}
	if member.UserName != "alice" {
		t.Errorf("got username=%q, want %q", member.UserName, "alice")
	}
}
