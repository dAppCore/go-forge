package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/core/forge/types"
)

func TestBranchService_ListBranches_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/branches" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.Branch{{Name: "main"}})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	branches, err := f.Branches.ListBranches(context.Background(), "core", "go-forge")
	if err != nil {
		t.Fatal(err)
	}
	if len(branches) != 1 || branches[0].Name != "main" {
		t.Fatalf("got %#v", branches)
	}
}

func TestBranchService_CreateBranch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/branches" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var body types.CreateBranchRepoOption
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body.BranchName != "release/v1" {
			t.Fatalf("unexpected body: %+v", body)
		}
		json.NewEncoder(w).Encode(types.Branch{Name: body.BranchName})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	branch, err := f.Branches.CreateBranch(context.Background(), "core", "go-forge", &types.CreateBranchRepoOption{
		BranchName: "release/v1",
		OldRefName: "main",
	})
	if err != nil {
		t.Fatal(err)
	}
	if branch.Name != "release/v1" {
		t.Fatalf("got name=%q", branch.Name)
	}
}
