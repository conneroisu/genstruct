// Package main shows how to export data from a package.
//
// The generated code will be in another package (./out) and
// will be named after the folder where the source code is generated.
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/genstruct"
	"github.com/conneroisu/genstruct/examples/exported-data-zoo/pkg"
)

// generateAnimalData generates the static animal data file
func generateAnimalData() error {
	// Define our array of animal data

	// Configure genstruct
	config := genstruct.Config{
		OutputFile:       "./out/zoo_animals.go",      // Output file name (absolute path from project root)
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

func main() {
	// First generate the animal data
	fmt.Println("Generating static animal data...")
	err := generateAnimalData()
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated static animal data in zoo_animals.go")

	// Show the content of the generated file
	content, err := os.ReadFile("./out/zoo_animals.go")
	if err != nil {
		fmt.Printf("Error reading generated file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nContents of generated file:")
	fmt.Println("---------------------------")
	fmt.Println(string(content))

	fmt.Println("\nTo use the generated code in your application you would:")
	fmt.Println("1. Import the generated file in your code by its package name (out)")
	fmt.Println("2. Use out.AnimalLeo, out.AnimalEllie, etc. to access specific animals")
	fmt.Println("3. Use out.AllAnimals slice for filtering and analysis")
}
