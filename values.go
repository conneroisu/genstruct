package genstruct

import (
	"fmt"
	"reflect"
	"time"

	"github.com/dave/jennifer/jen"
)

// getValueStatement generates code for a value based on its type
func (g *Generator) getValueStatement(value reflect.Value) *jen.Statement {
	switch value.Kind() {
	case reflect.Bool:
		return jen.Lit(value.Bool())
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		return jen.Lit(value.Int())
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		return jen.Lit(value.Uint())
	case reflect.Float32, reflect.Float64:
		return jen.Lit(value.Float())
	case reflect.Complex64, reflect.Complex128:
		return jen.Lit(value.Complex())
	case reflect.Array:
		// Handle arrays properly with their type and dimensions
		elemType := g.getTypeStatement(value.Type().Elem())
		dimensions := value.Len()

		// Create array type with dimensions
		arrayType := jen.Index(jen.Lit(dimensions)).Add(elemType)

		// Create values inside the array
		return arrayType.ValuesFunc(func(group *jen.Group) {
			for i := range value.Len() {
				group.Add(g.getValueStatement(value.Index(i)))
			}
		})
	case reflect.Slice:
		// Create a slice with proper syntax
		return jen.Index().Add(
			g.getTypeStatement(value.Type().Elem()),
		).ValuesFunc(func(group *jen.Group) {
			for i := range value.Len() {
				group.Add(g.getValueStatement(value.Index(i)))
			}
		})
	case reflect.Map:
		return g.getMapStatement(value)
	case reflect.String:
		return jen.Lit(value.String())
	case reflect.Struct:
		// Special case for time.Time
		if value.Type().String() == "time.Time" {
			t := value.Interface().(time.Time)
			return jen.Qual("time", "Date").Call(
				jen.Lit(t.Year()),
				jen.Qual("time", t.Month().String()),
				jen.Lit(t.Day()),
				jen.Lit(t.Hour()),
				jen.Lit(t.Minute()),
				jen.Lit(t.Second()),
				jen.Lit(t.Nanosecond()),
				jen.Qual("time", "UTC"),
			)
		}
		// For other structs, create a new values block with the struct fields
		return jen.Id(
			value.Type().Name(),
		).ValuesFunc(func(group *jen.Group) {
			g.generateStructValues(group, value)
		})
	case reflect.Pointer:
		if value.IsNil() {
			return jen.Nil()
		}
		return jen.Op("&").Add(g.getValueStatement(value.Elem()))
	case reflect.Interface:
		if value.IsNil() {
			return jen.Nil()
		}
		return g.getValueStatement(value.Elem())
	default:
		// For complex cases, fallback to string representation
		return jen.Lit(fmt.Sprintf("%v", value.Interface()))
	}
}

// getMapStatement generates code for a map
func (g *Generator) getMapStatement(mapValue reflect.Value) *jen.Statement {
	// Return empty map if there are no entries
	if mapValue.Len() == 0 {
		return jen.Map(
			g.getTypeStatement(mapValue.Type().Key()),
		).Add(
			g.getTypeStatement(mapValue.Type().Elem()),
		).Values()
	}

	// Use ValuesFunc for populated maps
	return jen.Map(
		g.getTypeStatement(mapValue.Type().Key()),
	).Add(
		g.getTypeStatement(mapValue.Type().Elem()),
	).ValuesFunc(func(group *jen.Group) {
		// Create a Dict for the map entries
		dict := jen.Dict{}
		var key reflect.Value

		// Add all key-value pairs to the Dict
		for _, key = range mapValue.MapKeys() {
			var stmt = g.getValueStatement(mapValue.MapIndex(key))
			dict[g.getValueStatement(key)] = stmt
		}

		// Add dict to group
		group.Add(dict)
	})
}

// generateStructValues adds values for a struct to a Dict
func (g *Generator) generateStructValues(group *jen.Group, structValue reflect.Value) {
	structType := structValue.Type()

	// Create a Dict for each field in the struct
	dict := jen.Dict{}

	for i := range structValue.NumField() {
		field := structValue.Field(i)
		fieldType := structType.Field(i)

		// Skip unexported fields
		if !fieldType.IsExported() {
			continue
		}

		// Add the field to the dict
		dict[jen.Id(fieldType.Name)] = g.getValueStatement(field)
	}

	// Add all fields to the group
	group.Add(dict)
}
