package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	json "github.com/goccy/go-json"
)

// AX-10: the CLI artifact is validated by building and running the driver
// binary, mirroring the Taskfile `test` path (go build ./tests/cli/forge then
// run it). These tests cover the skip path (no credentials), the success path
// (against a mock Forgejo server) and the error path (server 500).

// buildDriver compiles the CLI driver once per test into a temp binary.
func buildDriver(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "core-forge-driver")
	cmd := exec.Command("go", "build", "-o", bin, ".")
	cmd.Env = os.Environ()
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build driver: %v\n%s", err, out)
	}
	return bin
}

func TestCLIDriver_Skip_Good(t *testing.T) {
	bin := buildDriver(t)
	cmd := exec.Command(bin)
	// Strip any inherited FORGE_* so the skip path is taken.
	cmd.Env = append(os.Environ(), "FORGE_URL=", "FORGE_TOKEN=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("driver exited non-zero: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "skip:") {
		t.Fatalf("expected skip message, got %q", out)
	}
}

func TestCLIDriver_ListRepos_Good(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/orgs/core/repos" {
			t.Errorf("wrong path: %s", r.URL.Path)
		}
		w.Header().Set("X-Total-Count", "2")
		json.NewEncoder(w).Encode([]map[string]any{{"name": "go-forge"}, {"name": "go-scm"}})
	}))
	defer srv.Close()

	bin := buildDriver(t)
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "FORGE_URL="+srv.URL, "FORGE_TOKEN=test-token")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("driver exited non-zero: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "core repos: 2") {
		t.Fatalf("expected repo count, got %q", out)
	}
}

func TestCLIDriver_ListReposError_Bad(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": "boom"})
	}))
	defer srv.Close()

	bin := buildDriver(t)
	cmd := exec.Command(bin)
	cmd.Env = append(os.Environ(), "FORGE_URL="+srv.URL, "FORGE_TOKEN=test-token")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected non-zero exit on server error, got output %q", out)
	}
	if !strings.Contains(string(out), "list core repos") {
		t.Fatalf("expected error message, got %q", out)
	}
}
