package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestLabelService_ListRepoLabels_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/labels" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Label{
			{ID: 1, Name: "bug", Color: "#d73a4a"},
			{ID: 2, Name: "feature", Color: "#0075ca"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	labels, err := f.Labels.ListRepoLabels(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 2 {
		t.Errorf("got %d labels, want 2", len(labels))
	}
	if labels[0].Name != "bug" {
		t.Errorf("got name=%q, want %q", labels[0].Name, "bug")
	}
	if labels[1].Color != "#0075ca" {
		t.Errorf("got colour=%q, want %q", labels[1].Color, "#0075ca")
	}
}

func TestLabelService_CreateRepoLabel_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/labels" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateLabelOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Name != "enhancement" {
			t.Errorf("got name=%q, want %q", opts.Name, "enhancement")
		}
		if opts.Color != "#a2eeef" {
			t.Errorf("got colour=%q, want %q", opts.Color, "#a2eeef")
		}
		json.NewEncoder(w).Encode(types.Label{
			ID:    3,
			Name:  opts.Name,
			Color: opts.Color,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	label, err := f.Labels.CreateRepoLabel(context.Background(), "core", "go-forge", &types.CreateLabelOption{
		Name:  "enhancement",
		Color: "#a2eeef",
	})
	if err != nil {
		t.Fatal(err)
	}
	if label.ID != 3 {
		t.Errorf("got id=%d, want 3", label.ID)
	}
	if label.Name != "enhancement" {
		t.Errorf("got name=%q, want %q", label.Name, "enhancement")
	}
}

func TestLabelService_GetRepoLabel_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/labels/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Label{ID: 1, Name: "bug", Color: "#d73a4a"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	label, err := f.Labels.GetRepoLabel(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
	if label.Name != "bug" {
		t.Errorf("got name=%q, want %q", label.Name, "bug")
	}
}

func TestLabelService_EditRepoLabel_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/labels/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.EditLabelOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		json.NewEncoder(w).Encode(types.Label{ID: 1, Name: opts.Name, Color: opts.Color})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	label, err := f.Labels.EditRepoLabel(context.Background(), "core", "go-forge", 1, &types.EditLabelOption{
		Name:  "critical-bug",
		Color: "#ff0000",
	})
	if err != nil {
		t.Fatal(err)
	}
	if label.Name != "critical-bug" {
		t.Errorf("got name=%q, want %q", label.Name, "critical-bug")
	}
}

func TestLabelService_DeleteRepoLabel_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/labels/1" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Labels.DeleteRepoLabel(context.Background(), "core", "go-forge", 1)
	if err != nil {
		t.Fatal(err)
	}
}

func TestLabelService_ListOrgLabels_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/labels" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Label{
			{ID: 10, Name: "org-wide", Color: "#333333"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	labels, err := f.Labels.ListOrgLabels(context.Background(), "myorg")
	if err != nil {
		t.Fatal(err)
	}
	if len(labels) != 1 {
		t.Errorf("got %d labels, want 1", len(labels))
	}
	if labels[0].Name != "org-wide" {
		t.Errorf("got name=%q, want %q", labels[0].Name, "org-wide")
	}
}

func TestLabelService_CreateOrgLabel_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/orgs/myorg/labels" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateLabelOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		json.NewEncoder(w).Encode(types.Label{ID: 11, Name: opts.Name, Color: opts.Color})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	label, err := f.Labels.CreateOrgLabel(context.Background(), "myorg", &types.CreateLabelOption{
		Name:  "priority",
		Color: "#e4e669",
	})
	if err != nil {
		t.Fatal(err)
	}
	if label.ID != 11 {
		t.Errorf("got id=%d, want 11", label.ID)
	}
	if label.Name != "priority" {
		t.Errorf("got name=%q, want %q", label.Name, "priority")
	}
}

func TestLabelService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "label not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Labels.GetRepoLabel(context.Background(), "core", "go-forge", 999)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}
