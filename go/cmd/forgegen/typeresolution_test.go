package main

import (
	"testing"
)

// The swagger-to-Go type resolution helpers are pure functions. These tests
// pin every type/format branch of resolveGoType, definitionAliasType,
// definitionNeedsPointer, refGoType, pascalCase and enumConstName.

func TestResolveGoType_AllBranches_Good(t *testing.T) {
	defs := map[string]SchemaDefinition{}
	cases := []struct {
		name string
		prop SchemaProperty
		want string
	}{
		{"string", SchemaProperty{Type: "string"}, "string"},
		{"date-time", SchemaProperty{Type: "string", Format: "date-time"}, "time.Time"},
		{"binary", SchemaProperty{Type: "string", Format: "binary"}, "[]byte"},
		{"int", SchemaProperty{Type: "integer"}, "int"},
		{"int64", SchemaProperty{Type: "integer", Format: "int64"}, "int64"},
		{"int32", SchemaProperty{Type: "integer", Format: "int32"}, "int32"},
		{"float64", SchemaProperty{Type: "number"}, "float64"},
		{"float32", SchemaProperty{Type: "number", Format: "float"}, "float32"},
		{"bool", SchemaProperty{Type: "boolean"}, "bool"},
		{"array of strings", SchemaProperty{Type: "array", Items: &SchemaProperty{Type: "string"}}, "[]string"},
		{"array no items", SchemaProperty{Type: "array"}, "[]any"},
		{"object map", SchemaProperty{Type: "object", AdditionalProperties: &SchemaProperty{Type: "integer", Format: "int64"}}, "map[string]int64"},
		{"object bare", SchemaProperty{Type: "object"}, "map[string]any"},
		{"unknown", SchemaProperty{Type: "weird"}, "any"},
	}
	for _, tc := range cases {
		if got := resolveGoType(tc.prop, defs); got != tc.want {
			t.Errorf("%s: resolveGoType = %q, want %q", tc.name, got, tc.want)
		}
	}
}

func TestResolveGoType_Ref_Good(t *testing.T) {
	defs := map[string]SchemaDefinition{
		"Repository": {Type: "object", Properties: map[string]SchemaProperty{}},
		"State":      {Type: "string", Enum: []any{"open", "closed"}},
	}
	// An object ref needs a pointer; an enum (string) ref does not.
	if got := resolveGoType(SchemaProperty{Ref: "#/definitions/Repository"}, defs); got != "*Repository" {
		t.Errorf("object ref: got %q, want *Repository", got)
	}
	if got := resolveGoType(SchemaProperty{Ref: "#/definitions/State"}, defs); got != "State" {
		t.Errorf("enum ref: got %q, want State", got)
	}
	// An unknown ref defaults to a pointer to the named type.
	if got := resolveGoType(SchemaProperty{Ref: "#/definitions/Missing"}, defs); got != "*Missing" {
		t.Errorf("unknown ref: got %q, want *Missing", got)
	}
}

func TestDefinitionAliasType_AllBranches_Good(t *testing.T) {
	defs := map[string]SchemaDefinition{}
	cases := []struct {
		name string
		def  SchemaDefinition
		want string
		ok   bool
	}{
		{"ref", SchemaDefinition{Ref: "#/definitions/Foo"}, "Foo", true},
		{"string", SchemaDefinition{Type: "string"}, "string", true},
		{"int", SchemaDefinition{Type: "integer"}, "int", true},
		{"int64", SchemaDefinition{Type: "integer", Format: "int64"}, "int64", true},
		{"int32", SchemaDefinition{Type: "integer", Format: "int32"}, "int32", true},
		{"float64", SchemaDefinition{Type: "number"}, "float64", true},
		{"float32", SchemaDefinition{Type: "number", Format: "float"}, "float32", true},
		{"bool", SchemaDefinition{Type: "boolean"}, "bool", true},
		{"array items", SchemaDefinition{Type: "array", Items: &SchemaProperty{Type: "string"}}, "[]string", true},
		{"array no items", SchemaDefinition{Type: "array"}, "[]any", true},
		// NOTE: definitionAliasType passes *def.AdditionalProperties to
		// resolveMapType, which then reads THAT schema's own
		// AdditionalProperties for the value type. So a single level of
		// additionalProperties yields map[string]any; the value type is only
		// honoured when additionalProperties is itself nested.
		{"object map single-level", SchemaDefinition{Type: "object", AdditionalProperties: &SchemaProperty{Type: "string"}}, "map[string]any", true},
		{"object map nested", SchemaDefinition{Type: "object", AdditionalProperties: &SchemaProperty{Type: "object", AdditionalProperties: &SchemaProperty{Type: "string"}}}, "map[string]string", true},
		{"object bare", SchemaDefinition{Type: "object"}, "", false},
		{"unknown", SchemaDefinition{Type: "weird"}, "", false},
	}
	for _, tc := range cases {
		got, ok := definitionAliasType(tc.def, defs)
		if got != tc.want || ok != tc.ok {
			t.Errorf("%s: definitionAliasType = (%q, %v), want (%q, %v)", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func TestDefinitionNeedsPointer_Good(t *testing.T) {
	cases := []struct {
		name string
		def  SchemaDefinition
		want bool
	}{
		{"enum", SchemaDefinition{Type: "string", Enum: []any{"a"}}, false},
		{"ref", SchemaDefinition{Ref: "#/definitions/Foo"}, false},
		{"string", SchemaDefinition{Type: "string"}, false},
		{"array", SchemaDefinition{Type: "array"}, false},
		{"object", SchemaDefinition{Type: "object"}, true},
		{"unknown", SchemaDefinition{Type: "weird"}, false},
	}
	for _, tc := range cases {
		if got := definitionNeedsPointer(tc.def); got != tc.want {
			t.Errorf("%s: definitionNeedsPointer = %v, want %v", tc.name, got, tc.want)
		}
	}
}

func TestPascalCase_Acronyms_Good(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"created_by", "CreatedBy"},
		{"html_url", "HTMLURL"},
		{"ssh_url", "SSHURL"},
		{"api_token", "APIToken"},
		{"id", "ID"},
		{"repo-name", "RepoName"},
		{"", ""},
		// Consecutive delimiters produce empty segments that must be skipped.
		{"a__b", "AB"},
		{"--x--", "X"},
	}
	for _, tc := range cases {
		if got := pascalCase(tc.in); got != tc.want {
			t.Errorf("pascalCase(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestEnumConstName_Good(t *testing.T) {
	if got := enumConstName("StateType", "in_progress"); got != "StateTypeInProgress" {
		t.Fatalf("got %q, want StateTypeInProgress", got)
	}
}
