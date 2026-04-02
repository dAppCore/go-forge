package main

import (
	"cmp"
	json "github.com/goccy/go-json"
	"slices"

	core "dappco.re/go/core"
	coreio "dappco.re/go/core/io"
)

// Spec represents a Swagger 2.0 specification document.
//
// Usage:
//
//	spec, err := LoadSpec("testdata/swagger.v1.json")
//	_ = spec
type Spec struct {
	Swagger     string                      `json:"swagger"`
	Info        SpecInfo                    `json:"info"`
	Definitions map[string]SchemaDefinition `json:"definitions"`
	Paths       map[string]map[string]any   `json:"paths"`
}

// SpecInfo holds metadata about the API specification.
//
// Usage:
//
//	_ = SpecInfo{Title: "Forgejo API", Version: "1.0"}
type SpecInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// SchemaDefinition represents a single type definition in the swagger spec.
//
// Usage:
//
//	_ = SchemaDefinition{Type: "object"}
type SchemaDefinition struct {
	Description          string                    `json:"description"`
	Type                 string                    `json:"type"`
	Properties           map[string]SchemaProperty `json:"properties"`
	Required             []string                  `json:"required"`
	Enum                 []any                     `json:"enum"`
	AdditionalProperties *SchemaProperty           `json:"additionalProperties"`
	XGoName              string                    `json:"x-go-name"`
}

// SchemaProperty represents a single property within a schema definition.
//
// Usage:
//
//	_ = SchemaProperty{Type: "string"}
type SchemaProperty struct {
	Type                 string          `json:"type"`
	Format               string          `json:"format"`
	Description          string          `json:"description"`
	Ref                  string          `json:"$ref"`
	Items                *SchemaProperty `json:"items"`
	Enum                 []any           `json:"enum"`
	AdditionalProperties *SchemaProperty `json:"additionalProperties"`
	XGoName              string          `json:"x-go-name"`
}

// GoType is the intermediate representation for a Go type to be generated.
//
// Usage:
//
//	_ = GoType{Name: "Repository"}
type GoType struct {
	Name        string
	Description string
	Usage       string
	Fields      []GoField
	IsEnum      bool
	EnumValues  []string
	IsAlias     bool
	AliasType   string
}

// GoField is the intermediate representation for a single struct field.
//
// Usage:
//
//	_ = GoField{GoName: "ID", GoType: "int64"}
type GoField struct {
	GoName   string
	GoType   string
	JSONName string
	Comment  string
	Required bool
}

// CRUDPair groups a base type with its corresponding Create and Edit option types.
//
// Usage:
//
//	_ = CRUDPair{Base: "Repository", Create: "CreateRepoOption", Edit: "EditRepoOption"}
type CRUDPair struct {
	Base   string
	Create string
	Edit   string
}

// LoadSpec reads and parses a Swagger 2.0 JSON file from the given path.
//
// Usage:
//
//	spec, err := LoadSpec("testdata/swagger.v1.json")
//	_ = spec
func LoadSpec(path string) (*Spec, error) {
	content, err := coreio.Local.Read(path)
	if err != nil {
		return nil, core.E("LoadSpec", "read spec", err)
	}
	var spec Spec
	if err := json.Unmarshal([]byte(content), &spec); err != nil {
		return nil, core.E("LoadSpec", "parse spec", err)
	}
	return &spec, nil
}

