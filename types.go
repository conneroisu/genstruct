package genstruct

import (
	"reflect"
	"strings"

	"github.com/dave/jennifer/jen"
)

// getTypeStatement converts a reflect.Type to a jen.Statement
func (g *Generator) getTypeStatement(t reflect.Type) *jen.Statement {
	switch t.Kind() {
	case reflect.Bool:
		return jen.Bool()
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return jen.Id(t.String())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr:
		return jen.Id(t.String())
	case reflect.Float32, reflect.Float64:
		return jen.Id(t.String())
	case reflect.Complex64, reflect.Complex128:
		return jen.Id(t.String())
	case reflect.Array, reflect.Slice:
		elemType := t.Elem()
		// Special handling for []*Type pattern
		if elemType.Kind() == reflect.Pointer {
			return jen.Index().Add(jen.Op("*").Add(g.getTypeStatement(elemType.Elem())))
		}
		return jen.Index().Add(g.getTypeStatement(elemType))
	case reflect.Map:
		return jen.Map(
			g.getTypeStatement(t.Key()),
		).Add(g.getTypeStatement(t.Elem()))
	case reflect.String:
		return jen.String()
	case reflect.Struct:
		// Handle special types like time.Time
		if t.String() == "time.Time" {
			return jen.Qual("time", "Time")
		}
		
		// Check if this is from a different package (has a dot in the name)
		pkgPath := t.PkgPath()
		// Infer ExportDataMode by checking if output file contains package path separator
		isExportMode := strings.Contains(g.OutputFile, "/")
		if pkgPath != "" && pkgPath != "main" && pkgPath != g.PackageName && isExportMode {
			// If the type comes from a different package, reference it with the package name
			pkgName := t.String()
			if lastDot := strings.LastIndex(pkgName, "."); lastDot >= 0 {
				return jen.Qual(pkgPath, t.Name())
			}
		}
		return jen.Id(t.Name())
	case reflect.Pointer:
		return jen.Op("*").Add(g.getTypeStatement(t.Elem()))
	case reflect.Interface:
		if t.NumMethod() == 0 {
			return jen.Interface() // empty interface
		}
		// Complex interfaces would need more handling
		return jen.Interface()
	default:
		return jen.Id(t.String())
	}
}

// collectEmbeddedTypes finds all embedded struct types in a given struct type
// and adds them to the exportedTypes map
func (g *Generator) collectEmbeddedTypes(structType reflect.Type, exportedTypes map[string]reflect.Type) {
	if structType.Kind() == reflect.Pointer {
		structType = structType.Elem()
	}
	
	if structType.Kind() != reflect.Struct {
		return
	}
	
	// Check each field for embedded structs
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		
		// If the field is an embedded struct (anonymous field)
		if field.Anonymous {
			embeddedType := field.Type
			if embeddedType.Kind() == reflect.Pointer {
				embeddedType = embeddedType.Elem()
			}
			
			if embeddedType.Kind() == reflect.Struct {
				// Add this embedded type to our map if it's not from standard library
				if embeddedType.PkgPath() != "" && !strings.HasPrefix(embeddedType.PkgPath(), "time") {
					exportedTypes[embeddedType.Name()] = embeddedType
					// Recursively check this type for its own embedded types
					g.collectEmbeddedTypes(embeddedType, exportedTypes)
				}
			}
		}
	}
}

// parseTags turns a reflect.StructTag into a map[string]string for jennifer
func parseTags(tag reflect.StructTag) map[string]string {
	tagMap := make(map[string]string)
	
	// Extract supported tags
	if yaml, ok := tag.Lookup("yaml"); ok {
		tagMap["yaml"] = yaml
	}
	
	if structgen, ok := tag.Lookup("structgen"); ok {
		tagMap["structgen"] = structgen
	}
	
	return tagMap
}

// exportStructType exports a struct type definition to the generated code
func (g *Generator) exportStructType(typeName string, structType reflect.Type) {
	// Create the code block for the struct definition
	g.File.Type().Id(typeName).StructFunc(func(s *jen.Group) {
		for i := 0; i < structType.NumField(); i++ {
			field := structType.Field(i)
			
			// Skip unexported fields
			if !field.IsExported() {
				continue
			}
			
			// Handle embedded structs (anonymous fields)
			if field.Anonymous {
				if field.Type.Kind() == reflect.Pointer {
					s.Op("*").Id(field.Type.Elem().Name())
				} else {
					s.Id(field.Type.Name())
				}
				continue
			}
			
			// Regular field with type
			fieldDef := jen.Id(field.Name).Add(g.getTypeStatement(field.Type))
			
			// Add tags if present
			tagMap := parseTags(field.Tag)
			if len(tagMap) > 0 {
				fieldDef.Tag(tagMap)
			}
			
			s.Add(fieldDef)
		}
	})
}