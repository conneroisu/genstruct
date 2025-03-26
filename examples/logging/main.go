package main

import (
	"flag"
	"fmt"

	"github.com/conneroisu/genstruct"
)

// Animal represents an animal in our zoo
type Animal struct {
	Name        string
	Species     string
	Description string
	Age         int
}

func main() {
	// Initialize the logger (this will parse flags)
	logger := genstruct.InitLogger()
	logger.Info("Starting logging example")

	// Create some data - a slice of Animals
	animals := []Animal{
		{
			Name:        "Leo",
			Species:     "Lion",
			Description: "King of the jungle",
			Age:         5,
		},
		{
			Name:        "Jumbo",
			Species:     "Elephant",
			Description: "Largest land mammal",
			Age:         10,
		},
		{
			Name:        "Stripes",
			Species:     "Zebra",
			Description: "Has black and white stripes",
			Age:         3,
		},
	}

	// Create a generator config
	config := genstruct.Config{
		PackageName:  "zoo",
		OutputFile:   "animals_generated.go",
		VarPrefix:    "Animal",
		ConstantIdent: "Animal",
		// Logger is automatically set if not provided
	}

	// Create a new generator with our config and data
	generator, err := genstruct.NewGenerator(config, animals)
	if err != nil {
		logger.Error("Failed to create generator", "error", err)
		return
	}

	// Generate the code
	logger.Info("Generating code", "output", config.OutputFile)
	err = generator.Generate()
	if err != nil {
		logger.Error("Failed to generate code", "error", err)
		return
	}

	logger.Info("Code generation completed successfully")
	fmt.Println("Generated code successfully. See animals_generated.go")
}