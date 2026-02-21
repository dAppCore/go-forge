package forge

import (
	"os"
	"testing"
)

func TestResolveConfig_Good_EnvOverrides(t *testing.T) {
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

func TestResolveConfig_Good_FlagOverridesEnv(t *testing.T) {
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

func TestResolveConfig_Good_DefaultURL(t *testing.T) {
	os.Unsetenv("FORGE_URL")
	os.Unsetenv("FORGE_TOKEN")

	url, _, err := ResolveConfig("", "")
	if err != nil {
		t.Fatal(err)
	}
	if url != DefaultURL {
		t.Errorf("got url=%q, want %q", url, DefaultURL)
	}
}

func TestNewForgeFromConfig_Bad_NoToken(t *testing.T) {
	os.Unsetenv("FORGE_URL")
	os.Unsetenv("FORGE_TOKEN")

	_, err := NewForgeFromConfig("", "")
	if err == nil {
		t.Fatal("expected error for missing token")
	}
}
