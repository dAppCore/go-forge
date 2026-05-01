package forge

import "testing"

func TestResolvePath_Simple_Good(t *testing.T) {
	got := ResolvePath("/api/v1/repos/{owner}/{repo}", Params{"owner": "core", "repo": "go-forge"})
	want := "/api/v1/repos/core/go-forge"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolvePath_NoParams_Good(t *testing.T) {
	got := ResolvePath("/api/v1/user", nil)
	if got != "/api/v1/user" {
		t.Errorf("got %q", got)
	}
}

func TestResolvePath_WithID_Good(t *testing.T) {
	got := ResolvePath("/api/v1/repos/{owner}/{repo}/issues/{index}", Params{
		"owner": "core", "repo": "go-forge", "index": "42",
	})
	want := "/api/v1/repos/core/go-forge/issues/42"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestResolvePath_URLEncoding_Good(t *testing.T) {
	got := ResolvePath("/api/v1/repos/{owner}/{repo}", Params{"owner": "my org", "repo": "my repo"})
	want := "/api/v1/repos/my%20org/my%20repo"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
