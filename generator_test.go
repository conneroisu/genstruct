package genstruct

import (
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

// Tag is a test struct for reference embedding
type Tag struct {
	ID   string
	Name string
	Slug string
}

// Post is a test struct that references Tags
type Post struct {
	ID       string
	Title    string
	Date     time.Time
	TagSlugs []string
	Tags     []*Tag `structgen:"TagSlugs"`
}

func TestStructReferenceEmbedding(t *testing.T) {
	// Create test data
	tags := []Tag{
		{ID: "tag-1", Name: "Go", Slug: "go"},
		{ID: "tag-2", Name: "Programming", Slug: "programming"},
		{ID: "tag-3", Name: "Testing", Slug: "testing"},
	}

	posts := []Post{
		{
			ID:       "post-1",
			Title:    "Testing in Go",
			Date:     time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
			TagSlugs: []string{"go", "testing"},
		},
		{
			ID:       "post-2",
			Title:    "Go Programming",
			Date:     time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
			TagSlugs: []string{"go", "programming"},
		},
	}

	// Create config
	config := Config{
		PackageName:      "testdata",
		TypeName:         "Post",
		ConstantIdent:    "Post",
		VarPrefix:        "Post",
		OutputFile:       "test_posts.go",
		IdentifierFields: []string{"Slug", "ID"},
	}

	// Create generator with references
	generator, err := NewGenerator(config, posts, tags)
	if err != nil {
		t.Fatalf("Error creating generator: %v", err)
	}

	// Make sure Refs map is correctly populated
	if len(generator.Refs) != 1 {
		t.Errorf("Expected 1 reference type, got %d", len(generator.Refs))
	}

	tagRef, ok := generator.Refs["Tag"]
	if !ok {
		t.Fatal("Expected to find Tag in references")
	}

	// Verify the generator stored the correct reference data
	tagRefValue := reflect.ValueOf(tagRef)
	if tagRefValue.Len() != 3 {
		t.Errorf("Expected 3 tags in reference data, got %d", tagRefValue.Len())
	}

	// Generate the code
	err = generator.Generate()
	if err != nil {
		t.Fatalf("Error generating code: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile("test_posts.go")
	if err != nil {
		t.Fatalf("Error reading generated file: %v", err)
	}

	contentStr := string(content)

	// Test that the references are properly populated via variable references
	expectedRefs := []string{
		// Verify use of Tag variables
		"Tags:",
		"[]*Tag{&Tag",
	}

	for _, expected := range expectedRefs {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected to find %q in generated code", expected)
		}
	}

	// Clean up
	err = os.Remove("test_posts.go")
	if err != nil {
		return
	}
}

// TestErrorHandling tests that errors are properly propagated
func TestErrorHandling(t *testing.T) {
	// Test with non-slice data
	nonSliceData := "not a slice"
	_, err := NewGenerator(Config{PackageName: "test"}, nonSliceData)
	if err == nil {
		t.Error("Expected error for non-slice data, got nil")
	}

	// Test with empty slice
	emptySlice := []string{}
	_, err = NewGenerator(Config{PackageName: "test"}, emptySlice)
	if err == nil {
		t.Error("Expected error for empty slice, got nil")
	}

	// Test with non-struct slice
	stringSlice := []string{"not", "a", "struct"}
	_, err = NewGenerator(Config{PackageName: "test"}, stringSlice)
	if err == nil {
		t.Error("Expected error for non-struct slice, got nil")
	}
}

// TestConfigInference tests that configuration values are properly inferred
func TestConfigInference(t *testing.T) {
	// Create test data
	type Person struct {
		ID   string
		Name string
		Age  int
	}

	people := []Person{
		{ID: "person-1", Name: "Alice", Age: 30},
		{ID: "person-2", Name: "Bob", Age: 25},
	}

	// Test with minimal configuration
	minimalConfig := Config{
		// Only specify package name, let everything else be inferred
		PackageName: "testdata",
	}

	generator, err := NewGenerator(minimalConfig, people)
	if err != nil {
		t.Fatalf("Error creating generator with minimal config: %v", err)
	}

	// Check that values were properly inferred
	if generator.Config.TypeName != "Person" {
		t.Errorf("Expected TypeName to be 'Person', got %q", generator.Config.TypeName)
	}

	if generator.Config.ConstantIdent != "Person" {
		t.Errorf("Expected ConstantIdent to be 'Person', got %q", generator.Config.ConstantIdent)
	}

	if generator.Config.VarPrefix != "Person" {
		t.Errorf("Expected VarPrefix to be 'Person', got %q", generator.Config.VarPrefix)
	}

	if generator.Config.OutputFile != "person_generated.go" {
		t.Errorf("Expected OutputFile to be 'person_generated.go', got %q", generator.Config.OutputFile)
	}

	if len(generator.Config.IdentifierFields) == 0 {
		t.Error("Expected IdentifierFields to be set with defaults")
	}

	// Test that specified values are not overridden
	customConfig := Config{
		PackageName:      "custom",
		TypeName:         "CustomPerson",
		ConstantIdent:    "CPerson",
		VarPrefix:        "Person",
		OutputFile:       "custom_output.go",
		IdentifierFields: []string{"Name", "ID"},
	}

	customGenerator, err := NewGenerator(customConfig, people)
	if err != nil {
		t.Fatalf("Error creating generator with custom config: %v", err)
	}

	if customGenerator.Config.TypeName != "CustomPerson" {
		t.Errorf("Expected TypeName to be 'CustomPerson', got %q", customGenerator.Config.TypeName)
	}

	if customGenerator.Config.ConstantIdent != "CPerson" {
		t.Errorf("Expected ConstantIdent to be 'CPerson', got %q", customGenerator.Config.ConstantIdent)
	}

	if customGenerator.Config.OutputFile != "custom_output.go" {
		t.Errorf("Expected OutputFile to be 'custom_output.go', got %q", customGenerator.Config.OutputFile)
	}
}

