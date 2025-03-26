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
    
    // Configure genstruct
    config := genstruct.Config{
        PackageName:      "zoo",         // Target package name
        TypeName:         "Animal",      // The struct type name
        ConstantIdent:    "Animal",      // Prefix for constants
        VarPrefix:        "Animal",      // Prefix for variables
        OutputFile:       "animals.go",  // Output file name
        IdentifierFields: []string{"Name", "Species"}, // Fields used for naming
    }
    
    // Create generator
    generator := genstruct.NewGenerator(config, animals)
    
    // Generate the code
    err := generator.Generate()
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
    
    // Access all animals
    for _, animal := range zoo.AllAnimals {
        fmt.Printf("%s: %s\n", animal.Name, animal.Species)
    }
}
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
    Tags     []Tag     `structgen:"TagSlugs"` // Will be populated from TagSlugs
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

## Common Configuration Options

Here are some common configuration options:

| Option | Description | Example |
|--------|-------------|---------|
| PackageName | Target package name | `"models"` |
| TypeName | Struct type name | `"User"` |
| ConstantIdent | Prefix for constants | `"User"` |
| VarPrefix | Prefix for variables | `"User"` |
| OutputFile | Output file path | `"users.go"` |
| IdentifierFields | Priority fields for naming | `[]string{"ID", "Username"}` |

For more examples, see the [examples directory](https://github.com/conneroisu/genstruct/tree/main/examples) in the repository.