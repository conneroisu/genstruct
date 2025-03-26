package main

import (
	"fmt"
	"os"
	"time"

	"github.com/conneroisu/genstruct"
)

// Tag represents a blog post tag
type Tag struct {
	ID       string   // Unique identifier for the tag
	Name     string   // Name of the tag
	Slug     string   // URL-friendly slug for the tag
	RelatedTagSlugs []string // List of related tag slugs (used one direction only to avoid circular references)
	RelatedTags []*Tag   `structgen:"RelatedTagSlugs"` // Populated from RelatedTagSlugs
}

// Post represents a blog post
type Post struct {
	ID       string    // Unique identifier for the post
	Title    string    // Title of the post
	Slug     string    // URL-friendly slug for the post
	Content  string    // Content of the post
	TagSlugs []string  // List of tag slugs
	Tags     []*Tag    `structgen:"TagSlugs"` // Populated from TagSlugs using pointers
	Date     time.Time // Publication date
}

// generateBlogData generates the static blog data file
func generateBlogData() error {
	// Define our array of tag data
	tags := []Tag{
		{
			ID:       "tag-001",
			Name:     "Go Programming",
			Slug:     "go-programming",
			RelatedTagSlugs: []string{"code-generation", "developer-tools"},
		},
		{
			ID:       "tag-002",
			Name:     "Code Generation",
			Slug:     "code-generation",
			RelatedTagSlugs: []string{"developer-tools"},
		},
		{
			ID:       "tag-003",
			Name:     "Tutorials",
			Slug:     "tutorials",
			RelatedTagSlugs: []string{"developer-tools"},
		},
		{
			ID:       "tag-004",
			Name:     "Developer Tools",
			Slug:     "developer-tools",
			RelatedTagSlugs: []string{},
		},
	}

	// Define our array of post data
	posts := []Post{
		{
			ID:       "post-001",
			Title:    "Introduction to Go",
			Slug:     "introduction-to-go",
			Content:  "Go is a statically typed programming language...",
			TagSlugs: []string{"go-programming", "tutorials"},
			Date:     time.Date(2023, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       "post-002",
			Title:    "Code Generation in Go",
			Slug:     "code-generation-in-go",
			Content:  "Code generation can save time and reduce errors...",
			TagSlugs: []string{"go-programming", "code-generation", "developer-tools"},
			Date:     time.Date(2023, time.February, 20, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       "post-003",
			Title:    "Building Developer Tools",
			Slug:     "building-developer-tools",
			Content:  "Developer tools can greatly enhance productivity...",
			TagSlugs: []string{"developer-tools", "tutorials"},
			Date:     time.Date(2023, time.March, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	// Configure and generate for both post and tag data in one step
	genConfig := genstruct.Config{
		PackageName:      "main",
		OutputFile:       "blog_generated.go",
		IdentifierFields: []string{"Slug", "ID"},
	}
	gen, err := genstruct.NewGenerator(genConfig, posts, tags)
	if err != nil {
		return fmt.Errorf("error creating generator: %w", err)
	}
	err = gen.Generate()
	if err != nil {
		return fmt.Errorf("error generating: %w", err)
	}

	// Return nil as we've handled the generation
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

	fmt.Println("Successfully generated static blog post data in blog_generated.go")

	// Show the content of the generated file
	content, err := os.ReadFile("blog_generated.go")
	if err != nil {
		fmt.Printf("Error reading generated file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nContents of generated file:")
	fmt.Println("---------------------------")
	fmt.Println(string(content))

	fmt.Println("\nTo use the generated code in your application you would:")
	fmt.Println("1. Import the generated file in your code by its package name")
	fmt.Println("2. Use main.PostIntroductionToGo, main.PostCodeGenerationInGo, etc. to access specific posts")
	fmt.Println("3. Use main.AllPosts slice for filtering and analysis")
	fmt.Println("4. The Tags field in each post will be populated with pointers to Tag objects referenced by slug")
	fmt.Println("5. Similarly, the RelatedTags field in each Tag will be populated with pointers to related Tag objects")
}