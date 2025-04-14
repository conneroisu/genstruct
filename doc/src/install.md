# Installation and Getting Started

## Installation Options

### Using Go Modules (Recommended)

The recommended way to install genstruct is using Go modules:

```bash
# Initialize your module if you haven't already
go mod init example.com/myproject

# Add genstruct to your project
go get github.com/conneroisu/genstruct
```

### Manual Installation

You can also clone the repository directly:

```bash
git clone https://github.com/conneroisu/genstruct.git
cd genstruct
go install
```

## Quick Start Guide

Here's how to start using genstruct in your project:

### 1. Define Your Structs

Create the structs that you want to generate static data for:

```go
// Example: animal.go
package zoo

import "time"

type Animal struct {
    ID           string
    Name         string
    Species      string
    DateOfBirth  time.Time
    Diet         string
    Weight       float64
    IsEndangered bool
}
```

### 2. Create a Generator

Set up a generator in your code:

```go
package main

import (
    "github.com/conneroisu/genstruct"
    "time"
)

func main() {
    // Define sample data
    animals := []Animal{
        {
            ID:           "lion-001",
            Name:         "Leo",
            Species:      "African Lion",
            DateOfBirth:  time.Date(2018, time.March, 15, 0, 0, 0, 0, time.UTC),
            Diet:         "Carnivore",
            Weight:       180.5,
            IsEndangered: true,
        },
        // Add more animals...
    }
    
    // Configure genstruct (minimal configuration - many values are inferred automatically)
    config := genstruct.Config{
        PackageName: "zoo",        // Target package name
        OutputFile:  "animals.go", // Output file name
        // TypeName, ConstantIdent, and VarPrefix will be inferred as "Animal"
        // Customize which fields to prioritize for naming
        IdentifierFields: []string{"Name", "Species"},
    }
    
    // Create generator
    generator, err := genstruct.NewGenerator(config, animals)
    if err != nil {
        panic(err)
    }
    
    // Generate the code
    err = generator.Generate()
    if err != nil {
        panic(err)
    }
}
```

### 3. Run Your Generator

Execute your generator:

```bash
go run main.go
```

This will create `animals.go` with your static data.

### 4. Use the Generated Code

Import and use the generated code:

```go
package main

import (
    "fmt"
    "yourmodule/zoo" // Import the generated package
)

func main() {
    // Access specific animal
    fmt.Println(zoo.AnimalLeo.Species) // Prints: African Lion
    
    // Access all animals (which are pointers)
    for _, animal := range zoo.AllAnimals {
        fmt.Printf("%s: %s\n", animal.Name, animal.Species)
    }
    
    // Direct access to an animal from the slice
    fmt.Println(zoo.AllAnimals[0].Name) // Prints: Leo
}
```

The generated code will create a slice of pointers:

```go
var AllAnimals = []*Animal{&AnimalLeo, &AnimalEllie, &AnimalStripes, ...}
```

## Advanced: Using Struct References

For more complex data structures with relationships:

```go
// Define related structs
type Tag struct {
    ID   string
    Name string
    Slug string
}

type Post struct {
    ID       string
    Title    string
    TagSlugs []string  // References to tags
    Tags     []*Tag    `structgen:"TagSlugs"` // Will be populated from TagSlugs
}

// Create datasets
tags := []Tag{
    {ID: "tag-1", Name: "Go", Slug: "go"},
    // More tags...
}

posts := []Post{
    {
        ID:       "post-1",
        Title:    "Go Programming",
        TagSlugs: []string{"go"},
    },
    // More posts...
}

// Create generator with both datasets
generator := genstruct.NewGenerator(config, posts, tags)

// Generate code with relationships
err := generator.Generate()
```

## Configuration Options

Configuration is simplified with auto-inference of many values. Here are the available options:

| Option | Description | Default | Example |
|--------|-------------|---------|---------|
| PackageName | Target package name | `"{output-directory}"` | `"models"` |
| TypeName | Struct type name | *Inferred from data* | `"User"` |
| ConstantIdent | Prefix for constants | *Same as TypeName* | `"User"` |
| VarPrefix | Prefix for variables | *Same as TypeName* | `"User"` |
| OutputFile | Output file path | *typename_generated.go* | `"users.go"` |
| IdentifierFields | Priority fields for naming | `[]string{"ID", "Name", "Slug", "Title", "Key", "Code"}` | `[]string{"ID", "Username"}` |
| CustomVarNameFn | Custom naming function | *None* | *Custom function* |

In many cases, you only need to specify the `PackageName` and `OutputFile` - everything else will be inferred automatically:

```go
// Minimal configuration example
config := genstruct.Config{
    PackageName: "myapp",
    OutputFile:  "users_generated.go",
}

// The TypeName will be inferred as "User" from the data
generator, err := genstruct.NewGenerator(config, users)
if err != nil {
    // handle error
}
```

For more examples, see the [examples directory](https://github.com/conneroisu/genstruct/tree/main/examples) in the repository.
