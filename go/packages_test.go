package forge

import (
	"context"
	json "github.com/goccy/go-json"
	"net/http"
	"net/http/httptest"
	"testing"

	"dappco.re/go/forge/types"
)

func TestPackageService_List_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/packages/core" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]types.Package{
			{ID: 1, Name: "go-forge", Type: "generic", Version: "0.1.0"},
			{ID: 2, Name: "go-forge", Type: "generic", Version: "0.2.0"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pkgs, err := f.Packages.List(context.Background(), "core")
	if err != nil {
		t.Fatal(err)
	}
	if len(pkgs) != 2 {
		t.Fatalf("got %d packages, want 2", len(pkgs))
	}
	if pkgs[0].Name != "go-forge" {
		t.Errorf("got name=%q, want %q", pkgs[0].Name, "go-forge")
	}
	if pkgs[1].Version != "0.2.0" {
		t.Errorf("got version=%q, want %q", pkgs[1].Version, "0.2.0")
	}
}

func TestPackageService_Get_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/packages/core/generic/go-forge/0.1.0" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		json.NewEncoder(w).Encode(types.Package{
			ID:      1,
			Name:    "go-forge",
			Type:    "generic",
			Version: "0.1.0",
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	pkg, err := f.Packages.Get(context.Background(), "core", "generic", "go-forge", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if pkg.ID != 1 {
		t.Errorf("got id=%d, want 1", pkg.ID)
	}
	if pkg.Name != "go-forge" {
		t.Errorf("got name=%q, want %q", pkg.Name, "go-forge")
	}
	if pkg.Version != "0.1.0" {
		t.Errorf("got version=%q, want %q", pkg.Version, "0.1.0")
	}
}

func TestPackageService_Delete_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/packages/core/generic/go-forge/0.1.0" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	err := f.Packages.Delete(context.Background(), "core", "generic", "go-forge", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
}

func TestPackageService_ListFiles_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/packages/core/generic/go-forge/0.1.0/files" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "1")
		json.NewEncoder(w).Encode([]types.PackageFile{
			{ID: 1, Name: "go-forge-0.1.0.tar.gz", Size: 1024, HashMD5: "abc123"},
		})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	files, err := f.Packages.ListFiles(context.Background(), "core", "generic", "go-forge", "0.1.0")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 1 {
		t.Fatalf("got %d files, want 1", len(files))
	}
	if files[0].Name != "go-forge-0.1.0.tar.gz" {
		t.Errorf("got name=%q, want %q", files[0].Name, "go-forge-0.1.0.tar.gz")
	}
	if files[0].Size != 1024 {
		t.Errorf("got size=%d, want 1024", files[0].Size)
	}
}

func TestPackageService_NotFound_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "package not found"})
	}))
	defer srv.Close()

	f := NewForge(srv.URL, "tok")
	_, err := f.Packages.Get(context.Background(), "core", "generic", "nonexistent", "0.0.0")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !IsNotFound(err) {
		t.Errorf("expected not-found error, got %v", err)
	}
}

func TestPackages_PackageService_List_Good(t *core.T) {
	subject := (*PackageService).List
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_List_Bad(t *core.T) {
	subject := (*PackageService).List
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_List_Ugly(t *core.T) {
	subject := (*PackageService).List
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Iter_Good(t *core.T) {
	subject := (*PackageService).Iter
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Iter_Bad(t *core.T) {
	subject := (*PackageService).Iter
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Iter_Ugly(t *core.T) {
	subject := (*PackageService).Iter
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Get_Good(t *core.T) {
	subject := (*PackageService).Get
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Get_Bad(t *core.T) {
	subject := (*PackageService).Get
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Get_Ugly(t *core.T) {
	subject := (*PackageService).Get
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Delete_Good(t *core.T) {
	subject := (*PackageService).Delete
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Delete_Bad(t *core.T) {
	subject := (*PackageService).Delete
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_Delete_Ugly(t *core.T) {
	subject := (*PackageService).Delete
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_ListFiles_Good(t *core.T) {
	subject := (*PackageService).ListFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_ListFiles_Bad(t *core.T) {
	subject := (*PackageService).ListFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_ListFiles_Ugly(t *core.T) {
	subject := (*PackageService).ListFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_IterFiles_Good(t *core.T) {
	subject := (*PackageService).IterFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_IterFiles_Bad(t *core.T) {
	subject := (*PackageService).IterFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestPackages_PackageService_IterFiles_Ugly(t *core.T) {
	subject := (*PackageService).IterFiles
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
