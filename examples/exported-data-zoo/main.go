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

	// Create a generator with functional options
	// Note: ExportDataMode is inferred automatically from the output file path
	generator := genstruct.NewGenerator(
		genstruct.WithOutputFile("./out/zoo_animals.go"),       // Output file name (absolute path from project root)
		genstruct.WithIdentifierFields([]string{"Name", "Species"}), // Fields to use for naming variables
	)

	// Generate the code, passing animals data
	return generator.Generate(pkg.Animals)
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
