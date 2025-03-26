package genstruct

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/dave/jennifer/jen"
)

// Generator is responsible for generating code for static struct arrays
type Generator struct {
	Config Config
	Data   any            // The primary array of structs to generate code for
	Refs   map[string]any // Additional arrays that can be referenced
	File   *jen.File
}

// NewGenerator creates a new generator instance with support for struct references
//
// Parameters:
//   - config: Configuration options for code generation
//   - data: The primary array of structs to generate code for
//   - refs: Optional additional arrays that can be referenced by the primary data
//
// The refs parameters enable the use of struct tags with the `structgen` tag to reference
// data between structs. For example, a Post struct with a TagSlugs field can reference
// Tag structs using the tag `structgen:"TagSlugs"`.
//
// When Generate() is called, it will:
// 1. Generate constants, variables, and a slice for the primary data
// 2. Generate constants, variables, and a slice for each referenced data set
// 3. Create references between the primary data and the referenced data
//
// Example usage:
//
//	generator, err := genstruct.NewGenerator(config, posts, tags)
//	if err != nil {
//	    // handle error
//	}
//
// Many configuration options can be omitted and will be inferred automatically:
//   - TypeName: Inferred from the struct type in the data slice
//   - ConstantIdent: Defaults to TypeName if not specified
//   - VarPrefix: Defaults to TypeName if not specified
//   - OutputFile: Defaults to lowercase(typename_generated.go) if not specified
//   - IdentifierFields: Has reasonable defaults if not specified
//
// Returns an error if:
//   - The data is not a slice or array
//   - The data is empty (no elements to analyze)
//   - The data elements are not structs
//   - Required fields couldn't be inferred
func NewGenerator(config Config, data any, refs ...any) (*Generator, error) {
	// Validate and enhance the configuration with inferred values
	enhancedConfig, err := enhanceConfig(config, data)
	if err != nil {
		return nil, err
	}

	// Create a map of reference datasets
	refMap := make(map[string]any)
	for i, ref := range refs {
		// Get type name for this dataset
		refType := reflect.TypeOf(ref)
		if refType.Kind() == reflect.Slice || refType.Kind() == reflect.Array {
			elemType := refType.Elem()
			if elemType.Kind() == reflect.Struct {
				refMap[elemType.Name()] = ref
			} else {
				refMap[fmt.Sprintf("Ref%d", i)] = ref
			}
		}
	}

	return &Generator{
		Config: enhancedConfig,
		Data:   data,
		Refs:   refMap,
		File:   jen.NewFile(enhancedConfig.PackageName),
	}, nil
}

// enhanceConfig fills in missing configuration values using reflection
func enhanceConfig(config Config, data any) (Config, error) {
	// Get the element type from the data
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice && dataValue.Kind() != reflect.Array {
		// Can't determine type from non-slice/array, return as is
		return config, InvalidTypeError{Kind: dataValue.Kind()}
	}

	// Make sure we have at least one element to analyze
	if dataValue.Len() == 0 {
		// Can't determine type from empty slice, return as is
		return config, EmptyError{}
	}

	firstElem := dataValue.Index(0)
	if firstElem.Kind() != reflect.Struct {
		// Only struct slices are supported, return as is
		return config, InvalidTypeError{Kind: firstElem.Kind()}
	}

	// Get the struct type
	structType := firstElem.Type()
	typeName := structType.Name()

	// Infer TypeName if not specified
	if config.TypeName == "" {
		config.TypeName = typeName
	}

	// Infer ConstantIdent if not specified
	if config.ConstantIdent == "" {
		config.ConstantIdent = config.TypeName
	}

	// Infer VarPrefix if not specified
	if config.VarPrefix == "" {
		config.VarPrefix = config.TypeName
	}

	// Infer OutputFile if not specified
	if config.OutputFile == "" {
		config.OutputFile = strings.ToLower(config.TypeName) + "_generated.go"
	}

	// Set default identifier fields if none provided
	if config.IdentifierFields == nil {
		config.IdentifierFields = []string{"ID", "Name", "Slug", "Title", "Key", "Code"}
	}

	// If PackageName is not specified, use "generated"
	if config.PackageName == "" {
		config.PackageName = "generated"
	}

	return config, nil
}

