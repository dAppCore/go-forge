package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestOrgService_IsMember_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/members/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	member, err := f.Orgs.IsMember(context.Background(), "core", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if !member {
		t.Fatal("got member=false, want true")
	}
}

func TestOrgService_ListPublicMembers_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/public_members" {
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
	members, err := f.Orgs.ListPublicMembers(context.Background(), "core")
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
		if r.URL.Path != "/api/v1/orgs/core/list_blocked" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.BlockedUser{
			{BlockID: 1},
			{BlockID: 2},
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
	if blocked[0].BlockID != 1 {
		t.Errorf("got block_id=%d, want %d", blocked[0].BlockID, 1)
	}
}

func TestOrgService_PublicizeMember_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/public_members/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Orgs.PublicizeMember(context.Background(), "core", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestOrgService_ConcealMember_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/public_members/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Orgs.ConcealMember(context.Background(), "core", "alice"); err != nil {
		t.Fatal(err)
	}
}

func TestOrgService_Block_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/block/alice" {
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
		if r.URL.Path != "/api/v1/orgs/core/unblock/alice" {
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

func TestOrgService_ListActivityFeeds_Good(t *testing.T) {
	date := time.Date(2026, time.April, 2, 15, 4, 5, 0, time.UTC)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/activities/feeds" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("date"); got != "2026-04-02" {
			t.Errorf("wrong date: %s", got)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Activity{{
			ID:      9,
			OpType:  "create_org",
			Content: "created organisation",
		}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	activities, err := f.Orgs.ListActivityFeeds(context.Background(), "core", OrgActivityFeedListOptions{Date: &date})
	if err != nil {
		t.Fatal(err)
	}
	if len(activities) != 1 || activities[0].ID != 9 || activities[0].OpType != "create_org" {
		t.Fatalf("got %#v", activities)
	}
}

func TestOrgService_IterActivityFeeds_Good(t *testing.T) {
	var requests int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/activities/feeds" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Activity{{
			ID:      11,
			OpType:  "update_org",
			Content: "updated organisation",
		}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []int64
	for activity, err := range f.Orgs.IterActivityFeeds(context.Background(), "core") {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, activity.ID)
	}
	if requests != 1 {
		t.Fatalf("expected 1 request, got %d", requests)
	}
	if len(got) != 1 || got[0] != 11 {
		t.Fatalf("got %#v", got)
	}
}

func TestOrgService_IsBlocked_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/block/alice" {
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

func TestOrgService_IsPublicMember_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/core/public_members/alice" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	public, err := f.Orgs.IsPublicMember(context.Background(), "core", "alice")
	if err != nil {
		t.Fatal(err)
	}
	if !public {
		t.Fatal("got public=false, want true")
	}
}
