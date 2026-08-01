package forge

import (
	"os"
	"path/filepath"
	"testing"
)

// NewForgeFromConfig wraps ResolveConfig + NewForge with a Result contract and
// a panic-recovering deferred guard. These tests reach its error branches: a
// broken config file with no fallback token, and the no-token failure with a
// resolved URL.

func TestNewForgeFromConfig_BrokenConfig_Bad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("FORGE_URL", "")
	t.Setenv("FORGE_TOKEN", "")

	cfgPath := filepath.Join(home, ".config", "forge", "config.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, []byte("{not-json"), 0600); err != nil {
		t.Fatal(err)
	}

	// No flags, no env, broken file → ResolveConfig surfaces a decode error
	// which NewForgeFromConfig must propagate as a failed Result.
	r := NewForgeFromConfig("", "")
	if r.OK {
		t.Fatal("expected a failed Result from a broken config file")
	}
}

func TestNewForgeFromConfig_URLNoToken_Bad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "https://forge.example.com")
	t.Setenv("FORGE_TOKEN", "")

	// URL resolves but no token → the explicit no-token failure path.
	r := NewForgeFromConfig("", "")
	if r.OK {
		t.Fatal("expected a failed Result when no token is configured")
	}
}

func TestReadConfigFile_Decode_Bad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	cfgPath := filepath.Join(home, ".config", "forge", "config.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, []byte("{not-json"), 0600); err != nil {
		t.Fatal(err)
	}

	if _, _, err := readConfigFile(); err == nil {
		t.Fatal("expected a decode error from a malformed config file")
	}
}

func TestReadConfigFile_Missing_Good(t *testing.T) {
	// A missing config file is not an error — empty values are returned.
	t.Setenv("HOME", t.TempDir())
	url, token, err := readConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	if url != "" || token != "" {
		t.Fatalf("missing file should yield empty values, got %q %q", url, token)
	}
}

func TestSaveConfig_MkdirFail_Bad(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	// Plant a regular file where the .config directory needs to be, so the
	// MkdirAll for the config parent fails.
	if err := os.WriteFile(filepath.Join(home, ".config"), []byte("x"), 0600); err != nil {
		t.Fatal(err)
	}

	if err := SaveConfig("https://forge.example.com", "tok"); err == nil {
		t.Fatal("expected SaveConfig to fail when the parent path is a file")
	}
}

func TestSaveConfig_RoundTrip_Good(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SaveConfig("https://forge.example.com", "round-trip-token"); err != nil {
		t.Fatal(err)
	}
	url, token, err := readConfigFile()
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://forge.example.com" || token != "round-trip-token" {
		t.Fatalf("round-trip mismatch: got %q %q", url, token)
	}
}
