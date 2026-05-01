package forge

import (
	"fmt"
	"testing"
)

func TestResource_String_Good(t *testing.T) {
	res := NewResource[int, struct{}, struct{}](NewClient("https://forge.lthn.ai", "tok"), "/api/v1/repos/{owner}/{repo}")
	got := fmt.Sprint(res)
	want := `forge.Resource{path="/api/v1/repos/{owner}/{repo}", collection="/api/v1/repos/{owner}"}`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := res.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", res); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}
