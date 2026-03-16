package main

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	coreio "forge.lthn.ai/core/go-io"
	coreerr "forge.lthn.ai/core/go-log"
)

// Spec represents a Swagger 2.0 specification document.
type Spec struct {
	Swagger     string                      `json:"swagger"`
	Info        SpecInfo                    `json:"info"`
	Definitions map[string]SchemaDefinition `json:"definitions"`
	Paths       map[string]map[string]any   `json:"paths"`
}

// SpecInfo holds metadata about the API specification.
type SpecInfo struct {
	Title   string `json:"title"`
	Version string `json:"version"`
}

// SchemaDefinition represents a single type definition in the swagger spec.
type SchemaDefinition struct {
	Description string                    `json:"description"`
	Type        string                    `json:"type"`
	Properties  map[string]SchemaProperty `json:"properties"`
	Required    []string                  `json:"required"`
	Enum        []any                     `json:"enum"`
	XGoName     string                    `json:"x-go-name"`
}

// SchemaProperty represents a single property within a schema definition.
type SchemaProperty struct {
	Type        string          `json:"type"`
	Format      string          `json:"format"`
	Description string          `json:"description"`
	Ref         string          `json:"$ref"`
	Items       *SchemaProperty `json:"items"`
	Enum        []any           `json:"enum"`
	XGoName     string          `json:"x-go-name"`
}

// GoType is the intermediate representation for a Go type to be generated.
type GoType struct {
	Name        string
	Description string
	Fields      []GoField
	IsEnum      bool
	EnumValues  []string
}

// GoField is the intermediate representation for a single struct field.
type GoField struct {
	GoName   string
	GoType   string
	JSONName string
	Comment  string
	Required bool
}

// CRUDPair groups a base type with its corresponding Create and Edit option types.
type CRUDPair struct {
	Base   string
	Create string
	Edit   string
}

// LoadSpec reads and parses a Swagger 2.0 JSON file from the given path.
func LoadSpec(path string) (*Spec, error) {
	content, err := coreio.Local.Read(path)
	if err != nil {
		return nil, coreerr.E("LoadSpec", "read spec", err)
	}
	var spec Spec
	if err := json.Unmarshal([]byte(content), &spec); err != nil {
		return nil, coreerr.E("LoadSpec", "parse spec", err)
	}
	return &spec, nil
}

// ExtractTypes converts all swagger definitions into Go type intermediate representations.
func ExtractTypes(spec *Spec) map[string]*GoType {
	result := make(map[string]*GoType)
	for name, def := range spec.Definitions {
		gt := &GoType{Name: name, Description: def.Description}
		if len(def.Enum) > 0 {
			gt.IsEnum = true
			for _, v := range def.Enum {
				gt.EnumValues = append(gt.EnumValues, fmt.Sprintf("%v", v))
			}
			slices.Sort(gt.EnumValues)
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
			return strings.Compare(a.GoName, b.GoName)
		})
		result[name] = gt
	}
	return result
}

// DetectCRUDPairs finds Create*Option / Edit*Option pairs in the swagger definitions
// and maps them back to the base type name.
func DetectCRUDPairs(spec *Spec) []CRUDPair {
	var pairs []CRUDPair
	for name := range spec.Definitions {
		if !strings.HasPrefix(name, "Create") || !strings.HasSuffix(name, "Option") {
			continue
		}
		inner := strings.TrimPrefix(name, "Create")
		inner = strings.TrimSuffix(inner, "Option")
		editName := "Edit" + inner + "Option"
		pair := CRUDPair{Base: inner, Create: name}
		if _, ok := spec.Definitions[editName]; ok {
			pair.Edit = editName
		}
		pairs = append(pairs, pair)
	}
	slices.SortFunc(pairs, func(a, b CRUDPair) int {
		return strings.Compare(a.Base, b.Base)
	})
	return pairs
}

// resolveGoType maps a swagger schema property to a Go type string.
func resolveGoType(prop SchemaProperty) string {
	if prop.Ref != "" {
		parts := strings.Split(prop.Ref, "/")
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
		return "map[string]any"
	default:
		return "any"
	}
}

// pascalCase converts a snake_case or kebab-case string to PascalCase,
// with common acronyms kept uppercase.
func pascalCase(s string) string {
	var parts []string
	for p := range strings.FieldsFuncSeq(s, func(r rune) bool {
		return r == '_' || r == '-'
	}) {
		if len(p) == 0 {
			continue
		}
		upper := strings.ToUpper(p)
		switch upper {
		case "ID", "URL", "HTML", "SSH", "HTTP", "HTTPS", "API", "URI", "GPG", "IP", "CSS", "JS":
			parts = append(parts, upper)
		default:
			parts = append(parts, strings.ToUpper(p[:1])+p[1:])
		}
	}
	return strings.Join(parts, "")
}
