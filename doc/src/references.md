# Struct Reference Embedding

One of the powerful features of `genstruct` is the ability to automatically populate fields in one struct with references to other structs. This allows you to define relationships between your data structures and have them automatically resolved during code generation.

## How It Works

References between structs are established using Go struct tags with the `structgen` tag:

```go
type Post struct {
    ID       string    
    TagSlugs []string       // Contains slugs of tags
    Tags     []*Tag    `structgen:"TagSlugs"` // Will be populated based on TagSlugs
}
```

The `structgen` tag tells `genstruct` that:

1. The value of this field (`Tags`) should be populated based on the value of another field (`TagSlugs`)
2. It should look up matching values in the reference data provided

When references are resolved, genstruct now generates proper variable references instead of duplicating structs:

```go
// Instead of duplicating the entire Tag struct:
Tags: []*Tag{&Tag{
    ID:   "tag-001",
    Name: "Go Programming",
    Slug: "go-programming",
}},

// It now uses references to pre-generated variables:
Tags: []*Tag{&TagGoProgramming, &TagTutorials},
```

This improves efficiency and maintainability of the generated code.

## Supported Reference Types

Currently, `genstruct` supports several types of references (with pointer-based references being the default and recommended approach):

1. **String to Struct**: A string field referencing a single struct
   ```go
   type Post struct {
       AuthorID string            // Contains an author ID
       Author   Author     `structgen:"AuthorID"` // Will be populated from AuthorID
   }
   ```

2. **String to Struct Pointer**: A string field referencing a struct pointer
   ```go
   type Post struct {
       AuthorID string            // Contains an author ID
       Author   *Author    `structgen:"AuthorID"` // Will be populated from AuthorID
   }
   ```

3. **String Slice to Struct Slice**: A slice of strings referencing a slice of structs
   ```go
   type Post struct {
       TagSlugs []string       // Contains slugs of tags
       Tags     []Tag     `structgen:"TagSlugs"` // Will be populated as a slice of Tags
   }
   ```

4. **String Slice to Struct Pointer Slice**: A slice of strings referencing a slice of struct pointers (recommended)
   ```go
   type Post struct {
       TagSlugs []string       // Contains slugs of tags
       Tags     []*Tag    `structgen:"TagSlugs"` // Will be populated as a slice of Tag pointers
   }
   ```

## How to Use References

To use struct references:

1. Define your data structures with appropriate reference fields
2. Add `structgen` tags to fields that should be populated from references
3. Pass all datasets to the `Generate` method

### Example

```go
// Define tag struct
type Tag struct {
    ID   string
    Name string
    Slug string
}

// Define post struct with references to tags
type Post struct {
    ID       string
    Title    string
    TagSlugs []string  // Contains tag slugs
    Tags     []*Tag    `structgen:"TagSlugs"` // Will be populated from TagSlugs
}

// Create your datasets
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

// For separate files approach:

// Create generator for tags
tagGenerator := genstruct.NewGenerator(
    genstruct.WithPackageName("main"),
    genstruct.WithOutputFile("tags.go"),
)

// Generate tags
err := tagGenerator.Generate(tags)
if err != nil {
    // handle error
}

// Create generator for posts with references to tags
postGenerator := genstruct.NewGenerator(
    genstruct.WithPackageName("main"),
    genstruct.WithOutputFile("posts.go"),
)

// Generate posts, passing tags as a reference dataset
err = postGenerator.Generate(posts, tags)
if err != nil {
    // handle error
}

// Alternatively, generate everything in one file:
allInOneGenerator := genstruct.NewGenerator(
    genstruct.WithPackageName("main"),
    genstruct.WithOutputFile("blog_data.go"),
)

// Generate both posts and tags in one go
err = allInOneGenerator.Generate(posts, tags)
```

## How genstruct Finds Matching References

When looking for matching references, `genstruct` tries each of the identifier fields in order:

1. For each value in the source field (e.g., each string in `TagSlugs`)
2. It looks through each struct in the reference dataset (e.g., each `Tag`)
3. It tries each identifier field (`ID`, `Name`, `Slug`, etc.) to find a match
4. When a match is found, that struct is added to the result

The identifier fields can be set using `WithIdentifierFields`, with a default of `["ID", "Name", "Slug", "Title", "Key", "Code"]`.

## Edge Cases

- **Missing References**: If a referenced value doesn't exist in the reference dataset, it will be omitted from the result.
- **Empty Source**: If the source field is empty, the target field will be an empty slice or struct.
- **Invalid Types**: If the source and target field types don't match the supported reference patterns, the field will be ignored.

## Limitations

Currently, there are some limitations to be aware of:

1. Only string-based lookups are supported (string to struct, string slice to struct slice)
2. References must be between simple types (no nested struct references)
3. References are resolved during code generation, not at runtime

## Future Extensions

Future versions may include:

- Support for custom lookup logic
- Bidirectional references
- Deeper nested references
- More complex mapping relationships