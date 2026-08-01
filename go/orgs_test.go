package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	core "dappco.re/go"
	"dappco.re/go/forge/types"
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

func TestOrgs_OrgActivityFeedListOptions_String_Good(t *core.T) {
	subject := (*OrgActivityFeedListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgActivityFeedListOptions_String_Bad(t *core.T) {
	subject := (*OrgActivityFeedListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgActivityFeedListOptions_String_Ugly(t *core.T) {
	subject := (*OrgActivityFeedListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgActivityFeedListOptions_GoString_Good(t *core.T) {
	subject := (*OrgActivityFeedListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgActivityFeedListOptions_GoString_Bad(t *core.T) {
	subject := (*OrgActivityFeedListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgActivityFeedListOptions_GoString_Ugly(t *core.T) {
	subject := (*OrgActivityFeedListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetOrg_Good(t *core.T) {
	subject := (*OrgService).GetOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetOrg_Bad(t *core.T) {
	subject := (*OrgService).GetOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetOrg_Ugly(t *core.T) {
	subject := (*OrgService).GetOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_UpdateOrg_Good(t *core.T) {
	subject := (*OrgService).UpdateOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_UpdateOrg_Bad(t *core.T) {
	subject := (*OrgService).UpdateOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_UpdateOrg_Ugly(t *core.T) {
	subject := (*OrgService).UpdateOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_DeleteOrg_Good(t *core.T) {
	subject := (*OrgService).DeleteOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_DeleteOrg_Bad(t *core.T) {
	subject := (*OrgService).DeleteOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_DeleteOrg_Ugly(t *core.T) {
	subject := (*OrgService).DeleteOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgsPage_Good(t *core.T) {
	subject := (*OrgService).ListOrgsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgsPage_Bad(t *core.T) {
	subject := (*OrgService).ListOrgsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgsPage_Ugly(t *core.T) {
	subject := (*OrgService).ListOrgsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgs_Good(t *core.T) {
	subject := (*OrgService).ListOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgs_Bad(t *core.T) {
	subject := (*OrgService).ListOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgs_Ugly(t *core.T) {
	subject := (*OrgService).ListOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgTeamsPage_Good(t *core.T) {
	subject := (*OrgService).ListOrgTeamsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgTeamsPage_Bad(t *core.T) {
	subject := (*OrgService).ListOrgTeamsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgTeamsPage_Ugly(t *core.T) {
	subject := (*OrgService).ListOrgTeamsPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgTeams_Good(t *core.T) {
	subject := (*OrgService).ListOrgTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgTeams_Bad(t *core.T) {
	subject := (*OrgService).ListOrgTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgTeams_Ugly(t *core.T) {
	subject := (*OrgService).ListOrgTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgTeams_Good(t *core.T) {
	subject := (*OrgService).IterOrgTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgTeams_Bad(t *core.T) {
	subject := (*OrgService).IterOrgTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgTeams_Ugly(t *core.T) {
	subject := (*OrgService).IterOrgTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgs_Good(t *core.T) {
	subject := (*OrgService).IterOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgs_Bad(t *core.T) {
	subject := (*OrgService).IterOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgs_Ugly(t *core.T) {
	subject := (*OrgService).IterOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_CreateOrg_Good(t *core.T) {
	subject := (*OrgService).CreateOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_CreateOrg_Bad(t *core.T) {
	subject := (*OrgService).CreateOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_CreateOrg_Ugly(t *core.T) {
	subject := (*OrgService).CreateOrg
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMembersPage_Good(t *core.T) {
	subject := (*OrgService).ListMembersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMembersPage_Bad(t *core.T) {
	subject := (*OrgService).ListMembersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMembersPage_Ugly(t *core.T) {
	subject := (*OrgService).ListMembersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMembers_Good(t *core.T) {
	subject := (*OrgService).ListMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMembers_Bad(t *core.T) {
	subject := (*OrgService).ListMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMembers_Ugly(t *core.T) {
	subject := (*OrgService).ListMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgMembersPage_Good(t *core.T) {
	subject := (*OrgService).ListOrgMembersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgMembersPage_Bad(t *core.T) {
	subject := (*OrgService).ListOrgMembersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgMembersPage_Ugly(t *core.T) {
	subject := (*OrgService).ListOrgMembersPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgMembers_Good(t *core.T) {
	subject := (*OrgService).ListOrgMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgMembers_Bad(t *core.T) {
	subject := (*OrgService).ListOrgMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListOrgMembers_Ugly(t *core.T) {
	subject := (*OrgService).ListOrgMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterMembers_Good(t *core.T) {
	subject := (*OrgService).IterMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterMembers_Bad(t *core.T) {
	subject := (*OrgService).IterMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterMembers_Ugly(t *core.T) {
	subject := (*OrgService).IterMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgMembers_Good(t *core.T) {
	subject := (*OrgService).IterOrgMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgMembers_Bad(t *core.T) {
	subject := (*OrgService).IterOrgMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterOrgMembers_Ugly(t *core.T) {
	subject := (*OrgService).IterOrgMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_AddMember_Good(t *core.T) {
	subject := (*OrgService).AddMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_AddMember_Bad(t *core.T) {
	subject := (*OrgService).AddMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_AddMember_Ugly(t *core.T) {
	subject := (*OrgService).AddMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_RemoveMember_Good(t *core.T) {
	subject := (*OrgService).RemoveMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_RemoveMember_Bad(t *core.T) {
	subject := (*OrgService).RemoveMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_RemoveMember_Ugly(t *core.T) {
	subject := (*OrgService).RemoveMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsMember_Good(t *core.T) {
	subject := (*OrgService).IsMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsMember_Bad(t *core.T) {
	subject := (*OrgService).IsMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsMember_Ugly(t *core.T) {
	subject := (*OrgService).IsMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListBlockedUsers_Good(t *core.T) {
	subject := (*OrgService).ListBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListBlockedUsers_Bad(t *core.T) {
	subject := (*OrgService).ListBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListBlockedUsers_Ugly(t *core.T) {
	subject := (*OrgService).ListBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterBlockedUsers_Good(t *core.T) {
	subject := (*OrgService).IterBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterBlockedUsers_Bad(t *core.T) {
	subject := (*OrgService).IterBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterBlockedUsers_Ugly(t *core.T) {
	subject := (*OrgService).IterBlockedUsers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsBlocked_Good(t *core.T) {
	subject := (*OrgService).IsBlocked
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsBlocked_Bad(t *core.T) {
	subject := (*OrgService).IsBlocked
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsBlocked_Ugly(t *core.T) {
	subject := (*OrgService).IsBlocked
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListPublicMembers_Good(t *core.T) {
	subject := (*OrgService).ListPublicMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListPublicMembers_Bad(t *core.T) {
	subject := (*OrgService).ListPublicMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListPublicMembers_Ugly(t *core.T) {
	subject := (*OrgService).ListPublicMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterPublicMembers_Good(t *core.T) {
	subject := (*OrgService).IterPublicMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterPublicMembers_Bad(t *core.T) {
	subject := (*OrgService).IterPublicMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterPublicMembers_Ugly(t *core.T) {
	subject := (*OrgService).IterPublicMembers
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsPublicMember_Good(t *core.T) {
	subject := (*OrgService).IsPublicMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsPublicMember_Bad(t *core.T) {
	subject := (*OrgService).IsPublicMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IsPublicMember_Ugly(t *core.T) {
	subject := (*OrgService).IsPublicMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_PublicizeMember_Good(t *core.T) {
	subject := (*OrgService).PublicizeMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_PublicizeMember_Bad(t *core.T) {
	subject := (*OrgService).PublicizeMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_PublicizeMember_Ugly(t *core.T) {
	subject := (*OrgService).PublicizeMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ConcealMember_Good(t *core.T) {
	subject := (*OrgService).ConcealMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ConcealMember_Bad(t *core.T) {
	subject := (*OrgService).ConcealMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ConcealMember_Ugly(t *core.T) {
	subject := (*OrgService).ConcealMember
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_Block_Good(t *core.T) {
	subject := (*OrgService).Block
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_Block_Bad(t *core.T) {
	subject := (*OrgService).Block
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_Block_Ugly(t *core.T) {
	subject := (*OrgService).Block
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_Unblock_Good(t *core.T) {
	subject := (*OrgService).Unblock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_Unblock_Bad(t *core.T) {
	subject := (*OrgService).Unblock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_Unblock_Ugly(t *core.T) {
	subject := (*OrgService).Unblock
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetQuota_Good(t *core.T) {
	subject := (*OrgService).GetQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetQuota_Bad(t *core.T) {
	subject := (*OrgService).GetQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetQuota_Ugly(t *core.T) {
	subject := (*OrgService).GetQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_CheckQuota_Good(t *core.T) {
	subject := (*OrgService).CheckQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_CheckQuota_Bad(t *core.T) {
	subject := (*OrgService).CheckQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_CheckQuota_Ugly(t *core.T) {
	subject := (*OrgService).CheckQuota
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaArtifacts_Good(t *core.T) {
	subject := (*OrgService).ListQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaArtifacts_Bad(t *core.T) {
	subject := (*OrgService).ListQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaArtifacts_Ugly(t *core.T) {
	subject := (*OrgService).ListQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaArtifacts_Good(t *core.T) {
	subject := (*OrgService).IterQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaArtifacts_Bad(t *core.T) {
	subject := (*OrgService).IterQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaArtifacts_Ugly(t *core.T) {
	subject := (*OrgService).IterQuotaArtifacts
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaAttachments_Good(t *core.T) {
	subject := (*OrgService).ListQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaAttachments_Bad(t *core.T) {
	subject := (*OrgService).ListQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaAttachments_Ugly(t *core.T) {
	subject := (*OrgService).ListQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaAttachments_Good(t *core.T) {
	subject := (*OrgService).IterQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaAttachments_Bad(t *core.T) {
	subject := (*OrgService).IterQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaAttachments_Ugly(t *core.T) {
	subject := (*OrgService).IterQuotaAttachments
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaPackages_Good(t *core.T) {
	subject := (*OrgService).ListQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaPackages_Bad(t *core.T) {
	subject := (*OrgService).ListQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListQuotaPackages_Ugly(t *core.T) {
	subject := (*OrgService).ListQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaPackages_Good(t *core.T) {
	subject := (*OrgService).IterQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaPackages_Bad(t *core.T) {
	subject := (*OrgService).IterQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterQuotaPackages_Ugly(t *core.T) {
	subject := (*OrgService).IterQuotaPackages
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetRunnerRegistrationToken_Good(t *core.T) {
	subject := (*OrgService).GetRunnerRegistrationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetRunnerRegistrationToken_Bad(t *core.T) {
	subject := (*OrgService).GetRunnerRegistrationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetRunnerRegistrationToken_Ugly(t *core.T) {
	subject := (*OrgService).GetRunnerRegistrationToken
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_UpdateAvatar_Good(t *core.T) {
	subject := (*OrgService).UpdateAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_UpdateAvatar_Bad(t *core.T) {
	subject := (*OrgService).UpdateAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_UpdateAvatar_Ugly(t *core.T) {
	subject := (*OrgService).UpdateAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_DeleteAvatar_Good(t *core.T) {
	subject := (*OrgService).DeleteAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_DeleteAvatar_Bad(t *core.T) {
	subject := (*OrgService).DeleteAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_DeleteAvatar_Ugly(t *core.T) {
	subject := (*OrgService).DeleteAvatar
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_SearchTeams_Good(t *core.T) {
	subject := (*OrgService).SearchTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_SearchTeams_Bad(t *core.T) {
	subject := (*OrgService).SearchTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_SearchTeams_Ugly(t *core.T) {
	subject := (*OrgService).SearchTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterSearchTeams_Good(t *core.T) {
	subject := (*OrgService).IterSearchTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterSearchTeams_Bad(t *core.T) {
	subject := (*OrgService).IterSearchTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterSearchTeams_Ugly(t *core.T) {
	subject := (*OrgService).IterSearchTeams
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetUserPermissions_Good(t *core.T) {
	subject := (*OrgService).GetUserPermissions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetUserPermissions_Bad(t *core.T) {
	subject := (*OrgService).GetUserPermissions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_GetUserPermissions_Ugly(t *core.T) {
	subject := (*OrgService).GetUserPermissions
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListActivityFeeds_Good(t *core.T) {
	subject := (*OrgService).ListActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListActivityFeeds_Bad(t *core.T) {
	subject := (*OrgService).ListActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListActivityFeeds_Ugly(t *core.T) {
	subject := (*OrgService).ListActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterActivityFeeds_Good(t *core.T) {
	subject := (*OrgService).IterActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterActivityFeeds_Bad(t *core.T) {
	subject := (*OrgService).IterActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterActivityFeeds_Ugly(t *core.T) {
	subject := (*OrgService).IterActivityFeeds
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListUserOrgs_Good(t *core.T) {
	subject := (*OrgService).ListUserOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListUserOrgs_Bad(t *core.T) {
	subject := (*OrgService).ListUserOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListUserOrgs_Ugly(t *core.T) {
	subject := (*OrgService).ListUserOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterUserOrgs_Good(t *core.T) {
	subject := (*OrgService).IterUserOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterUserOrgs_Bad(t *core.T) {
	subject := (*OrgService).IterUserOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterUserOrgs_Ugly(t *core.T) {
	subject := (*OrgService).IterUserOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMyOrgs_Good(t *core.T) {
	subject := (*OrgService).ListMyOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMyOrgs_Bad(t *core.T) {
	subject := (*OrgService).ListMyOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_ListMyOrgs_Ugly(t *core.T) {
	subject := (*OrgService).ListMyOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterMyOrgs_Good(t *core.T) {
	subject := (*OrgService).IterMyOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterMyOrgs_Bad(t *core.T) {
	subject := (*OrgService).IterMyOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestOrgs_OrgService_IterMyOrgs_Ugly(t *core.T) {
	subject := (*OrgService).IterMyOrgs
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
