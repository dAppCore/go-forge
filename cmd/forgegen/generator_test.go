package main

import (
	"os"
	"testing"

	core "dappco.re/go"
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

	entries, _ := os.ReadDir(outDir)
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

	entries, _ := os.ReadDir(outDir)
	var content string
	for _, e := range entries {
		if core.HasSuffix(e.Name(), ".go") {
			data, readErr := os.ReadFile(core.JoinPath(outDir, e.Name()))
			if readErr == nil {
				content = string(data)
				err = nil
				break
			}
			err = readErr
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
	entries, _ := os.ReadDir(outDir)
	for _, e := range entries {
		data, _ := os.ReadFile(core.JoinPath(outDir, e.Name()))
		if core.Contains(string(data), "type Repository struct") {
			content = string(data)
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

	entries, _ := os.ReadDir(outDir)
	for _, e := range entries {
		data, _ := os.ReadFile(core.JoinPath(outDir, e.Name()))
		content := string(data)
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

	entries, _ := os.ReadDir(outDir)
	var hookContent string
	var teamContent string
	for _, e := range entries {
		data, _ := os.ReadFile(core.JoinPath(outDir, e.Name()))
		content := string(data)
		if core.Contains(content, "type CreateHookOptionConfig") {
			hookContent = content
		}
		if core.Contains(content, "UnitsMap map[string]string `json:\"units_map,omitempty\"`") {
			teamContent = content
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

	entries, _ := os.ReadDir(outDir)
	var content string
	for _, e := range entries {
		data, _ := os.ReadFile(core.JoinPath(outDir, e.Name()))
		if core.Contains(string(data), "type CreateIssueOption struct") {
			content = string(data)
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

func TestGenerate_UsageExamples_AllKinds_Good(t *testing.T) {
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

	entries, _ := os.ReadDir(outDir)
	var content string
	for _, e := range entries {
		data, _ := os.ReadFile(core.JoinPath(outDir, e.Name()))
		if core.Contains(string(data), "type CommitStatusState string") {
			content = string(data)
			break
		}
	}
	if content == "" {
		t.Fatal("CommitStatusState type not found in any generated file")
	}
	if !core.Contains(content, "type CommitStatusState string") {
		t.Fatalf("CommitStatusState type not generated:\n%s", content)
	}
	if !core.Contains(content, "// Usage:") {
		t.Fatalf("generated enum type is missing usage documentation:\n%s", content)
	}

	content = ""
	for _, e := range entries {
		data, _ := os.ReadFile(core.JoinPath(outDir, e.Name()))
		if core.Contains(string(data), "type CreateHookOptionConfig map[string]any") {
			content = string(data)
			break
		}
	}
	if content == "" {
		t.Fatal("CreateHookOptionConfig type not found in any generated file")
	}
	if !core.Contains(content, "CreateHookOptionConfig(map[string]any{\"key\": \"value\"})") {
		t.Fatalf("generated alias type is missing a valid usage example:\n%s", content)
	}
}
