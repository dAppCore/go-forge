package forge

import (
	"fmt"
	"testing"
)

func TestParams_String_Good(t *testing.T) {
	params := Params{"repo": "go-forge", "owner": "core"}
	want := `forge.Params{owner="core", repo="go-forge"}`
	if got := params.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(params); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", params); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestParams_String_NilSafe(t *testing.T) {
	var params Params
	want := "forge.Params{<nil>}"
	if got := params.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(params); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", params); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestListOptions_String_Good(t *testing.T) {
	opts := ListOptions{Page: 2, Limit: 25}
	want := "forge.ListOptions{page=2, limit=25}"
	if got := opts.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(opts); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", opts); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestRateLimit_String_Good(t *testing.T) {
	rl := RateLimit{Limit: 80, Remaining: 79, Reset: 1700000003}
	want := "forge.RateLimit{limit=80, remaining=79, reset=1700000003}"
	if got := rl.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(rl); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", rl); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}

func TestPagedResult_String_Good(t *testing.T) {
	page := PagedResult[int]{
		Items:      []int{1, 2, 3},
		TotalCount: 10,
		Page:       2,
		HasMore:    true,
	}
	want := "forge.PagedResult{items=3, totalCount=10, page=2, hasMore=true}"
	if got := page.String(); got != want {
		t.Fatalf("got String()=%q, want %q", got, want)
	}
	if got := fmt.Sprint(page); got != want {
		t.Fatalf("got fmt.Sprint=%q, want %q", got, want)
	}
	if got := fmt.Sprintf("%#v", page); got != want {
		t.Fatalf("got GoString=%q, want %q", got, want)
	}
}
