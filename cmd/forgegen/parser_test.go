package main

import (
	"testing"
)

func TestParser_Good_LoadSpec(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}
	if spec.Swagger != "2.0" {
		t.Errorf("got swagger=%q", spec.Swagger)
	}
	if len(spec.Definitions) < 200 {
		t.Errorf("got %d definitions, expected 200+", len(spec.Definitions))
	}
}

func TestParser_Good_ExtractTypes(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	if len(types) < 200 {
		t.Errorf("got %d types", len(types))
	}

	// Check a known type
	repo, ok := types["Repository"]
	if !ok {
		t.Fatal("Repository type not found")
	}
	if len(repo.Fields) < 50 {
		t.Errorf("Repository has %d fields, expected 50+", len(repo.Fields))
	}
}

func TestParser_Good_FieldTypes(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	repo := types["Repository"]

	// Check specific field mappings
	for _, f := range repo.Fields {
		switch f.JSONName {
		case "id":
			if f.GoType != "int64" {
				t.Errorf("id: got %q, want int64", f.GoType)
			}
		case "name":
			if f.GoType != "string" {
				t.Errorf("name: got %q, want string", f.GoType)
			}
		case "private":
			if f.GoType != "bool" {
				t.Errorf("private: got %q, want bool", f.GoType)
			}
		case "created_at":
			if f.GoType != "time.Time" {
				t.Errorf("created_at: got %q, want time.Time", f.GoType)
			}
		case "owner":
			if f.GoType != "*User" {
				t.Errorf("owner: got %q, want *User", f.GoType)
			}
		}
	}
}

func TestParser_Good_DetectCreateEditPairs(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	pairs := DetectCRUDPairs(spec)
	if len(pairs) < 10 {
		t.Errorf("got %d pairs, expected 10+", len(pairs))
	}

	found := false
	for _, p := range pairs {
		if p.Base == "Repo" {
			found = true
			if p.Create != "CreateRepoOption" {
				t.Errorf("repo create=%q", p.Create)
			}
			if p.Edit != "EditRepoOption" {
				t.Errorf("repo edit=%q", p.Edit)
			}
		}
	}
	if !found {
		t.Fatal("Repo pair not found")
	}
}
