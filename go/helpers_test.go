package forge

import (
	"testing"
	"time"
)

// The option-string helpers underpin every option type's String()/GoString().
// These tests pin the per-type behaviour of isZeroOptionValue and
// formatOptionValue, plus the query-assembly helpers.

func TestIsZeroOptionValue_AllTypes_Good(t *testing.T) {
	yes := true
	day := time.Date(2026, 5, 31, 0, 0, 0, 0, time.UTC)
	cases := []struct {
		name string
		v    any
		zero bool
	}{
		{"nil", nil, true},
		{"empty string", "", true},
		{"set string", "x", false},
		{"false bool", false, true},
		{"true bool", true, false},
		{"zero int", 0, true},
		{"set int", 1, false},
		{"zero int64", int64(0), true},
		{"set int64", int64(2), false},
		{"empty slice", []string{}, true},
		{"set slice", []string{"a"}, false},
		{"nil *time.Time", (*time.Time)(nil), true},
		{"set *time.Time", &day, false},
		{"nil *bool", (*bool)(nil), true},
		{"set *bool", &yes, false},
		{"zero time.Time", time.Time{}, true},
		{"set time.Time", day, false},
		{"unknown type", 3.14, false},
	}
	for _, tc := range cases {
		if got := isZeroOptionValue(tc.v); got != tc.zero {
			t.Errorf("%s: isZeroOptionValue = %v, want %v", tc.name, got, tc.zero)
		}
	}
}

func TestFormatOptionValue_AllTypes_Good(t *testing.T) {
	yes := true
	day := time.Date(2026, 5, 31, 1, 2, 3, 0, time.UTC)
	cases := []struct {
		name string
		v    any
		want string
	}{
		{"string", "hi", `"hi"`},
		{"bool", true, "true"},
		{"int", 7, "7"},
		{"int64", int64(8), "8"},
		{"slice", []string{"a", "b"}, `[]string{"a", "b"}`},
		{"set *time.Time", &day, `"` + day.Format(time.RFC3339) + `"`},
		{"set *bool", &yes, "true"},
		{"time.Time", day, `"` + day.Format(time.RFC3339) + `"`},
	}
	for _, tc := range cases {
		if got := formatOptionValue(tc.v); got != tc.want {
			t.Errorf("%s: formatOptionValue = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestFormatOptionValue_NilPointers_Ugly(t *testing.T) {
	if got := formatOptionValue((*time.Time)(nil)); got != "<nil>" {
		t.Errorf("nil *time.Time: got %q, want <nil>", got)
	}
	if got := formatOptionValue((*bool)(nil)); got != "<nil>" {
		t.Errorf("nil *bool: got %q, want <nil>", got)
	}
}

func TestFormatOptionValue_UnknownType_Ugly(t *testing.T) {
	// The default branch falls back to a Go-syntax representation.
	if got := formatOptionValue(3.5); got == "" {
		t.Error("default branch should produce a non-empty representation")
	}
}

func TestOptionString_SkipsZeroFields_Good(t *testing.T) {
	got := optionString("forge.Demo", "a", "set", "b", "", "c", 0, "d", int64(9))
	want := `forge.Demo{a="set", d=9}`
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAppendQuery_NoQuery_Good(t *testing.T) {
	// A builder that sets nothing must leave the path untouched.
	got := appendQuery("/api/v1/x", func(_ *queryBuilder) {})
	if got != "/api/v1/x" {
		t.Fatalf("got %q, want unchanged path", got)
	}
}

func TestAppendQuery_FreshQuery_Good(t *testing.T) {
	got := appendQuery("/api/v1/x", func(q *queryBuilder) { q.Set("k", "v") })
	if got != "/api/v1/x?k=v" {
		t.Fatalf("got %q, want /api/v1/x?k=v", got)
	}
}

func TestAppendQuery_ExistingQuery_Ugly(t *testing.T) {
	// When the path already carries a query string, the new params are joined
	// with '&' rather than a second '?'.
	got := appendQuery("/api/v1/x?a=1", func(q *queryBuilder) { q.Set("b", "2") })
	if got != "/api/v1/x?a=1&b=2" {
		t.Fatalf("got %q, want /api/v1/x?a=1&b=2", got)
	}
}

func TestQueryBuilder_AddMulti_Good(t *testing.T) {
	q := newQueryBuilder()
	q.Add("k", "1")
	q.Add("k", "2")
	if got := q.Encode(); got != "k=1&k=2" {
		t.Fatalf("got %q, want k=1&k=2", got)
	}
}

func TestQueryBuilder_Encode_EmptyAndNil_Bad(t *testing.T) {
	if got := newQueryBuilder().Encode(); got != "" {
		t.Errorf("empty builder: got %q, want empty", got)
	}
	var nilQ *queryBuilder
	if got := nilQ.Encode(); got != "" {
		t.Errorf("nil builder: got %q, want empty", got)
	}
}

func TestTrimTrailingSlashes_Good(t *testing.T) {
	if got := trimTrailingSlashes("http://x///"); got != "http://x" {
		t.Fatalf("got %q, want http://x", got)
	}
	if got := trimTrailingSlashes("http://x"); got != "http://x" {
		t.Fatalf("got %q, want http://x", got)
	}
}

func TestLastIndexByte_Good(t *testing.T) {
	if got := lastIndexByte("a.b.c", '.'); got != 3 {
		t.Fatalf("got %d, want 3", got)
	}
}

func TestLastIndexByte_NotFound_Bad(t *testing.T) {
	if got := lastIndexByte("abc", '.'); got != -1 {
		t.Fatalf("got %d, want -1", got)
	}
}

func TestPathParams_OddTrailing_Ugly(t *testing.T) {
	// A trailing key with no value is dropped rather than panicking.
	p := pathParams("owner", "core", "repo")
	if p["owner"] != "core" {
		t.Fatalf("got %v", p)
	}
	if _, ok := p["repo"]; ok {
		t.Fatalf("dangling key should be dropped, got %v", p)
	}
}

func TestParseInt_Bad(t *testing.T) {
	if got := parseInt("nope"); got != 0 {
		t.Fatalf("got %d, want 0", got)
	}
}
