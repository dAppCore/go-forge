package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

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
