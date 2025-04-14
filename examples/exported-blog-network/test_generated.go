// This file isn't part of the actual package - it has a build tag to exclude it
// but shows how generated code would be used
//go:build ignore
package main

import (
	"fmt"

	"github.com/conneroisu/genstruct/examples/exported-blog-network/out"
	"github.com/conneroisu/genstruct/examples/exported-blog-network/pkg"
)

func main() {
	// Use the generated variables
	post := out.PostGettingStartedWithGo
	fmt.Printf("Post title: %s\n", post.Title)
	
	// Access embedded struct fields
	fmt.Printf("Slug: %s\n", post.Slug)
	fmt.Printf("Description: %s\n", post.Description)
	
	// Access all posts
	fmt.Printf("Total posts: %d\n", len(out.AllPosts))
	
	// Access all tags
	fmt.Printf("Total tags: %d\n", len(out.AllTags))
	
	// Use the embedded struct from the package
	var embedded pkg.Embedded
	embedded.Title = "New Post"
	embedded.Slug = "new-post"
	embedded.Description = "A new post created programmatically"
	fmt.Printf("Created new post: %s\n", embedded.Title)
	
	// This proves the code compiles and the structs can be used without "undefined: Embedded" errors
}