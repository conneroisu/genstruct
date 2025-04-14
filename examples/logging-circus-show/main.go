// Package main shows how to use the generator with a custom logger.
package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/conneroisu/genstruct"
)

// Trick represents a circus trick that can be performed
type Trick struct {
	ID          string
	Name        string
	Description string
	Difficulty  int    // 1-10 scale
	Equipment   string // Equipment needed for this trick
}

// Performer represents a circus performer
type Performer struct {
	ID         string
	Name       string
	Role       string
	Experience int      // Years of experience
	Specialty  string   // Main specialty
	TrickIDs   []string `structgen:"Trick"` // This will be used to reference tricks
}

func main() {
	// Initialize the logger (this will parse flags)
	logger := genstruct.InitLogger()
	logger.Info("Starting circus show example")

	// Define some tricks
	tricks := []Trick{
		{
			ID:          "trick-fire-breathing",
			Name:        "Fire Breathing",
			Description: "Performer breathes out fire in a controlled manner",
			Difficulty:  8,
			Equipment:   "Fire fuel, torch",
		},
		{
			ID:          "trick-triple-somersault",
			Name:        "Triple Somersault",
			Description: "Three consecutive aerial somersaults",
			Difficulty:  9,
			Equipment:   "Trapeze, safety net",
		},
		{
			ID:          "trick-juggling-knives",
			Name:        "Knife Juggling",
			Description: "Juggling sharp knives in complex patterns",
			Difficulty:  7,
			Equipment:   "Balanced throwing knives",
		},
		{
			ID:          "trick-sword-swallowing",
			Name:        "Sword Swallowing",
			Description: "Inserting a sword down the throat and into the esophagus",
			Difficulty:  8,
			Equipment:   "Special performance sword",
		},
		{
			ID:          "trick-tightrope-walk",
			Name:        "Tightrope Walking",
			Description: "Walking across a thin wire high above the ground",
			Difficulty:  6,
			Equipment:   "Tightrope, balance pole",
		},
	}

	// Define some performers with their tricks
	performers := []Performer{
		{
			ID:         "performer-alex",
			Name:       "Alex the Great",
			Role:       "Fire Performer",
			Experience: 15,
			Specialty:  "Fire manipulation",
			TrickIDs:   []string{"trick-fire-breathing"},
		},
		{
			ID:         "performer-bella",
			Name:       "Bella Air",
			Role:       "Acrobat",
			Experience: 8,
			Specialty:  "Aerial acrobatics",
			TrickIDs:   []string{"trick-triple-somersault", "trick-tightrope-walk"},
		},
		{
			ID:         "performer-carlos",
			Name:       "Carlos Danger",
			Role:       "Knife Expert",
			Experience: 12,
			Specialty:  "Dangerous props",
			TrickIDs:   []string{"trick-juggling-knives", "trick-sword-swallowing"},
		},
	}

	// Create a generator with functional options
	logger.Info("Creating generator with functional options")
	generator := genstruct.NewGenerator(
		genstruct.WithPackageName("main"),
		genstruct.WithOutputFile("circus_generated.go"),
		genstruct.WithVarPrefix("Circus"),
		genstruct.WithConstantIdent("Circus"),
		genstruct.WithLogger(createCustomLogger()), // Use our custom logger for demonstration
	)

	// Generate the code, passing performers and tricks
	logger.Info("Generating code", "output", "circus_generated.go")
	err := generator.Generate(performers, tricks)
	if err != nil {
		logger.Error("Failed to generate code", "error", err)
		return
	}

	logger.Info("Code generation completed successfully")
	fmt.Println("Generated circus show code successfully. See circus_generated.go")
}

// createCustomLogger creates a logger with a specific format for this example
func createCustomLogger() *slog.Logger {
	// Get the current verbosity level from flags

	// Determine log level
	var level = slog.LevelDebug

	// Create a handler with custom attributes for the circus show
	h := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}).WithAttrs([]slog.Attr{
		slog.String("component", "circus-generator"),
	})

	return slog.New(h)
}
