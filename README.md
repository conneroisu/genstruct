# genstruct

A Go library for generating statically defined Go structs with support for references between structures.

## Features

- Generate Go code for static structs and arrays
- Automatic creation of constants, variables, and slices
- Smart naming of variables based on identifier fields
- Reference embedding via struct tags to connect related structs
- Customizable code generation

## Installation

```bash
go get github.com/conneroisu/genstruct
```

## Basic Usage

```go
// Define your struct type
type Animal struct {
    ID     string
    Name   string
    Species string
    Diet   string
}

// Create a slice of structs
animals := []Animal{
    {ID: "lion-001", Name: "Leo", Species: "Lion", Diet: "Carnivore"},
    {ID: "tiger-001", Name: "Stripes", Species: "Tiger", Diet: "Carnivore"},
}

// Configure the generator
config := genstruct.Config{
    PackageName:   "zoo",          // Target package name
    TypeName:      "Animal",       // Struct type name
    ConstantIdent: "Animal",       // Prefix for constants
    VarPrefix:     "Animal",       // Prefix for variables
    OutputFile:    "animals.go",   // Output file path
}

// Generate the code
generator := genstruct.NewGenerator(config, animals)
err := generator.Generate()
```

## Struct Reference Embedding

A powerful feature of genstruct is the ability to automatically populate fields in one struct by referencing values from another struct.

### How it works

1. Define your structs with reference fields
2. Use the `structgen` tag to specify the source field
3. Pass additional reference datasets to `NewGenerator`

### Example

```go
// Define your tag struct
type Tag struct {
    ID   string
    Name string
    Slug string
}

// Define your post struct with references to tags
type Post struct {
    ID       string
    Title    string
    TagSlugs []string  // Contains tag slugs
    Tags     []Tag     `structgen:"TagSlugs"` // Will be populated from TagSlugs
}

// Create your data
tags := []Tag{
    {ID: "tag-001", Name: "Go", Slug: "go"},
    {ID: "tag-002", Name: "Programming", Slug: "programming"},
}

posts := []Post{
    {
        ID: "post-001", 
        Title: "Introduction to Go",
        TagSlugs: []string{"go", "programming"},
    },
}

// Generate code with both datasets
generator := genstruct.NewGenerator(config, posts, tags)
err := generator.Generate()
```

The generated code will include:
1. Constant definitions for all Post IDs
2. Variable definitions for each Post
3. A slice containing all Posts 
4. **Constant definitions for all Tag IDs**
5. **Variable definitions for each Tag**
6. **A slice containing all Tags**
7. **Cross-references between Posts and Tags** (the `Tags` field in each Post will reference the generated Tag variables)

All of this is generated in a single file, with a single generator call.

## Config Options

- `PackageName`: The package name for the generated file
- `TypeName`: The name of the struct type
- `ConstantIdent`: Prefix for generated constants
- `VarPrefix`: Prefix for generated variables
- `OutputFile`: The output file path
- `IdentifierFields`: Fields to use for naming (default: "ID", "Name", "Slug", "Title", "Key", "Code")
- `CustomVarNameFn`: Optional function to customize variable naming

## Dependencies

- [jennifer](https://github.com/dave/jennifer) for code generation