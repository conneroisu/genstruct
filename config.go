package genstruct

import (
	"log/slog"
	"reflect"
	"regexp"
	"strings"
)

// Config holds the configuration for code generation of static structs and arrays.
// Many fields are optional and will be automatically inferred if not specified.
type Config struct {
	// PackageName defines the target package name
	// If not provided, defaults to "generated"
	PackageName string

	// TypeName is the name of the struct type to generate
	// If not provided, inferred from the struct type in the data
	TypeName string
	
	// ConstantIdent is the prefix for constants (e.g., "Post" for "PostMyPostID")
	// If not provided, defaults to the TypeName
	ConstantIdent string
	
	// VarPrefix is the prefix for variables (e.g., "Post" for "PostMyPost")
	// If not provided, defaults to the TypeName
	VarPrefix string
	
	// OutputFile is the output file name
	// If not provided, defaults to lowercase(typename_generated.go)
	OutputFile string
	
	// IdentifierFields are the fields to try using for naming, in priority order
	// If not provided, defaults to ["ID", "Name", "Slug", "Title", "Key", "Code"]
	IdentifierFields []string
	
	// CustomVarNameFn is a custom function to generate variable names (optional)
	// If provided, this takes precedence over IdentifierFields
	CustomVarNameFn func(structValue reflect.Value) string
	
	// Logger is the slog.Logger instance to use for logging
	// If not provided, defaults to a no-op logger
	Logger *slog.Logger
}

// slugToIdentifier converts a string to a valid Go identifier
func slugToIdentifier(s string) string {
	// Replace non-alphanumeric characters with spaces
	reg := regexp.MustCompile("[^a-zA-Z0-9]+")
	processed := reg.ReplaceAllString(s, " ")

	// Title case each word and remove spaces
	words := strings.Fields(processed)
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[0:1]) + strings.ToLower(word[1:])
		}
	}

	return strings.Join(words, "")
}
