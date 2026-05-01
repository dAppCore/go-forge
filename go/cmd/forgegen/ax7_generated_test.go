package main

import (
	"os"
	"path/filepath"

	. "dappco.re/go"
)

func TestAX7_Generate_Good(t *T) {
	types := map[string]*GoType{"Repository": {Name: "Repository", Fields: []GoField{{GoName: "Name", GoType: "string", JSONName: "name"}}}}
	err := Generate(types, nil, t.TempDir())
	AssertNoError(t, err)
}

func TestAX7_Generate_Bad(t *T) {
	file := filepath.Join(t.TempDir(), "not-dir")
	AssertNoError(t, os.WriteFile(file, []byte("x"), 0600))
	err := Generate(map[string]*GoType{"Repository": {Name: "Repository"}}, nil, file)
	AssertError(t, err)
}

func TestAX7_Generate_Ugly(t *T) {
	dir := t.TempDir()
	err := Generate(map[string]*GoType{}, nil, dir)
	AssertNoError(t, err)
	_, statErr := os.Stat(dir)
	AssertNoError(t, statErr)
}

func TestAX7_LoadSpec_Good(t *T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	AssertNoError(t, err)
	AssertNotNil(t, spec)
	AssertNotEmpty(t, spec.Definitions)
}

func TestAX7_LoadSpec_Bad(t *T) {
	spec, err := LoadSpec(filepath.Join(t.TempDir(), "missing.json"))
	AssertError(t, err)
	AssertNil(t, spec)
}

func TestAX7_LoadSpec_Ugly(t *T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	AssertNoError(t, os.WriteFile(path, []byte("{bad"), 0600))
	spec, err := LoadSpec(path)
	AssertError(t, err)
	AssertNil(t, spec)
}

func TestAX7_ExtractTypes_Good(t *T) {
	spec := &Spec{Definitions: map[string]SchemaDefinition{"Repository": {Type: "object", Properties: map[string]SchemaProperty{"name": {Type: "string"}}}}}
	got := ExtractTypes(spec)
	AssertNotNil(t, got["Repository"])
	AssertEqual(t, "Repository", got["Repository"].Name)
}

func TestAX7_ExtractTypes_Bad(t *T) {
	spec := &Spec{Definitions: map[string]SchemaDefinition{}}
	got := ExtractTypes(spec)
	AssertEmpty(t, got)
	AssertNotNil(t, got)
}

func TestAX7_ExtractTypes_Ugly(t *T) {
	AssertPanics(t, func() {
		_ = ExtractTypes(nil)
	})
	AssertTrue(t, true)
}

func TestAX7_DetectCRUDPairs_Good(t *T) {
	spec := &Spec{Definitions: map[string]SchemaDefinition{"CreateRepoOption": {}, "EditRepoOption": {}}}
	got := DetectCRUDPairs(spec)
	AssertLen(t, got, 1)
	AssertEqual(t, "Repo", got[0].Base)
}

func TestAX7_DetectCRUDPairs_Bad(t *T) {
	spec := &Spec{Definitions: map[string]SchemaDefinition{"Repository": {}}}
	got := DetectCRUDPairs(spec)
	AssertEmpty(t, got)
	AssertEqual(t, 0, len(got))
}

func TestAX7_DetectCRUDPairs_Ugly(t *T) {
	AssertPanics(t, func() {
		_ = DetectCRUDPairs(nil)
	})
	AssertTrue(t, true)
}
