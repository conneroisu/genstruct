// Package main demonstrates how to export blog post data with tag references.
//
// The generated code will be in the "./out" package and will maintain references
// between posts and tags using pointer slices.
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/genstruct"
	"github.com/conneroisu/genstruct/examples/exported-blog-posts-tags/pkg"
)

// generateBlogData generates the static blog data file in the out package
func generateBlogData() error {
	// Define our blog data in the pkg package
	posts := pkg.Posts
	tags := pkg.Tags

	// Configure and generate for both post and tag data in one step
	// Use absolute path for the output file
	outPath := "/home/connerohnesorge/Documents/001Repos/genstruct/examples/exported-blog-posts-tags/out/blog_generated.go"
	gen := genstruct.NewGenerator(
		genstruct.WithOutputFile(outPath),
		genstruct.WithIdentifierFields([]string{"Slug", "ID"}),
	)

	// Generate code - this will maintain references between posts and tags
	err := gen.Generate(posts, tags)
	if err != nil {
		return fmt.Errorf("error generating: %w", err)
	}

	return nil
}

func main() {
	// Generate the blog post data
	fmt.Println("Generating static blog post data...")
	err := generateBlogData()
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated static blog post data in ./out/blog_generated.go")

	// Show the content of the generated file
	outPath := "/home/connerohnesorge/Documents/001Repos/genstruct/examples/exported-blog-posts-tags/out/blog_generated.go"
	content, err := os.ReadFile(outPath)
	if err != nil {
		fmt.Printf("Error reading generated file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nContents of generated file:")
	fmt.Println("---------------------------")
	fmt.Println(string(content))

	fmt.Println("\nTo use the generated code in your application you would:")
	fmt.Println("1. Import the generated file with \"github.com/conneroisu/genstruct/examples/exported-blog-posts-tags/out\"")
	fmt.Println("2. Use out.PostIntroductionToGo, out.PostCodeGenerationInGo, etc. to access specific posts")
	fmt.Println("3. Use out.AllPosts slice for filtering and analysis")
	fmt.Println("4. The Tags field in each post will be populated with pointers to Tag objects referenced by slug")
	fmt.Println("5. Similarly, the RelatedTags field in each Tag will be populated with pointers to related Tag objects")
}