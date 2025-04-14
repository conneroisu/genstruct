package genstruct

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"runtime/debug"
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
// Reference fields can be either direct structs or pointers to structs:
//   - []Tag `structgen:"TagSlugs"` - Direct struct references
//   - []*Tag `structgen:"TagSlugs"` - Pointer-based struct references (recommended)
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
			} else if elemType.Kind() == reflect.Pointer &&
				elemType.Elem().Kind() == reflect.Struct {
				// Handle pointer slice ([]*Type)
				refMap[elemType.Elem().Name()] = ref
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
	var structType reflect.Type

	// Support both direct struct slices and pointer slices
	if firstElem.Kind() == reflect.Struct {
		structType = firstElem.Type()
	} else if firstElem.Kind() == reflect.Pointer && firstElem.Elem().Kind() == reflect.Struct {
		structType = firstElem.Elem().Type()
	} else {
		// Only struct or struct pointer slices are supported
		return config, InvalidTypeError{Kind: firstElem.Kind()}
	}

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
		config.IdentifierFields = []string{
			"ID",
			"Name",
			"Slug",
			"Title",
			"Key",
			"Code",
		}
	}

	// If PackageName is not specified, use "generated"
	if config.PackageName == "" {
		config.PackageName = GetPackageNameFromPath(config.OutputFile)
	}

	// If Logger is not specified, use the default logger
	if config.Logger == nil {
		config.Logger = GetLogger()
	}

	// Log the configuration
	config.Logger.Debug(
		"Configuration enhanced",
		slog.String("typeName", config.TypeName),
		slog.String("packageName", config.PackageName),
		slog.String("outputFile", config.OutputFile),
	)

	return config, nil
}

// GetPackageNameFromPath extracts the containing folder name from a file path
// This can be used to determine the package name for a given Go file
// Example: "./out/penguin/gen.go" would return "penguin"
func GetPackageNameFromPath(filePath string) string {
	// Clean the path to handle any OS-specific separators and normalize it
	cleanPath := filepath.Clean(filePath)

	// Get the directory containing the file
	dir := filepath.Dir(cleanPath)

	// Split the directory path into components
	components := strings.Split(dir, string(filepath.Separator))

	// The package name is the last component of the directory path
	// If the path ends with a separator, the last component will be empty
	if len(components) > 0 {
		lastComponent := components[len(components)-1]
		if lastComponent != "" {
			return lastComponent
		}

		// If the last component is empty, try the second-to-last one
		if len(components) > 1 {
			return components[len(components)-2]
		}
	}

	// Default to "main" if we couldn't extract a package name
	return "main"
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
	g.Config.Logger.Info(
		"Starting code generation",
		slog.String("package", g.Config.PackageName),
		slog.String("type", g.Config.TypeName),
		slog.String("output", g.Config.OutputFile),
	)
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return fmt.Errorf("failed to read build info for version number")
	}
	/// find github.com/conneroisu/genstruct dep
	var dep *debug.Module
	for _, dep = range bi.Deps {
		if dep.Path == "github.com/conneroisu/genstruct" {
			break
		}
	}
	if dep == nil {
		dep = &debug.Module{
			Path:    "github.com/conneroisu/genstruct",
			Version: "Unknown",
		}
	}
	g.File.PackageComment(fmt.Sprintf(
		"// Code generated by genstruct. DO NOT EDIT.\n// Package %s contains auto-generated %s data\n//\n// genstruct Version: %s\n//",
		g.Config.PackageName,
		g.Config.TypeName,
		dep.Version,
	))

	// Validate that we have an array or slice
	dataValue := reflect.ValueOf(g.Data)
	if dataValue.Kind() != reflect.Slice &&
		dataValue.Kind() != reflect.Array {
		g.Config.Logger.Error(
			"Invalid data type",
			"expected",
			"slice or array",
			"got",
			dataValue.Kind().String(),
		)
		return NonSliceOrArrayError{dataValue.Kind()}
	}

	// Make sure we have at least one element to analyze the type
	if dataValue.Len() == 0 {
		g.Config.Logger.Error("Empty data slice", "type", g.Config.TypeName)
		return EmptyError{}
	}

	// Get the type of the first element
	firstElem := dataValue.Index(0)
	// Support both direct struct slices and pointer slices
	if firstElem.Kind() != reflect.Struct &&
		(firstElem.Kind() != reflect.Pointer ||
			firstElem.Elem().Kind() != reflect.Struct) {
		g.Config.Logger.Error(
			"Invalid element type",
			slog.String("expected", "struct or pointer to struct"),
			slog.String("got", firstElem.Kind().String()),
		)
		return InvalidTypeError{firstElem.Kind()}
	}

	// Generate constants for IDs if there's an ID field
	g.Config.Logger.Debug(
		"Generating constants",
		"type",
		g.Config.TypeName,
	)
	g.generateConstants(dataValue)

	// Generate variables for each struct
	g.Config.Logger.Debug(
		"Generating variables",
		"type",
		g.Config.TypeName,
		"count",
		dataValue.Len(),
	)
	g.generateVariables(dataValue)

	// Generate a slice with all structs
	g.Config.Logger.Debug(
		"Generating slice",
		"type",
		g.Config.TypeName,
	)
	g.generateSlice(dataValue)

	// Process reference datasets to generate their constants and variables
	// This ensures that all referenced types (like Tag in Post.Tags) are properly defined
	// in the generated code, making the references fully usable.
	g.Config.Logger.Debug(
		"Processing reference datasets",
		slog.Int("count", len(g.Refs)),
	)
	for typeName, refDataObj := range g.Refs {
		g.Config.Logger.Debug(
			"Processing reference dataset",
			slog.String("type", typeName),
		)
		refDataValue := reflect.ValueOf(refDataObj)
		if refDataValue.Kind() == reflect.Slice ||
			refDataValue.Kind() == reflect.Array {
			if refDataValue.Len() > 0 {
				refElem := refDataValue.Index(0)
				// Support both direct structs and pointer-to-structs
				if refElem.Kind() == reflect.Struct ||
					(refElem.Kind() == reflect.Pointer &&
						refElem.Elem().Kind() == reflect.Struct) {
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
	g.Config.Logger.Debug("Rendering generated code")
	buf := &bytes.Buffer{}
	if err := g.File.Render(buf); err != nil {
		g.Config.Logger.Error("Failed to render code", "error", err)
		return err
	}

	// Save the formatted code to file
	g.Config.Logger.Debug(
		"Writing generated code to file",
		slog.String("file", g.Config.OutputFile),
	)
	return os.WriteFile(g.Config.OutputFile, buf.Bytes(), 0644)
}

// getStructIdentifier returns a string to identify this struct instance
func (g *Generator) getStructIdentifier(structValue reflect.Value) string {
	// Handle pointer to struct case
	if structValue.Kind() == reflect.Pointer {
		structValue = structValue.Elem()
	}

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
