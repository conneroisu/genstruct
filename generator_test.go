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
	Tags     []Tag `structgen:"TagSlugs"`
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
	generator := NewGenerator(config, posts, tags)

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
	err := generator.Generate()
	if err != nil {
		t.Fatalf("Error generating code: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile("test_posts.go")
	if err != nil {
		t.Fatalf("Error reading generated file: %v", err)
	}

	contentStr := string(content)

	// Test that the references are properly populated
	expectedRefs := []string{
		// Fields from Tag structs
		"Name: \"Go\"",
		"Name: \"Testing\"",
		"Name: \"Programming\"",
		// Verify structs are properly referenced
		"Tags: []Tag{",
	}

	for _, expected := range expectedRefs {
		if !strings.Contains(contentStr, expected) {
			t.Errorf("Expected to find %q in generated code", expected)
		}
	}

	// Clean up
	os.Remove("test_posts.go")
}