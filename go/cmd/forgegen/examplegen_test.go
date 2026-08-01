package main

import (
	"testing"
)

// The usage-example generators are pure functions over GoType / GoField.
// These tests pin every branch of exampleTypeExpression, exampleValue,
// exampleStringValue and usageExample.

func TestExampleTypeExpression_AllBranches_Good(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"string", `"example"`},
		{"bool", "true"},
		{"int", "1"},
		{"int32", "1"},
		{"int64", "1"},
		{"uint", "1"},
		{"uint32", "1"},
		{"uint64", "1"},
		{"float32", "1.0"},
		{"float64", "1.0"},
		{"time.Time", "time.Now()"},
		{"[]string", `[]string{"example"}`},
		{"[]int64", "[]int64{1}"},
		{"[]int", "[]int{1}"},
		{"map[string]string", `map[string]string{"key": "value"}`},
		{"SomeStruct", "SomeStruct{}"},
	}
	for _, tc := range cases {
		if got := exampleTypeExpression(tc.in); got != tc.want {
			t.Errorf("exampleTypeExpression(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestExampleValue_AllBranches_Good(t *testing.T) {
	cases := []struct {
		name  string
		field GoField
		want  string
	}{
		{"pointer", GoField{GoName: "Owner", GoType: "*User"}, "&User{}"},
		{"string", GoField{GoName: "Name", GoType: "string"}, `"example"`},
		{"time", GoField{GoName: "Created", GoType: "time.Time"}, "time.Now()"},
		{"bool", GoField{GoName: "Active", GoType: "bool"}, "true"},
		{"int64", GoField{GoName: "ID", GoType: "int64"}, "1"},
		{"int", GoField{GoName: "Count", GoType: "int"}, "1"},
		{"slice string", GoField{GoName: "Labels", GoType: "[]string"}, `[]string{"example"}`},
		// NOTE: the int64/int suffix cases precede the []int64/[]int prefix
		// cases, so a []int64 GoType matches the suffix branch first and
		// resolves to "1" — the dedicated slice-int branches are effectively
		// shadowed. Pinned as production behaviour, not the intuitive result.
		{"slice int64 (suffix shadow)", GoField{GoName: "IDs", GoType: "[]int64"}, "1"},
		{"slice int (suffix shadow)", GoField{GoName: "Nums", GoType: "[]int"}, "1"},
		{"map", GoField{GoName: "Meta", GoType: "map[string]string"}, `map[string]string{"key": "value"}`},
		{"unknown", GoField{GoName: "Weird", GoType: "rune"}, "{}"},
	}
	for _, tc := range cases {
		if got := exampleValue(tc.field); got != tc.want {
			t.Errorf("%s: exampleValue = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestExampleStringValue_AllBranches_Good(t *testing.T) {
	cases := []struct {
		name string
		want string
	}{
		{"AvatarURL", `"https://example.com"`},
		{"Email", `"alice@example.com"`},
		{"TagName", `"v1.0.0"`},
		{"Branch", `"main"`},
		{"Ref", `"main"`},
		{"Name", `"example"`},
	}
	for _, tc := range cases {
		if got := exampleStringValue(tc.name); got != tc.want {
			t.Errorf("exampleStringValue(%q) = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestUsageExample_AllBranches_Good(t *testing.T) {
	enum := &GoType{Name: "StateType", IsEnum: true, EnumValues: []string{"open", "closed"}}
	if got := usageExample(enum); got != "StateTypeOpen" {
		t.Errorf("enum: got %q, want StateTypeOpen", got)
	}

	alias := &GoType{Name: "RepoID", IsAlias: true, AliasType: "int64"}
	if got := usageExample(alias); got != "RepoID(1)" {
		t.Errorf("alias: got %q, want RepoID(1)", got)
	}

	withField := &GoType{Name: "CreateRepoOption", Fields: []GoField{{GoName: "Name", GoType: "string", Required: true}}}
	if got := usageExample(withField); got != `CreateRepoOption{Name: "example"}` {
		t.Errorf("struct: got %q", got)
	}

	empty := &GoType{Name: "EmptyOption"}
	if got := usageExample(empty); got != "EmptyOption{}" {
		t.Errorf("empty: got %q, want EmptyOption{}", got)
	}
}

func TestUsageExample_EnumNoValues_Ugly(t *testing.T) {
	// An enum flag with no values falls through to the default literal path.
	enum := &GoType{Name: "StateType", IsEnum: true}
	if got := usageExample(enum); got != "StateType{}" {
		t.Errorf("got %q, want StateType{}", got)
	}
}
