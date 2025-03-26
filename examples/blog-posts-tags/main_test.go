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
	content, err := os.ReadFile("blog_generated.go")
	if err != nil {
		t.Fatalf("Error reading blog_generated.go file: %v", err)
	}

	contentStr := string(content)
	
	// For compatibility with test, derive blog post and tags content from the single file
	blogPostStr := contentStr
	blogTagsStr := contentStr

	// Test that generated code contains expected elements
	expectedBlogPostTests := []struct {
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
			expected: "Tags:     []Tag{Tag",
			message:  "Should contain Tags field with array",
		},
		{
			name:     "Complex post reference",
			expected: "CodeGeneration",
			message:  "Should contain values from posts with multiple tags",
		},
	}

	// Test the blog_posts.go file
	for _, tc := range expectedBlogPostTests {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(blogPostStr, tc.expected) {
				t.Errorf("%s: %q not found in generated posts code", tc.message, tc.expected)
			}
		})
	}
	
	// Test the blog_tags.go file for tag values
	expectedBlogTagsTests := []struct {
		name     string
		expected string
		message  string
	}{
		{
			name:     "Tag reference values",
			expected: "Name: \"Go Programming\"",
			message:  "Should contain tag values",
		},
		{
			name:     "Second tag reference values",
			expected: "Name: \"Tutorials\"",
			message:  "Should contain all tag values",
		},
	}
	
	// Test the blog_tags.go file
	for _, tc := range expectedBlogTagsTests {
		t.Run(tc.name, func(t *testing.T) {
			if !strings.Contains(blogTagsStr, tc.expected) {
				t.Errorf("%s: %q not found in generated tags code", tc.message, tc.expected)
			}
		})
	}
}

func TestCleanup(t *testing.T) {
	// This runs after all tests to clean up
	os.Remove("blog_generated.go")
}