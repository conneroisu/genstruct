// Package main demonstrates how to export complex embedded structs.
//
// The generated code will be in another package (./out) and
// will be named after the folder where the source code is generated.
package main

import (
	"fmt"
	"os"

	"github.com/conneroisu/genstruct"
	"github.com/conneroisu/genstruct/examples/exported-blog-network/pkg"
)

// generateBlogNetwork generates the static blog network data file
func generateBlogNetwork() error {
	// Create a generator with functional options
	generator := genstruct.NewGenerator(
		genstruct.WithOutputFile("./out/blog_network.go"),
		genstruct.WithIdentifierFields([]string{"Slug", "Title"}),
	)

	// Generate the code, passing all our data collections
	return generator.Generate(pkg.Posts, pkg.Tags, pkg.Projects, pkg.Employments)
}

func main() {
	// Generate the blog network data
	fmt.Println("Generating blog network data...")
	err := generateBlogNetwork()
	if err != nil {
		fmt.Printf("Error generating code: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully generated blog network data in out/blog_network.go")

	// Show the content of the generated file
	content, err := os.ReadFile("./out/blog_network.go")
	if err != nil {
		fmt.Printf("Error reading generated file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nContents of generated file:")
	fmt.Println("---------------------------")
	fmt.Println(string(content))

	fmt.Println("\nTo use the generated code in your application you would:")
	fmt.Println("1. Import the generated file in your code by its package name (out)")
	fmt.Println("2. Use out.PostGettingStartedWithGo, out.TagGo, etc. to access specific items")
	fmt.Println("3. Use out.AllPosts, out.AllTags, etc. slices for filtering and analysis")
}
