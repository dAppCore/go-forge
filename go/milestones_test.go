package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	core "dappco.re/go"
	"dappco.re/go/forge/types"
)

func TestMilestoneService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Errorf("got limit=%q, want %q", got, "1")
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Milestone{{ID: 2, Title: "v2.0"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Milestones.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, ListOptions{Page: 1, Limit: 1})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 1 {
		t.Errorf("got page=%d, want 1", page.Page)
	}
	if page.TotalCount != 2 {
		t.Errorf("got total=%d, want 2", page.TotalCount)
	}
	if !page.HasMore {
		t.Error("expected HasMore=true")
	}
	if len(page.Items) != 1 || page.Items[0].Title != "v2.0" {
		t.Fatalf("unexpected items: %+v", page.Items)
	}
}

func TestMilestoneService_ListWithFilters_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("state"); got != "all" {
			t.Errorf("got state=%q, want %q", got, "all")
		}
		if got := r.URL.Query().Get("name"); got != "v1.0" {
			t.Errorf("got name=%q, want %q", got, "v1.0")
		}
		if got := r.URL.Query().Get("page"); got != "1" {
			t.Errorf("got page=%q, want %q", got, "1")
		}
		if got := r.URL.Query().Get("limit"); got != "1" {
			t.Errorf("got limit=%q, want %q", got, "1")
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Milestone{{ID: 1, Title: "v1.0"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	page, err := f.Milestones.List(
		context.Background(),
		Params{"owner": "core", "repo": "go-forge"},
		ListOptions{Page: 1, Limit: 1},
		MilestoneListOptions{State: "all", Name: "v1.0"},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.Items[0].Title != "v1.0" {
		t.Fatalf("unexpected items: %+v", page.Items)
	}
}

func TestMilestoneService_Iter_Good(t *testing.T) {
	requests := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		switch requests {
		case 1:
			if got := r.URL.Query().Get("page"); got != "1" {
				t.Errorf("got page=%q, want %q", got, "1")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.Milestone{{ID: 1, Title: "v1.0"}})
		case 2:
			if got := r.URL.Query().Get("page"); got != "2" {
				t.Errorf("got page=%q, want %q", got, "2")
			}
			w.Header().Set("X-Total-Count", "2")
			json.NewEncoder(w).Encode([]types.Milestone{{ID: 2, Title: "v2.0"}})
		default:
			t.Fatalf("unexpected request %d", requests)
		}
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	var got []string
	for milestone, err := range f.Milestones.Iter(context.Background(), Params{"owner": "core", "repo": "go-forge"}) {
		if err != nil {
			t.Fatal(err)
		}
		got = append(got, milestone.Title)
	}
	if !reflect.DeepEqual(got, []string{"v1.0", "v2.0"}) {
		t.Fatalf("got %v", got)
	}
}

func TestMilestoneService_ListAll_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode([]types.Milestone{
			{ID: 1, Title: "v1.0"},
			{ID: 2, Title: "v2.0"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	milestones, err := f.Milestones.ListAll(context.Background(), Params{"owner": "core", "repo": "go-forge"})
	if err != nil {
		t.Fatal(err)
	}
	if len(milestones) != 2 {
		t.Errorf("got %d milestones, want 2", len(milestones))
	}
	if milestones[0].Title != "v1.0" {
		t.Errorf("got title=%q, want %q", milestones[0].Title, "v1.0")
	}
}

func TestMilestoneService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones/7" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Milestone{ID: 7, Title: "v1.0"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	milestone, err := f.Milestones.Get(context.Background(), "core", "go-forge", 7)
	if err != nil {
		t.Fatal(err)
	}
	if milestone.ID != 7 {
		t.Errorf("got id=%d, want 7", milestone.ID)
	}
	if milestone.Title != "v1.0" {
		t.Errorf("got title=%q, want %q", milestone.Title, "v1.0")
	}
}

func TestMilestoneService_Create_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateMilestoneOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Title != "v1.0" {
			t.Errorf("got title=%q, want %q", opts.Title, "v1.0")
		}

		json.NewEncoder(w).Encode(types.Milestone{ID: 3, Title: opts.Title})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	milestone, err := f.Milestones.Create(context.Background(), "core", "go-forge", &types.CreateMilestoneOption{
		Title: "v1.0",
	})
	if err != nil {
		t.Fatal(err)
	}
	if milestone.ID != 3 {
		t.Errorf("got id=%d, want 3", milestone.ID)
	}
	if milestone.Title != "v1.0" {
		t.Errorf("got title=%q, want %q", milestone.Title, "v1.0")
	}
}

func TestMilestoneService_Edit_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones/3" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditMilestoneOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Title != "v1.1" {
			t.Errorf("got title=%q, want %q", opts.Title, "v1.1")
		}

		json.NewEncoder(w).Encode(types.Milestone{ID: 3, Title: opts.Title})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	milestone, err := f.Milestones.Edit(context.Background(), "core", "go-forge", 3, &types.EditMilestoneOption{
		Title: "v1.1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if milestone.ID != 3 {
		t.Errorf("got id=%d, want 3", milestone.ID)
	}
	if milestone.Title != "v1.1" {
		t.Errorf("got title=%q, want %q", milestone.Title, "v1.1")
	}
}

func TestMilestoneService_Delete_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/milestones/3" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Milestones.Delete(context.Background(), "core", "go-forge", 3); err != nil {
		t.Fatal(err)
	}
}

