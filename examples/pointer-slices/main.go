package main

import (
	"fmt"
	"os"
	"time"

	"github.com/conneroisu/genstruct"
)

// Author represents a content creator
type Author struct {
	ID       string    // Unique identifier
	Name     string    // Full name
	Email    string    // Contact email
	JoinDate time.Time // When they joined
}

// Comment represents user feedback
type Comment struct {
	ID        string    // Unique identifier
	Content   string    // Comment text
	AuthorID  string    // ID of the author
	Author    *Author   `structgen:"AuthorID"` // Populated from AuthorID - single pointer reference
	CreatedAt time.Time // When comment was created
}

// Article represents a blog article
type Article struct {
	ID          string     // Unique identifier
	Title       string     // Article title
	Slug        string     // URL-friendly slug
	Content     string     // Article content
	AuthorID    string     // ID of primary author
	Author      *Author    `structgen:"AuthorID"` // Single author (pointer reference)
	CommentIDs  []string   // IDs of comments
	Comments    []*Comment `structgen:"CommentIDs"` // List of pointer references to comments
	PublishedAt time.Time  // Publication date
}

// generateArticlesData generates the static data file
func generateArticlesData() error {
	// Define our array of authors as pointer slice
	authors := []*Author{
		{
			ID:       "author-001",
			Name:     "Alice Johnson",
			Email:    "alice@example.com",
			JoinDate: time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:       "author-002",
			Name:     "Bob Smith",
			Email:    "bob@example.com",
			JoinDate: time.Date(2021, time.March, 10, 0, 0, 0, 0, time.UTC),
		},
	}

	// Define our array of comments as pointer slice
	comments := []*Comment{
		{
			ID:        "comment-001",
			Content:   "Great article, very informative!",
			AuthorID:  "author-002", // Bob's comment
			CreatedAt: time.Date(2023, time.February, 5, 14, 30, 0, 0, time.UTC),
		},
		{
			ID:        "comment-002",
			Content:   "I learned a lot from this, thanks!",
			AuthorID:  "author-001", // Alice's comment
			CreatedAt: time.Date(2023, time.February, 6, 9, 15, 0, 0, time.UTC),
		},
		{
			ID:        "comment-003",
			Content:   "Could you expand on the second point?",
			AuthorID:  "author-002", // Bob's comment
			CreatedAt: time.Date(2023, time.February, 7, 17, 45, 0, 0, time.UTC),
		},
	}

	// Define our array of articles as pointer slice
	articles := []*Article{
		{
			ID:          "article-001",
			Title:       "Understanding Pointer Slices in Go",
			Slug:        "understanding-pointer-slices",
			Content:     "Pointers in Go provide a way to share data across your program...",
			AuthorID:    "author-001",                           // Written by Alice
			CommentIDs:  []string{"comment-001", "comment-003"}, // Has comments from Bob
			PublishedAt: time.Date(2023, time.February, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			ID:          "article-002",
			Title:       "Code Generation Techniques",
			Slug:        "code-generation-techniques",
			Content:     "Code generation can dramatically improve productivity...",
			AuthorID:    "author-002",            // Written by Bob
			CommentIDs:  []string{"comment-002"}, // Has a comment from Alice
			PublishedAt: time.Date(2023, time.March, 15, 0, 0, 0, 0, time.UTC),
		},
	}

	// Configure and generate the code
	genConfig := genstruct.Config{
		PackageName:      "main",
		OutputFile:       "articles_generated.go",
		IdentifierFields: []string{"ID", "Slug"},
	}
	gen, err := genstruct.NewGenerator(genConfig, articles, authors, comments)
	if err != nil {
		return fmt.Errorf("error creating generator: %w", err)
	}
	err = gen.Generate()
	if err != nil {
		return fmt.Errorf("error generating: %w", err)
	}

	return nil
}

func main() {
	// Generate the static data
	fmt.Println("Generating static article data with pointer references...")
	err := generateArticlesData()
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated static data in articles_generated.go")

	// Show the content of the generated file
	content, err := os.ReadFile("articles_generated.go")
	if err != nil {
		fmt.Printf("Error reading generated file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nContents of generated file:")
	fmt.Println("---------------------------")
	fmt.Println(string(content))

	fmt.Println("\nThis example demonstrates:")
	fmt.Println("1. Single pointer references (*Author)")
	fmt.Println("2. Slices of pointer references ([]*Comment)")
	fmt.Println("3. Bi-directional references between different types")
}
