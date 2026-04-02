package main

import (
	"testing"

	core "dappco.re/go/core"
	coreio "dappco.re/go/core/io"
)

func TestGenerate_CreatesFiles_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	outDir := t.TempDir()
	if err := Generate(types, pairs, outDir); err != nil {
		t.Fatal(err)
	}

	entries, _ := coreio.Local.List(outDir)
	goFiles := 0
	for _, e := range entries {
		if core.HasSuffix(e.Name(), ".go") {
			goFiles++
		}
	}
	if goFiles == 0 {
		t.Fatal("no .go files generated")
	}
}

func TestGenerate_ValidGoSyntax_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	outDir := t.TempDir()
	if err := Generate(types, pairs, outDir); err != nil {
		t.Fatal(err)
	}

	entries, _ := coreio.Local.List(outDir)
	var content string
	for _, e := range entries {
		if core.HasSuffix(e.Name(), ".go") {
			content, err = coreio.Local.Read(core.JoinPath(outDir, e.Name()))
			if err == nil {
				break
			}
		}
	}
	if err != nil || content == "" {
		t.Fatal("could not read any generated file")
	}
	if !core.Contains(content, "package types") {
		t.Error("missing package declaration")
	}
	if !core.Contains(content, "// Code generated") {
		t.Error("missing generated comment")
	}
}

func TestGenerate_RepositoryType_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	outDir := t.TempDir()
	if err := Generate(types, pairs, outDir); err != nil {
		t.Fatal(err)
	}

	var content string
	entries, _ := coreio.Local.List(outDir)
	for _, e := range entries {
		data, _ := coreio.Local.Read(core.JoinPath(outDir, e.Name()))
		if core.Contains(data, "type Repository struct") {
			content = data
			break
		}
	}

	if content == "" {
		t.Fatal("Repository type not found in any generated file")
	}

	// Repository has no required fields in the swagger spec,
	// so all fields get the ,omitempty suffix.
	checks := []string{
		"`json:\"id,omitempty\"`",
		"`json:\"name,omitempty\"`",
		"`json:\"full_name,omitempty\"`",
		"`json:\"private,omitempty\"`",
	}
	for _, check := range checks {
		if !core.Contains(content, check) {
			t.Errorf("missing field with tag %s", check)
		}
	}
}

func TestGenerate_TimeImport_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	outDir := t.TempDir()
	if err := Generate(types, pairs, outDir); err != nil {
		t.Fatal(err)
	}

	entries, _ := coreio.Local.List(outDir)
	for _, e := range entries {
		content, _ := coreio.Local.Read(core.JoinPath(outDir, e.Name()))
		if core.Contains(content, "time.Time") && !core.Contains(content, "\"time\"") {
			t.Errorf("file %s uses time.Time but doesn't import time", e.Name())
		}
	}
}

func TestGenerate_AdditionalProperties_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	outDir := t.TempDir()
	if err := Generate(types, pairs, outDir); err != nil {
		t.Fatal(err)
	}

	entries, _ := coreio.Local.List(outDir)
	var hookContent string
	var teamContent string
	for _, e := range entries {
		data, _ := coreio.Local.Read(core.JoinPath(outDir, e.Name()))
		if core.Contains(data, "type CreateHookOptionConfig") {
			hookContent = data
		}
		if core.Contains(data, "UnitsMap map[string]string `json:\"units_map,omitempty\"`") {
			teamContent = data
		}
	}
	if hookContent == "" {
		t.Fatal("CreateHookOptionConfig type not found in any generated file")
	}
	if !core.Contains(hookContent, "type CreateHookOptionConfig map[string]any") {
		t.Fatalf("generated alias not found in file:\n%s", hookContent)
	}
	if teamContent == "" {
		t.Fatal("typed units_map field not found in any generated file")
	}
}

func TestGenerate_UsageExamples_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	pairs := DetectCRUDPairs(spec)

	outDir := t.TempDir()
	if err := Generate(types, pairs, outDir); err != nil {
		t.Fatal(err)
	}

	entries, _ := coreio.Local.List(outDir)
	var content string
	for _, e := range entries {
		data, _ := coreio.Local.Read(core.JoinPath(outDir, e.Name()))
		if core.Contains(data, "type CreateIssueOption struct") {
			content = data
			break
		}
	}
	if content == "" {
		t.Fatal("CreateIssueOption type not found in any generated file")
	}
	if !core.Contains(content, "// Usage:") {
		t.Fatalf("generated option type is missing usage documentation:\n%s", content)
	}
	if !core.Contains(content, "opts :=") {
		t.Fatalf("generated usage example is missing assignment syntax:\n%s", content)
	}
}
