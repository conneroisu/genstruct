// Package main shows how to export data from a package.
//
// The generated code will be in another package (./out) and
// will be named after the folder where the source code is generated.
package main_test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	
	"github.com/conneroisu/genstruct"
	"github.com/conneroisu/genstruct/examples/exported-data-zoo/pkg"
)

// generateAnimalData generates the static animal data file
func generateAnimalData() error {
	// Define our array of animal data

	// Configure genstruct
	config := genstruct.Config{
		PackageName:      "out",                       // Target package name
		TypeName:         "Animal",                    // The struct type name
		ConstantIdent:    "Animal",                    // Prefix for constants
		VarPrefix:        "Animal",                    // Prefix for variables
		OutputFile:       "out/zoo_animals.go",        // Output file name (relative to test directory)
		IdentifierFields: []string{"Name", "Species"}, // Fields to use for naming variables
		ExportDataMode:   true,                        // Enable referencing types from other packages
	}

	// Create a new generator with our config and animal data
	generator, err := genstruct.NewGenerator(config, pkg.Animals)
	if err != nil {
		return fmt.Errorf("error creating generator: %w", err)
	}

	// Generate the code
	return generator.Generate()
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
