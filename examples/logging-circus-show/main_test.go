package main

import (
	"bytes"
	"log/slog"
	"os"
	"testing"

	"github.com/conneroisu/genstruct"
)

// TestCircusGeneration tests the generation of circus code
func TestCircusGeneration(t *testing.T) {
	// Redirect logging to a buffer for testing
	var buf bytes.Buffer
	logger := slog.New(slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))

	// Define simple test data
	tricks := []Trick{
		{
			ID:          "trick-test",
			Name:        "Test Trick",
			Description: "A trick for testing",
			Difficulty:  5,
			Equipment:   "Test equipment",
		},
	}

	performers := []Performer{
		{
			ID:         "performer-test",
			Name:       "Test Performer",
			Role:       "Tester",
			Experience: 10,
			Specialty:  "Testing",
			TrickIDs:   []string{"trick-test"},
		},
	}

	// Create a temporary file for output
	tempFile, err := os.CreateTemp("", "circus_test_*.go")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Set up the config
	config := genstruct.Config{
		PackageName:   "circus_test",
		OutputFile:    tempFile.Name(),
		VarPrefix:     "Test",
		ConstantIdent: "Test",
		Logger:        logger,
	}

	// Create the generator
	generator, err := genstruct.NewGenerator(config, performers, tricks)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	// Generate the code
	err = generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Verify that the file exists and has content
	fileInfo, err := os.Stat(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	if fileInfo.Size() == 0 {
		t.Errorf("Generated file is empty")
	}

	expectedLogMsgs := []string{
		"Starting code generation",
		"Generating constants",
		"Generating variables",
		"Generating slice",
		"Processing reference datasets",
	}

	for _, msg := range expectedLogMsgs {
		if !bytes.Contains(buf.Bytes(), []byte(msg)) {
			t.Errorf("Log does not contain expected message: %s", msg)
		}
	}
}
