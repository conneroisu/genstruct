package main

import (
	"os"
	"strings"
	"testing"
)

func TestBlogPostsGeneration(t *testing.T) {
	// Run the generation
	err := generateBlogData()
	if err != nil {
		t.Fatalf("Error generating blog posts: %v", err)
	}

	// Read the generated file
	content, err := os.ReadFile("blog_posts.go")
	if err != nil {
		t.Fatalf("Error reading generated file: %v", err)
	}

	contentStr := string(content)

	// Test that generated code contains expected elements
	expectedTests := []struct {
		name     string
		expected string
		message  string
	}{
		{
			name:     "Post constants",
			expected: "PostIntroductionToGoID",
			message:  "Should contain post ID constants",
		},
		{
			name:     "Post variables",
			expected: "var PostIntroductionToGo = Post{",
			message:  "Should contain post variables",
		},
		{
			name:     "Post slice",
			expected: "var AllPosts = []Post{",
			message:  "Should contain AllPosts slice",
		},
		{
			name:     "Has first post with tags",
			expected: "Tags: []Tag{",
			message:  "Should contain Tags field with array",
		},
		{
			name:     "Tag reference values",
			expected: "Name: \"Go Programming\"",
			message:  "Should contain referenced tag values",
		},
		{
			name:     "Second tag reference values",
			expected: "Name: \"Tutorials\"",
			message:  "Should contain all referenced tag values",
		},
		{
			name:     "Complex post reference",
			expected: "CodeGeneration",
			message:  "Should contain values from posts with multiple tags",
		},
	}

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
	os.Remove("blog_posts.go")
}