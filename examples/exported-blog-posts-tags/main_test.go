package main

import (
	"os"
	"strings"
	"testing"
)

func TestExportedBlogPostsGeneration(t *testing.T) {
	// Ensure the output directory exists
	outDir := "/home/connerohnesorge/Documents/001Repos/genstruct/examples/exported-blog-posts-tags/out"
	err := os.MkdirAll(outDir, 0755)
	if err != nil {
		t.Fatalf("Error creating output directory: %v", err)
	}

	// Run the generation
	err = generateBlogData()
	if err != nil {
		t.Fatalf("Error generating blog posts: %v", err)
	}

	// Read the generated file
	outPath := "/home/connerohnesorge/Documents/001Repos/genstruct/examples/exported-blog-posts-tags/out/blog_generated.go"
	content, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("Error reading blog_generated.go file: %v", err)
	}

	contentStr := string(content)

	// Test that generated code contains expected elements
	expectedTests := []struct {
		name     string
		expected string
		message  string
	}{
		{
			name:     "Package name",
			expected: "package out",
			message:  "Should have correct package name",
		},
		{
			name:     "Import pkg",
			expected: "github.com/conneroisu/genstruct/examples/exported-blog-posts-tags/pkg",
			message:  "Should import the pkg package for type definitions",
		},
		{
			name:     "Post constants",
			expected: "PostIntroductionToGoID",
			message:  "Should contain post ID constants",
		},
		{
			name:     "Post variables",
			expected: "var PostIntroductionToGo = pkg.Post{",
			message:  "Should contain post variables with pkg qualification",
		},
		{
			name:     "Post slice",
			expected: "var AllPosts = []*pkg.Post{",
			message:  "Should contain AllPosts pointer slice with pkg qualification",
		},
		{
			name:     "Has first post with pointer tags",
			expected: "Tags:     []*pkg.Tag{&",
			message:  "Should contain Tags field with pointer array and pkg qualification",
		},
		{
			name:     "Tag constants",
			expected: "TagGoProgrammingID",
			message:  "Should contain tag ID constants",
		},
		{
			name:     "Tag variables",
			expected: "var TagGoProgramming = pkg.Tag{",
			message:  "Should contain tag variables with pkg qualification",
		},
		{
			name:     "Tag slice",
			expected: "var AllTags = []*pkg.Tag{",
			message:  "Should contain AllTags pointer slice with pkg qualification",
		},
		{
			name:     "Related tags pointer field",
			expected: "RelatedTags:     []*pkg.Tag{",
			message:  "Should contain RelatedTags pointer field with pkg qualification",
		},
	}

	// Test the generated file
	for _, tc := range expectedTests {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(contentStr, tc.expected) {
				t.Errorf("%s: %q not found in generated code", tc.message, tc.expected)
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	// This runs after all tests to clean up
	// Comment out during development if you want to inspect the generated file
	outPath := "/home/connerohnesorge/Documents/001Repos/genstruct/examples/exported-blog-posts-tags/out/blog_generated.go"
	err := os.Remove(outPath)
	if err != nil && !os.IsNotExist(err) {
		t.Fatalf("Error removing blog_generated.go file: %v", err)
	}
}