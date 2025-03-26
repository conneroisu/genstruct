package main

import (
	"fmt"
	"os"
	"time"

	"github.com/conneroisu/genstruct"
)

// Animal represents an animal in a zoo
type Animal struct {
	ID           string    // Unique identifier for the animal
	Name         string    // Name of the animal
	Species      string    // Species of the animal
	DateOfBirth  time.Time // Date when the animal was born
	Diet         string    // Diet type (e.g. Carnivore, Herbivore)
	Weight       float64   // Weight in kilograms
	IsEndangered bool      // Whether the species is endangered
	Habitat      string    // Natural habitat type
	Region       string    // Geographic region of origin
}

// generateAnimalData generates the static animal data file
func generateAnimalData() error {
	// Define our array of animal data
	animals := []Animal{
		{
			ID:           "lion-001",
			Name:         "Leo",
			Species:      "African Lion",
			DateOfBirth:  time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
			Diet:         "Carnivore",
			Weight:       180.5,
			IsEndangered: true,
			Habitat:      "Savanna",
			Region:       "Africa",
		},
		{
			ID:           "elephant-002",
			Name:         "Ellie",
			Species:      "African Elephant",
			DateOfBirth:  time.Date(2012, time.June, 22, 0, 0, 0, 0, time.UTC),
			Diet:         "Herbivore",
			Weight:       3200.75,
			IsEndangered: false,
			Habitat:      "Savanna",
			Region:       "Africa",
		},
		{
			ID:           "tiger-003",
			Name:         "Stripes",
			Species:      "Bengal Tiger",
			DateOfBirth:  time.Date(2019, time.February, 8, 0, 0, 0, 0, time.UTC),
			Diet:         "Carnivore",
			Weight:       160.3,
			IsEndangered: true,
			Habitat:      "Tropical Forest",
			Region:       "Asia",
		},
		{
			ID:           "penguin-004",
			Name:         "Penny",
			Species:      "Humboldt Penguin",
			DateOfBirth:  time.Date(2020, time.November, 30, 0, 0, 0, 0, time.UTC),
			Diet:         "Carnivore",
			Weight:       4.2,
			IsEndangered: false,
			Habitat:      "Coastal",
			Region:       "South America",
		},
		{
			ID:           "giraffe-005",
			Name:         "George",
			Species:      "Giraffe",
			DateOfBirth:  time.Date(2016, time.April, 12, 0, 0, 0, 0, time.UTC),
			Diet:         "Herbivore",
			Weight:       1100.0,
			IsEndangered: false,
			Habitat:      "Savanna",
			Region:       "Africa",
		},
	}

	// Configure genstruct
	config := genstruct.Config{
		PackageName:      "main",                      // Target package name
		TypeName:         "Animal",                    // The struct type name
		ConstantIdent:    "Animal",                    // Prefix for constants
		VarPrefix:        "Animal",                    // Prefix for variables
		OutputFile:       "zoo_animals.go",            // Output file name
		IdentifierFields: []string{"Name", "Species"}, // Fields to use for naming variables
	}

	// Create a new generator with our config and animal data
	generator, err := genstruct.NewGenerator(config, animals)
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
	content, err := os.ReadFile("zoo_animals.go")
	if err != nil {
		fmt.Printf("Error reading generated file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nContents of generated file:")
	fmt.Println("---------------------------")
	fmt.Println(string(content))

	fmt.Println("\nTo use the generated code in your application you would:")
	fmt.Println("1. Import the generated file in your code by it's package name")
	fmt.Println("2. Use zoo.AnimalLeo, zoo.AnimalEllie, etc. to access specific animals")
	fmt.Println("3. Use zoo.AllAnimals slice for filtering and analysis")
}
