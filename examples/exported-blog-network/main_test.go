package main

import (
	"os"
	"testing"
)

func TestGenerateBlogNetwork(t *testing.T) {
	// Ensure the out directory exists
	err := os.MkdirAll("./out", 0755)
	if err != nil {
		t.Fatalf("Failed to create out directory: %v", err)
	}

	// Generate the blog network data
	err = generateBlogNetwork()
	if err != nil {
		t.Fatalf("Failed to generate blog network data: %v", err)
	}

	// Verify the file was created
	_, err = os.Stat("./out/blog_network.go")
	if err != nil {
		t.Fatalf("Generated file does not exist: %v", err)
	}
}