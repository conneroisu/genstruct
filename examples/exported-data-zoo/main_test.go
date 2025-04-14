// Package main shows how to export data from a package.
//
// The generated code will be in another package (./out) and
// will be named after the folder where the source code is generated.
package main_test

import (
	"os"
	"strings"
	"testing"
	
	"github.com/conneroisu/genstruct"
	"github.com/conneroisu/genstruct/examples/exported-data-zoo/pkg"
)

// generateAnimalData generates the static animal data file
func generateAnimalData() error {
	// Define our array of animal data

	// Create a generator with functional options
	// Note: ExportDataMode is inferred automatically from the output file path
	generator := genstruct.NewGenerator(
		genstruct.WithPackageName("out"),                      // Target package name
		genstruct.WithTypeName("Animal"),                      // The struct type name
		genstruct.WithConstantIdent("Animal"),                 // Prefix for constants
		genstruct.WithVarPrefix("Animal"),                     // Prefix for variables
		genstruct.WithOutputFile("out/zoo_animals.go"),        // Output file name (relative to test directory)
		genstruct.WithIdentifierFields([]string{"Name", "Species"}), // Fields to use for naming variables
	)

	// Generate the code, passing animals data
	return generator.Generate(pkg.Animals)
}
func TestExportedDataZoo(t *testing.T) {
	// Generate the animal data
	err := generateAnimalData()
	if err != nil {
		t.Fatalf("Error generating animal data: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile("out/zoo_animals.go")
	if err != nil {
		t.Fatalf("Error reading generated file: %v", err)
	}

	contentStr := string(content)

	// Test that generated code contains expected elements
	expectedTests := []struct {
		name     string
		expected string
		message  string
	}{
		{
			name:     "Package declaration",
			expected: "package out",
			message:  "Should declare output package",
		},
		{
			name:     "Import pkg package",
			expected: "github.com/conneroisu/genstruct/examples/exported-data-zoo/pkg",
			message:  "Should import the pkg package",
		},
		{
			name:     "Animal variable with package qualification",
			expected: "pkg.Animal{",
			message:  "Should use pkg.Animal type",
		},
		{
			name:     "Animal slice with package qualification",
			expected: "[]*pkg.Animal{",
			message:  "Should use pkg.Animal for slice",
		},
	}

	// Run all tests
	for _, tc := range expectedTests {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(contentStr, tc.expected) {
				t.Errorf("%s: %q not found in generated code", tc.message, tc.expected)
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	// Clean up the generated file
	err := os.Remove("out/zoo_animals.go")
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error removing generated file: %v", err)
	}
}
