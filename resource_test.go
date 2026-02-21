package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test types
type testItem struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type testCreate struct {
	Name string `json:"name"`
}

type testUpdate struct {
	Name *string `json:"name,omitempty"`
}

func TestResource_Good_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/orgs/core/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]testItem{{1, "a"}, {2, "b"}})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	res := NewResource[testItem, testCreate, testUpdate](c, "/api/v1/orgs/{org}/repos")

	items, err := res.List(context.Background(), Params{"org": "core"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(items.Items) != 2 {
		t.Errorf("got %d items", len(items.Items))
	}
}

func TestResource_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/repos/core/go-forge" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(testItem{1, "go-forge"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	res := NewResource[testItem, testCreate, testUpdate](c, "/api/v1/repos/{owner}/{repo}")

	item, err := res.Get(context.Background(), Params{"owner": "core", "repo": "go-forge"})
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "go-forge" {
		t.Errorf("got name=%q", item.Name)
	}
}

func TestResource_Good_Create(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		var body testCreate
		json.NewDecoder(r.Body).Decode(&body)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(testItem{1, body.Name})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	res := NewResource[testItem, testCreate, testUpdate](c, "/api/v1/orgs/{org}/repos")

	item, err := res.Create(context.Background(), Params{"org": "core"}, &testCreate{Name: "new-repo"})
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "new-repo" {
		t.Errorf("got name=%q", item.Name)
	}
}

func TestResource_Good_Update(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		json.NewEncoder(w).Encode(testItem{1, "updated"})
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	res := NewResource[testItem, testCreate, testUpdate](c, "/api/v1/repos/{owner}/{repo}")

	name := "updated"
	item, err := res.Update(context.Background(), Params{"owner": "core", "repo": "old"}, &testUpdate{Name: &name})
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "updated" {
		t.Errorf("got name=%q", item.Name)
	}
}

func TestResource_Good_Delete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	res := NewResource[testItem, testCreate, testUpdate](c, "/api/v1/repos/{owner}/{repo}")

	err := res.Delete(context.Background(), Params{"owner": "core", "repo": "old"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestResource_Good_ListAll(t *testing.T) {
	page := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page++
		w.Header().Set("X-Total-Count", "3")
		if page == 1 {
			json.NewEncoder(w).Encode([]testItem{{1, "a"}, {2, "b"}})
		} else {
			json.NewEncoder(w).Encode([]testItem{{3, "c"}})
		}
	}))
	defer srv.Close()

	c := NewClient(srv.URL, "tok")
	res := NewResource[testItem, testCreate, testUpdate](c, "/api/v1/repos")

	items, err := res.ListAll(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 3 {
		t.Errorf("got %d items, want 3", len(items))
	}
}
