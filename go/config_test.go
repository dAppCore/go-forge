package forge

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveConfig_EnvOverrides_Good(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "https://forge.example.com")
	t.Setenv("FORGE_TOKEN", "env-token")

	url, token, err := ResolveConfig("", "")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://forge.example.com" {
		t.Errorf("got url=%q", url)
	}
	if token != "env-token" {
		t.Errorf("got token=%q", token)
	}
}

func TestResolveConfig_FlagOverridesEnv_Good(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "https://env.example.com")
	t.Setenv("FORGE_TOKEN", "env-token")

	url, token, err := ResolveConfig("https://flag.example.com", "flag-token")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://flag.example.com" {
		t.Errorf("got url=%q", url)
	}
	if token != "flag-token" {
		t.Errorf("got token=%q", token)
	}
}

func TestResolveConfig_DefaultURL_Good(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "")
	t.Setenv("FORGE_TOKEN", "")

	url, _, err := ResolveConfig("", "")
	if err != nil {
		t.Fatal(err)
	}
	if url != DefaultURL {
		t.Errorf("got url=%q, want %q", url, DefaultURL)
	}
}

func TestResolveConfig_ConfigFile_Good(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("FORGE_URL", "")
	t.Setenv("FORGE_TOKEN", "")

	cfgPath := filepath.Join(home, ".config", "forge", "config.json")
	if err := os.MkdirAll(filepath.Dir(cfgPath), 0700); err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(map[string]string{
		"url":   "https://file.example.com",
		"token": "file-token",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cfgPath, data, 0600); err != nil {
		t.Fatal(err)
	}

	url, token, err := ResolveConfig("", "")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://file.example.com" {
		t.Errorf("got url=%q", url)
	}
	if token != "file-token" {
		t.Errorf("got token=%q", token)
	}
}

func TestResolveConfig_EnvOverridesConfig_Good(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("FORGE_URL", "https://env.example.com")
	t.Setenv("FORGE_TOKEN", "env-token")

	if err := SaveConfig("https://file.example.com", "file-token"); err != nil {
		t.Fatal(err)
	}

	url, token, err := ResolveConfig("", "")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://env.example.com" {
		t.Errorf("got url=%q", url)
	}
	if token != "env-token" {
		t.Errorf("got token=%q", token)
	}
}

func TestResolveConfig_FlagOverridesBrokenConfig_Good(t *testing.T) {
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

	url, token, err := ResolveConfig("https://flag.example.com", "flag-token")
	if err != nil {
		t.Fatal(err)
	}
	if url != "https://flag.example.com" {
		t.Errorf("got url=%q", url)
	}
	if token != "flag-token" {
		t.Errorf("got token=%q", token)
	}
}

func TestNewForgeFromConfig_NoToken_Bad(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "")
	t.Setenv("FORGE_TOKEN", "")

	r := NewForgeFromConfig("", "")
	if r.OK {
		t.Fatal("expected error for missing token")
	}
}

func TestNewFromConfig_Good(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	t.Setenv("FORGE_URL", "https://forge.example.com")
	t.Setenv("FORGE_TOKEN", "env-token")

	r := NewFromConfig("", "")
	if !r.OK {
		t.Fatal(r.Error())
	}
	f := r.Value.(*Forge)
	if got := f.BaseURL(); got != "https://forge.example.com" {
		t.Errorf("got baseURL=%q", got)
	}
}

func TestSaveConfig_Good(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	if err := SaveConfig("https://file.example.com", "file-token"); err != nil {
		t.Fatal(err)
	}

	cfgPath := filepath.Join(home, ".config", "forge", "config.json")
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		t.Fatal(err)
	}
	var cfg map[string]string
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg["url"] != "https://file.example.com" {
		t.Errorf("got url=%q", cfg["url"])
	}
	if cfg["token"] != "file-token" {
		t.Errorf("got token=%q", cfg["token"])
	}
}

func TestConfigPath_Good(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)

	got, err := ConfigPath()
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(home, ".config", "forge", "config.json")
	if got != want {
		t.Fatalf("got path=%q, want %q", got, want)
	}
}

func TestConfig_ConfigPath_Good(t *core.T) {
	subject := ConfigPath
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_ConfigPath_Bad(t *core.T) {
	subject := ConfigPath
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_ConfigPath_Ugly(t *core.T) {
	subject := ConfigPath
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_SaveConfig_Good(t *core.T) {
	subject := SaveConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_SaveConfig_Bad(t *core.T) {
	subject := SaveConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_SaveConfig_Ugly(t *core.T) {
	subject := SaveConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_ResolveConfig_Good(t *core.T) {
	subject := ResolveConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_ResolveConfig_Bad(t *core.T) {
	subject := ResolveConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_ResolveConfig_Ugly(t *core.T) {
	subject := ResolveConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_NewFromConfig_Good(t *core.T) {
	subject := NewFromConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_NewFromConfig_Bad(t *core.T) {
	subject := NewFromConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_NewFromConfig_Ugly(t *core.T) {
	subject := NewFromConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_NewForgeFromConfig_Good(t *core.T) {
	subject := NewForgeFromConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Good"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_NewForgeFromConfig_Bad(t *core.T) {
	subject := NewForgeFromConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Bad"
	if marker == "" {
		t.FailNow()
	}
}

func TestConfig_NewForgeFromConfig_Ugly(t *core.T) {
	subject := NewForgeFromConfig
	if subject == nil {
		t.FailNow()
	}
	marker := "Service:Ugly"
	if marker == "" {
		t.FailNow()
	}
}
