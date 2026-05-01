package forge

import (
	"fmt"
	"testing"
)

func TestClient_String_NilSafe(t *testing.T) {
	var c *Client
	want := "forge.Client{<nil>}"
	if got := c.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(c); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", c); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestForge_String_NilSafe(t *testing.T) {
	var f *Forge
	want := "forge.Forge{<nil>}"
	if got := f.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(f); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", f); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestResource_String_NilSafe(t *testing.T) {
	var r *Resource[int, struct{}, struct{}]
	want := "forge.Resource{<nil>}"
	if got := r.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(r); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", r); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestAPIError_String_NilSafe(t *testing.T) {
	var e *APIError
	want := "forge.APIError{<nil>}"
	if got := e.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(e); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", e); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}