// Generate performs the code generation for both primary data and reference data
//
// This method generates:
// 1. Constants for the primary data's IDs
// 2. Variables for each item in the primary data
// 3. A slice containing all primary data items
// 4. Constants for each reference data set's IDs
// 5. Variables for each item in each reference data set
// 6. A slice for each reference data set
// 7. Creates references between primary data and reference data as specified by structgen tags
//
// All generated code is written to a single output file specified in the Config.
func (g *Generator) Generate() error {
	g.File.PackageComment(fmt.Sprintf(
		"// Code generated by genstruct. DO NOT EDIT.\n// Package %s contains auto-generated %s data",
		g.Config.PackageName,
		g.Config.TypeName,
	))

	// Validate that we have an array or slice
	dataValue := reflect.ValueOf(g.Data)
	if dataValue.Kind() != reflect.Slice &&
		dataValue.Kind() != reflect.Array {
		return NonSliceOrArrayError{dataValue.Kind()}
	}

	// Make sure we have at least one element to analyze the type
	if dataValue.Len() == 0 {
		return EmptyError{}
	}

	// Get the type of the first element
	firstElem := dataValue.Index(0)
	if firstElem.Kind() != reflect.Struct {
		return InvalidTypeError{firstElem.Kind()}
	}

	// Generate constants for IDs if there's an ID field
	g.generateConstants(dataValue)

	// Generate variables for each struct
	g.generateVariables(dataValue)

	// Generate a slice with all structs
	g.generateSlice(dataValue)

	// Process reference datasets to generate their constants and variables
	// This ensures that all referenced types (like Tag in Post.Tags) are properly defined
	// in the generated code, making the references fully usable.
	for typeName, refDataObj := range g.Refs {
		refDataValue := reflect.ValueOf(refDataObj)
		if refDataValue.Kind() == reflect.Slice || refDataValue.Kind() == reflect.Array {
			if refDataValue.Len() > 0 {
				refElem := refDataValue.Index(0)
				if refElem.Kind() == reflect.Struct {
					// Store original config values so we can restore them after
					// processing this reference type
					originalTypeName := g.Config.TypeName
					originalVarPrefix := g.Config.VarPrefix
					originalConstantIdent := g.Config.ConstantIdent

					// Temporarily set config values for the reference type
					// This ensures that constants and variables are named correctly
					// (e.g., TagGoProgramming instead of PostGoProgramming)
					g.Config.TypeName = typeName
					g.Config.VarPrefix = typeName
					g.Config.ConstantIdent = typeName

					// Generate constants, variables, and slice for this reference dataset
					// using the same generation methods as for the primary dataset
					g.generateConstants(refDataValue)
					g.generateVariables(refDataValue)
					g.generateSlice(refDataValue)

					// Restore original config values for processing the next reference dataset
					g.Config.TypeName = originalTypeName
					g.Config.VarPrefix = originalVarPrefix
					g.Config.ConstantIdent = originalConstantIdent
				}
			}
		}
	}

	// Generate the code as a string
	buf := &bytes.Buffer{}
	if err := g.File.Render(buf); err != nil {
		return err
	}

	// Format the code with gofmt
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	// Save the formatted code to file
	return os.WriteFile(g.Config.OutputFile, formatted, 0644)
}

// getStructIdentifier returns a string to identify this struct instance
func (g *Generator) getStructIdentifier(structValue reflect.Value) string {
	// If a custom name function is provided, use it
	if g.Config.CustomVarNameFn != nil {
		return g.Config.CustomVarNameFn(structValue)
	}

	// Try all configured identifier fields
	for _, fieldName := range g.Config.IdentifierFields {
		field := structValue.FieldByName(fieldName)
		if field.IsValid() && field.Kind() == reflect.String && field.String() != "" {
			return field.String()
		}
	}

	// Fallback 1: Look for any string field
	for i := range structValue.NumField() {
		field := structValue.Field(i)
		if field.Kind() == reflect.String && field.String() != "" {
			return field.String()
		}
	}

	// Fallback 2: Generate a name based on the type
	return fmt.Sprintf("%s-%d", g.Config.TypeName, time.Now().UnixNano())
}
