package forge

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"forge.lthn.ai/core/go-forge/types"
)

func TestBranchService_Good_List(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Branch{
			{Name: "main", Protected: true},
			{Name: "develop", Protected: false},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	result, err := f.Branches.List(context.Background(), Params{"owner": "core", "repo": "go-forge"}, DefaultList)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Items) != 2 {
		t.Errorf("got %d items, want 2", len(result.Items))
	}
	if result.Items[0].Name != "main" {
		t.Errorf("got name=%q, want %q", result.Items[0].Name, "main")
	}
}

func TestBranchService_Good_Get(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/branches/main" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Branch{Name: "main", Protected: true})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	branch, err := f.Branches.Get(context.Background(), Params{"owner": "core", "repo": "go-forge", "branch": "main"})
	if err != nil {
		t.Fatal(err)
	}
	if branch.Name != "main" {
		t.Errorf("got name=%q, want %q", branch.Name, "main")
	}
	if !branch.Protected {
		t.Error("expected branch to be protected")
	}
}

func TestBranchService_Good_CreateProtection(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/branch_protections" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.CreateBranchProtectionOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.RuleName != "main" {
			t.Errorf("got rule_name=%q, want %q", opts.RuleName, "main")
		}
		json.NewEncoder(w).Encode(types.BranchProtection{
			RuleName:          "main",
			EnablePush:        true,
			RequiredApprovals: 2,
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	bp, err := f.Branches.CreateBranchProtection(context.Background(), "core", "go-forge", &types.CreateBranchProtectionOption{
		RuleName:          "main",
		EnablePush:        true,
		RequiredApprovals: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if bp.RuleName != "main" {
		t.Errorf("got rule_name=%q, want %q", bp.RuleName, "main")
	}
	if bp.RequiredApprovals != 2 {
		t.Errorf("got required_approvals=%d, want 2", bp.RequiredApprovals)
	}
}
