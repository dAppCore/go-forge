package main

import (
	"testing"
)

func TestParser_LoadSpec_Good(t *testing.T) {
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

func TestParser_ExtractTypes_Good(t *testing.T) {
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

func TestParser_FieldTypes_Good(t *testing.T) {
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
		case "units_map":
			if f.GoType != "map[string]string" {
				t.Errorf("units_map: got %q, want map[string]string", f.GoType)
			}
		}
	}
}

func TestParser_DetectCreateEditPairs_Good(t *testing.T) {
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

func TestParser_AdditionalPropertiesAlias_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)
	alias, ok := types["CreateHookOptionConfig"]
	if !ok {
		t.Fatal("CreateHookOptionConfig type not found")
	}
	if !alias.IsAlias {
		t.Fatal("expected CreateHookOptionConfig to be emitted as an alias")
	}
	if alias.AliasType != "map[string]any" {
		t.Fatalf("got alias type %q, want map[string]any", alias.AliasType)
	}
}

func TestParser_PrimitiveAndCollectionAliases_Good(t *testing.T) {
	spec, err := LoadSpec("../../testdata/swagger.v1.json")
	if err != nil {
		t.Fatal(err)
	}

	types := ExtractTypes(spec)

	cases := []struct {
		name     string
		wantType string
	}{
		{name: "CommitStatusState", wantType: "string"},
		{name: "IssueFormFieldType", wantType: "string"},
		{name: "IssueFormFieldVisible", wantType: "string"},
		{name: "NotifySubjectType", wantType: "string"},
		{name: "ReviewStateType", wantType: "string"},
		{name: "StateType", wantType: "string"},
		{name: "TimeStamp", wantType: "int64"},
		{name: "IssueTemplateLabels", wantType: "[]string"},
		{name: "QuotaGroupList", wantType: "[]*QuotaGroup"},
		{name: "QuotaUsedArtifactList", wantType: "[]*QuotaUsedArtifact"},
		{name: "QuotaUsedAttachmentList", wantType: "[]*QuotaUsedAttachment"},
		{name: "QuotaUsedPackageList", wantType: "[]*QuotaUsedPackage"},
		{name: "CreatePullReviewCommentOptions", wantType: "CreatePullReviewComment"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			gt, ok := types[tc.name]
			if !ok {
				t.Fatalf("type %q not found", tc.name)
			}
			if !gt.IsAlias {
				t.Fatalf("type %q should be emitted as an alias", tc.name)
			}
			if gt.AliasType != tc.wantType {
				t.Fatalf("type %q: got alias %q, want %q", tc.name, gt.AliasType, tc.wantType)
			}
		})
	}
}
