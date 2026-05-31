package types

import (
	"testing"
)

// normaliseLabelRefs coerces a heterogeneous []any of label references (the
// Forgejo API accepts either label names or numeric ids) into a homogeneous
// []string or []int64, passing non-slice or unrecognised shapes through
// unchanged. These tests pin every branch.

func TestNormaliseLabelRefs_NonSlice_Passthrough(t *testing.T) {
	// A non-slice value is returned unchanged.
	in := "not a slice"
	if got := normaliseLabelRefs(in); got != in {
		t.Fatalf("got %v, want passthrough", got)
	}
}

func TestNormaliseLabelRefs_Empty_Good(t *testing.T) {
	got, ok := normaliseLabelRefs([]any{}).([]string)
	if !ok || len(got) != 0 {
		t.Fatalf("empty slice should become []string{}, got %#v", got)
	}
}

func TestNormaliseLabelRefs_Strings_Good(t *testing.T) {
	got, ok := normaliseLabelRefs([]any{"bug", "help"}).([]string)
	if !ok || len(got) != 2 || got[0] != "bug" || got[1] != "help" {
		t.Fatalf("got %#v, want []string{bug, help}", got)
	}
}

func TestNormaliseLabelRefs_Floats_Good(t *testing.T) {
	// JSON numbers decode to float64; whole numbers become int64.
	got, ok := normaliseLabelRefs([]any{float64(1), float64(2)}).([]int64)
	if !ok || len(got) != 2 || got[0] != 1 || got[1] != 2 {
		t.Fatalf("got %#v, want []int64{1, 2}", got)
	}
}

func TestNormaliseLabelRefs_IntKinds_Good(t *testing.T) {
	// Explicit int / int64 elements are also accepted.
	got, ok := normaliseLabelRefs([]any{int(3), int64(4)}).([]int64)
	if !ok || len(got) != 2 || got[0] != 3 || got[1] != 4 {
		t.Fatalf("got %#v, want []int64{3, 4}", got)
	}
}

func TestNormaliseLabelRefs_FractionalFloat_Passthrough(t *testing.T) {
	// A non-integral float cannot be a label id, so the input is returned
	// unchanged for the caller to handle.
	in := []any{float64(1.5)}
	got, ok := normaliseLabelRefs(in).([]any)
	if !ok || len(got) != 1 {
		t.Fatalf("fractional float should pass through unchanged, got %#v", got)
	}
}

func TestNormaliseLabelRefs_UnknownElement_Passthrough(t *testing.T) {
	// A slice mixing non-string, non-numeric elements passes through.
	in := []any{true}
	got, ok := normaliseLabelRefs(in).([]any)
	if !ok || len(got) != 1 {
		t.Fatalf("unknown element type should pass through, got %#v", got)
	}
}

func TestNormaliseLabelRefs_MixedStringThenNumber_Ugly(t *testing.T) {
	// A string followed by a number: the string pass fails on the number, so
	// the numeric pass runs and the string is unrecognised -> passthrough.
	in := []any{"bug", float64(2)}
	got, ok := normaliseLabelRefs(in).([]any)
	if !ok || len(got) != 2 {
		t.Fatalf("mixed string+number should pass through, got %#v", got)
	}
}