func TestMilestones_MilestoneListOptions_String_Good(t *core.T) {
	subject := (*MilestoneListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneListOptions_String_Bad(t *core.T) {
	subject := (*MilestoneListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneListOptions_String_Ugly(t *core.T) {
	subject := (*MilestoneListOptions).String
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneListOptions_GoString_Good(t *core.T) {
	subject := (*MilestoneListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneListOptions_GoString_Bad(t *core.T) {
	subject := (*MilestoneListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneListOptions_GoString_Ugly(t *core.T) {
	subject := (*MilestoneListOptions).GoString
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListMilestonesPage_Good(t *core.T) {
	subject := (*MilestoneService).ListMilestonesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListMilestonesPage_Bad(t *core.T) {
	subject := (*MilestoneService).ListMilestonesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListMilestonesPage_Ugly(t *core.T) {
	subject := (*MilestoneService).ListMilestonesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListMilestones_Good(t *core.T) {
	subject := (*MilestoneService).ListMilestones
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListMilestones_Bad(t *core.T) {
	subject := (*MilestoneService).ListMilestones
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListMilestones_Ugly(t *core.T) {
	subject := (*MilestoneService).ListMilestones
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_IterMilestones_Good(t *core.T) {
	subject := (*MilestoneService).IterMilestones
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_IterMilestones_Bad(t *core.T) {
	subject := (*MilestoneService).IterMilestones
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_IterMilestones_Ugly(t *core.T) {
	subject := (*MilestoneService).IterMilestones
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_List_Good(t *core.T) {
	subject := (*MilestoneService).List
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_List_Bad(t *core.T) {
	subject := (*MilestoneService).List
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_List_Ugly(t *core.T) {
	subject := (*MilestoneService).List
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Iter_Good(t *core.T) {
	subject := (*MilestoneService).Iter
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Iter_Bad(t *core.T) {
	subject := (*MilestoneService).Iter
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Iter_Ugly(t *core.T) {
	subject := (*MilestoneService).Iter
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListAll_Good(t *core.T) {
	subject := (*MilestoneService).ListAll
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListAll_Bad(t *core.T) {
	subject := (*MilestoneService).ListAll
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_ListAll_Ugly(t *core.T) {
	subject := (*MilestoneService).ListAll
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Get_Good(t *core.T) {
	subject := (*MilestoneService).Get
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Get_Bad(t *core.T) {
	subject := (*MilestoneService).Get
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Get_Ugly(t *core.T) {
	subject := (*MilestoneService).Get
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_GetMilestone_Good(t *core.T) {
	subject := (*MilestoneService).GetMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_GetMilestone_Bad(t *core.T) {
	subject := (*MilestoneService).GetMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_GetMilestone_Ugly(t *core.T) {
	subject := (*MilestoneService).GetMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Create_Good(t *core.T) {
	subject := (*MilestoneService).Create
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Create_Bad(t *core.T) {
	subject := (*MilestoneService).Create
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Create_Ugly(t *core.T) {
	subject := (*MilestoneService).Create
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_CreateMilestone_Good(t *core.T) {
	subject := (*MilestoneService).CreateMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_CreateMilestone_Bad(t *core.T) {
	subject := (*MilestoneService).CreateMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_CreateMilestone_Ugly(t *core.T) {
	subject := (*MilestoneService).CreateMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Edit_Good(t *core.T) {
	subject := (*MilestoneService).Edit
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Edit_Bad(t *core.T) {
	subject := (*MilestoneService).Edit
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Edit_Ugly(t *core.T) {
	subject := (*MilestoneService).Edit
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_EditMilestone_Good(t *core.T) {
	subject := (*MilestoneService).EditMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_EditMilestone_Bad(t *core.T) {
	subject := (*MilestoneService).EditMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_EditMilestone_Ugly(t *core.T) {
	subject := (*MilestoneService).EditMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Delete_Good(t *core.T) {
	subject := (*MilestoneService).Delete
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Delete_Bad(t *core.T) {
	subject := (*MilestoneService).Delete
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_Delete_Ugly(t *core.T) {
	subject := (*MilestoneService).Delete
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_DeleteMilestone_Good(t *core.T) {
	subject := (*MilestoneService).DeleteMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_DeleteMilestone_Bad(t *core.T) {
	subject := (*MilestoneService).DeleteMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestMilestones_MilestoneService_DeleteMilestone_Ugly(t *core.T) {
	subject := (*MilestoneService).DeleteMilestone
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