// ExtractTypes converts all swagger definitions into Go type intermediate representations.
//
// Usage:
//
//	types := ExtractTypes(spec)
//	_ = types["Repository"]
func ExtractTypes(spec *Spec) map[string]*GoType {
	result := make(map[string]*GoType)
	for name, def := range spec.Definitions {
		gt := &GoType{Name: name, Description: def.Description}
		if len(def.Enum) > 0 {
			gt.IsEnum = true
			for _, v := range def.Enum {
				gt.EnumValues = append(gt.EnumValues, core.Sprint(v))
			}
			slices.Sort(gt.EnumValues)
			result[name] = gt
			continue
		}
		if len(def.Properties) == 0 && def.AdditionalProperties != nil {
			gt.IsAlias = true
			gt.AliasType = resolveMapType(*def.AdditionalProperties)
			result[name] = gt
			continue
		}
		required := make(map[string]bool)
		for _, r := range def.Required {
			required[r] = true
		}
		for fieldName, prop := range def.Properties {
			goName := prop.XGoName
			if goName == "" {
				goName = pascalCase(fieldName)
			}
			gf := GoField{
				GoName:   goName,
				GoType:   resolveGoType(prop),
				JSONName: fieldName,
				Comment:  prop.Description,
				Required: required[fieldName],
			}
			gt.Fields = append(gt.Fields, gf)
		}
		slices.SortFunc(gt.Fields, func(a, b GoField) int {
			return cmp.Compare(a.GoName, b.GoName)
		})
		result[name] = gt
	}
	return result
}

// DetectCRUDPairs finds Create*Option / Edit*Option pairs in the swagger definitions
// and maps them back to the base type name.
//
// Usage:
//
//	pairs := DetectCRUDPairs(spec)
//	_ = pairs
func DetectCRUDPairs(spec *Spec) []CRUDPair {
	var pairs []CRUDPair
	for name := range spec.Definitions {
		if !core.HasPrefix(name, "Create") || !core.HasSuffix(name, "Option") {
			continue
		}
		inner := core.TrimPrefix(name, "Create")
		inner = core.TrimSuffix(inner, "Option")
		editName := core.Concat("Edit", inner, "Option")
		pair := CRUDPair{Base: inner, Create: name}
		if _, ok := spec.Definitions[editName]; ok {
			pair.Edit = editName
		}
		pairs = append(pairs, pair)
	}
	slices.SortFunc(pairs, func(a, b CRUDPair) int {
		return cmp.Compare(a.Base, b.Base)
	})
	return pairs
}

// resolveGoType maps a swagger schema property to a Go type string.
func resolveGoType(prop SchemaProperty) string {
	if prop.Ref != "" {
		parts := core.Split(prop.Ref, "/")
		return "*" + parts[len(parts)-1]
	}
	switch prop.Type {
	case "string":
		switch prop.Format {
		case "date-time":
			return "time.Time"
		case "binary":
			return "[]byte"
		default:
			return "string"
		}
	case "integer":
		switch prop.Format {
		case "int64":
			return "int64"
		case "int32":
			return "int32"
		default:
			return "int"
		}
	case "number":
		switch prop.Format {
		case "float":
			return "float32"
		default:
			return "float64"
		}
	case "boolean":
		return "bool"
	case "array":
		if prop.Items != nil {
			return "[]" + resolveGoType(*prop.Items)
		}
		return "[]any"
	case "object":
		return resolveMapType(prop)
	default:
		return "any"
	}
}

// resolveMapType maps a swagger object with additionalProperties to a Go map type.
func resolveMapType(prop SchemaProperty) string {
	valueType := "any"
	if prop.AdditionalProperties != nil {
		valueType = resolveGoType(*prop.AdditionalProperties)
	}
	return "map[string]" + valueType
}

// pascalCase converts a snake_case or kebab-case string to PascalCase,
// with common acronyms kept uppercase.
func pascalCase(s string) string {
	var parts []string
	for _, p := range splitSnakeKebab(s) {
		if len(p) == 0 {
			continue
		}
		upper := core.Upper(p)
		switch upper {
		case "ID", "URL", "HTML", "SSH", "HTTP", "HTTPS", "API", "URI", "GPG", "IP", "CSS", "JS":
			parts = append(parts, upper)
		default:
			parts = append(parts, core.Concat(core.Upper(p[:1]), p[1:]))
		}
	}
	return core.Concat(parts...)
}
