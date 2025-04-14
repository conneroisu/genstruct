package pkg

import (
	"time"
)

// Tag represents a blog post tag
type Tag struct {
	ID              string   // Unique identifier for the tag
	Name            string   // Name of the tag
	Slug            string   // URL-friendly slug for the tag
	RelatedTagSlugs []string // List of related tag slugs (used one direction only to avoid circular references)
	RelatedTags     []*Tag   `structgen:"RelatedTagSlugs"` // Populated from RelatedTagSlugs
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

// Tags is the slice of sample tags
var Tags = []Tag{
	{
		ID:              "tag-001",
		Name:            "Go Programming",
		Slug:            "go-programming",
		RelatedTagSlugs: []string{"code-generation", "developer-tools"},
	},
	{
		ID:              "tag-002",
		Name:            "Code Generation",
		Slug:            "code-generation",
		RelatedTagSlugs: []string{"developer-tools"},
	},
	{
		ID:              "tag-003",
		Name:            "Tutorials",
		Slug:            "tutorials",
		RelatedTagSlugs: []string{"developer-tools"},
	},
	{
		ID:              "tag-004",
		Name:            "Developer Tools",
		Slug:            "developer-tools",
		RelatedTagSlugs: []string{},
	},
}

// Posts is the slice of sample posts
var Posts = []Post{
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