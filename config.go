package genstruct

import (
	"reflect"
	"regexp"
	"strings"
)

// Config holds the configuration for code generation of static structs and arrays.
type Config struct {
	// PackageName defines the Target package name (optional)
	// If not provided, the package name will be the same as the given struct.
	PackageName      string
	TypeName         string   // The name of the struct type to generate
	ConstantIdent    string   // Prefix for constants (e.g., "Post" for "PostMyPostID")
	VarPrefix        string   // Prefix for variables (e.g., "Post" for "PostMyPost")
	OutputFile       string   // Output file name
	IdentifierFields []string // Fields to try using for naming, in priority order (optional)
	// Custom function to generate variable names (optional)
	// If provided, this takes precedence over IdentifierFields
	CustomVarNameFn func(structValue reflect.Value) string
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
