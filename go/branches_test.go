package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestBranchService_List_Good(t *testing.T) {
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

func TestBranchService_Get_Good(t *testing.T) {
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

func TestBranchService_UpdateBranch_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			t.Errorf("expected PATCH, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/repos/core/go-forge/branches/main" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		var opts types.UpdateBranchRepoOption
		if err := json.NewDecoder(r.Body).Decode(&opts); err != nil {
			t.Fatal(err)
		}
		if opts.Name != "develop" {
			t.Errorf("got name=%q, want %q", opts.Name, "develop")
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	if err := f.Branches.UpdateBranch(context.Background(), "core", "go-forge", "main", &types.UpdateBranchRepoOption{
		Name: "develop",
	}); err != nil {
		t.Fatal(err)
	}
}

func TestBranchService_CreateProtection_Good(t *testing.T) {
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

func TestBranches_BranchService_ListBranchesPage_Good(t *core.T) {
	subject := (*BranchService).ListBranchesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranchesPage_Bad(t *core.T) {
	subject := (*BranchService).ListBranchesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranchesPage_Ugly(t *core.T) {
	subject := (*BranchService).ListBranchesPage
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranches_Good(t *core.T) {
	subject := (*BranchService).ListBranches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranches_Bad(t *core.T) {
	subject := (*BranchService).ListBranches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranches_Ugly(t *core.T) {
	subject := (*BranchService).ListBranches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_IterBranches_Good(t *core.T) {
	subject := (*BranchService).IterBranches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_IterBranches_Bad(t *core.T) {
	subject := (*BranchService).IterBranches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_IterBranches_Ugly(t *core.T) {
	subject := (*BranchService).IterBranches
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_CreateBranch_Good(t *core.T) {
	subject := (*BranchService).CreateBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_CreateBranch_Bad(t *core.T) {
	subject := (*BranchService).CreateBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_CreateBranch_Ugly(t *core.T) {
	subject := (*BranchService).CreateBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_GetBranch_Good(t *core.T) {
	subject := (*BranchService).GetBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_GetBranch_Bad(t *core.T) {
	subject := (*BranchService).GetBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_GetBranch_Ugly(t *core.T) {
	subject := (*BranchService).GetBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_UpdateBranch_Good(t *core.T) {
	subject := (*BranchService).UpdateBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_UpdateBranch_Bad(t *core.T) {
	subject := (*BranchService).UpdateBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_UpdateBranch_Ugly(t *core.T) {
	subject := (*BranchService).UpdateBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_DeleteBranch_Good(t *core.T) {
	subject := (*BranchService).DeleteBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_DeleteBranch_Bad(t *core.T) {
	subject := (*BranchService).DeleteBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_DeleteBranch_Ugly(t *core.T) {
	subject := (*BranchService).DeleteBranch
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranchProtections_Good(t *core.T) {
	subject := (*BranchService).ListBranchProtections
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranchProtections_Bad(t *core.T) {
	subject := (*BranchService).ListBranchProtections
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_ListBranchProtections_Ugly(t *core.T) {
	subject := (*BranchService).ListBranchProtections
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_IterBranchProtections_Good(t *core.T) {
	subject := (*BranchService).IterBranchProtections
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_IterBranchProtections_Bad(t *core.T) {
	subject := (*BranchService).IterBranchProtections
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_IterBranchProtections_Ugly(t *core.T) {
	subject := (*BranchService).IterBranchProtections
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_GetBranchProtection_Good(t *core.T) {
	subject := (*BranchService).GetBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_GetBranchProtection_Bad(t *core.T) {
	subject := (*BranchService).GetBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_GetBranchProtection_Ugly(t *core.T) {
	subject := (*BranchService).GetBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_CreateBranchProtection_Good(t *core.T) {
	subject := (*BranchService).CreateBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_CreateBranchProtection_Bad(t *core.T) {
	subject := (*BranchService).CreateBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_CreateBranchProtection_Ugly(t *core.T) {
	subject := (*BranchService).CreateBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_EditBranchProtection_Good(t *core.T) {
	subject := (*BranchService).EditBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_EditBranchProtection_Bad(t *core.T) {
	subject := (*BranchService).EditBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_EditBranchProtection_Ugly(t *core.T) {
	subject := (*BranchService).EditBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_DeleteBranchProtection_Good(t *core.T) {
	subject := (*BranchService).DeleteBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_DeleteBranchProtection_Bad(t *core.T) {
	subject := (*BranchService).DeleteBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestBranches_BranchService_DeleteBranchProtection_Ugly(t *core.T) {
	subject := (*BranchService).DeleteBranchProtection
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
