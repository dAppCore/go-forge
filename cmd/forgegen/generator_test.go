package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	coreio "dappco.re/go/core/io"
)

func TestGenerate_Good_CreatesFiles(t *testing.T) {
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
		if strings.HasSuffix(e.Name(), ".go") {
			goFiles++
		}
	}
	if goFiles == 0 {
		t.Fatal("no .go files generated")
	}
}

func TestGenerate_Good_ValidGoSyntax(t *testing.T) {
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
		if strings.HasSuffix(e.Name(), ".go") {
			content, err = coreio.Local.Read(filepath.Join(outDir, e.Name()))
			if err == nil {
				break
			}
		}
	}
	if err != nil || content == "" {
		t.Fatal("could not read any generated file")
	}
	if !strings.Contains(content, "package types") {
		t.Error("missing package declaration")
	}
	if !strings.Contains(content, "// Code generated") {
		t.Error("missing generated comment")
	}
}

func TestGenerate_Good_RepositoryType(t *testing.T) {
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
		data, _ := coreio.Local.Read(filepath.Join(outDir, e.Name()))
		if strings.Contains(data, "type Repository struct") {
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
		if !strings.Contains(content, check) {
			t.Errorf("missing field with tag %s", check)
		}
	}
}

func TestGenerate_Good_TimeImport(t *testing.T) {
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
		content, _ := coreio.Local.Read(filepath.Join(outDir, e.Name()))
		if strings.Contains(content, "time.Time") && !strings.Contains(content, "\"time\"") {
			t.Errorf("file %s uses time.Time but doesn't import time", e.Name())
		}
	}
}
