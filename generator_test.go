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

	// Create generator with functional options
	generator := NewGenerator(
		WithPackageName("testdata"),
		WithTypeName("Post"),
		WithConstantIdent("Post"),
		WithVarPrefix("Post"),
		WithOutputFile("test_posts.go"),
		WithIdentifierFields([]string{"Slug", "ID"}),
	)

	// Generate will populate Refs
	err := generator.Generate(posts, tags)
	if err != nil {
		t.Fatalf("Error generating code: %v", err)
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
	generator := NewGenerator(WithPackageName("test"))
	err := generator.Generate(nonSliceData)
	if err == nil {
		t.Error("Expected error for non-slice data, got nil")
	}

	// Test with empty slice
	emptySlice := []string{}
	generator = NewGenerator(WithPackageName("test"))
	err = generator.Generate(emptySlice)
	if err == nil {
		t.Error("Expected error for empty slice, got nil")
	}

	// Test with non-struct slice
	stringSlice := []string{"not", "a", "struct"}
	generator = NewGenerator(WithPackageName("test"))
	err = generator.Generate(stringSlice)
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
	minimalGenerator := NewGenerator(
		WithPackageName("testdata"),
	)
	
	// Try to infer values from data
	err := minimalGenerator.inferConfig(people)
	if err != nil {
		t.Fatalf("Error inferring config: %v", err)
	}

	// Check that values were properly inferred
	if minimalGenerator.TypeName != "Person" {
		t.Errorf("Expected TypeName to be 'Person', got %q", minimalGenerator.TypeName)
	}

	if minimalGenerator.ConstantIdent != "Person" {
		t.Errorf("Expected ConstantIdent to be 'Person', got %q", minimalGenerator.ConstantIdent)
	}

	if minimalGenerator.VarPrefix != "Person" {
		t.Errorf("Expected VarPrefix to be 'Person', got %q", minimalGenerator.VarPrefix)
	}

	if minimalGenerator.OutputFile != "person_generated.go" {
		t.Errorf("Expected OutputFile to be 'person_generated.go', got %q", minimalGenerator.OutputFile)
	}

	if len(minimalGenerator.IdentifierFields) == 0 {
		t.Error("Expected IdentifierFields to be set with defaults")
	}

	// Test that specified values are not overridden
	customGenerator := NewGenerator(
		WithPackageName("custom"),
		WithTypeName("CustomPerson"),
		WithConstantIdent("CPerson"),
		WithVarPrefix("Person"),
		WithOutputFile("custom_output.go"),
		WithIdentifierFields([]string{"Name", "ID"}),
	)
	
	// Try to infer values from data
	err = customGenerator.inferConfig(people)
	if err != nil {
		t.Fatalf("Error inferring config: %v", err)
	}

	if customGenerator.TypeName != "CustomPerson" {
		t.Errorf("Expected TypeName to be 'CustomPerson', got %q", customGenerator.TypeName)
	}

	if customGenerator.ConstantIdent != "CPerson" {
		t.Errorf("Expected ConstantIdent to be 'CPerson', got %q", customGenerator.ConstantIdent)
	}

	if customGenerator.OutputFile != "custom_output.go" {
		t.Errorf("Expected OutputFile to be 'custom_output.go', got %q", customGenerator.OutputFile)
	}
}

