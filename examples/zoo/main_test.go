package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestGeneratedAnimalData(t *testing.T) {
	if err := generateAnimalData(); err != nil {
		t.Fatalf("Failed to generate animal data: %v", err)
	}
	// validate the generated go file
	if err := validateGeneratedFile(); err != nil {
		t.Fatalf("Failed to validate generated file: %v", err)
	}
}

func validateGeneratedFile() error {
	// Use go/ast to validate the generated file
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "zoo_animals.go", nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("failed to parse generated file: %v", err)
	}

	// Validate package name
	if file.Name.Name != "main" {
		return fmt.Errorf("expected package name to be 'zoo', got '%s'", file.Name.Name)
	}

	// Expected animal names and IDs
	expectedIDConstants := []string{"AnimalLeoID", "AnimalEllieID", "AnimalStripesID", "AnimalPennyID", "AnimalGeorgeID"}
	expectedVariables := []string{"AnimalLeo", "AnimalEllie", "AnimalStripes", "AnimalPenny", "AnimalGeorge"}

	// Track what we've found
	foundIDConstants := make(map[string]bool)
	foundVariables := make(map[string]bool)
	foundAllAnimalsSlice := false

	// Walk the AST to validate declarations
	ast.Inspect(file, func(n ast.Node) bool {
		// Check for the constants
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.CONST {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						// Check if it's an animal ID constant
						for _, expectedID := range expectedIDConstants {
							if name.Name == expectedID {
								foundIDConstants[expectedID] = true
							}
						}
					}
				}
			}
		}

		// Check for variables
		if genDecl, ok := n.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if valueSpec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range valueSpec.Names {
						// Check if it's an animal variable
						for _, expectedVar := range expectedVariables {
							if name.Name == expectedVar {
								foundVariables[expectedVar] = true
							}
						}

						// Check for AllAnimals slice
						if name.Name == "AllAnimals" {
							foundAllAnimalsSlice = true
						}
					}
				}
			}
		}

		return true
	})

	// Validate all expected ID constants were found
	missingIDs := []string{}
	for _, expectedID := range expectedIDConstants {
		if !foundIDConstants[expectedID] {
			missingIDs = append(missingIDs, expectedID)
		}
	}
	if len(missingIDs) > 0 {
		return fmt.Errorf("missing expected ID constants: %s", strings.Join(missingIDs, ", "))
	}

	// Validate all expected variables were found
	missingVars := []string{}
	for _, expectedVar := range expectedVariables {
		if !foundVariables[expectedVar] {
			missingVars = append(missingVars, expectedVar)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("missing expected variables: %s", strings.Join(missingVars, ", "))
	}

	// Validate AllAnimals slice was found
	if !foundAllAnimalsSlice {
		return fmt.Errorf("missing AllAnimals slice")
	}

	return nil
}
