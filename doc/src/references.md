# Struct Reference Embedding

One of the powerful features of `genstruct` is the ability to automatically populate fields in one struct with references to other structs. This allows you to define relationships between your data structures and have them automatically resolved during code generation.

## How It Works

References between structs are established using Go struct tags with the `structgen` tag:

```go
type Post struct {
    ID       string    
    TagSlugs []string       // Contains slugs of tags
    Tags     []Tag     `structgen:"TagSlugs"` // Will be populated based on TagSlugs
}
```

The `structgen` tag tells `genstruct` that:

1. The value of this field (`Tags`) should be populated based on the value of another field (`TagSlugs`)
2. It should look up matching values in the reference data provided

## Supported Reference Types

Currently, `genstruct` supports two types of references:

1. **String to Struct**: A string field referencing a single struct
   ```go
   type Post struct {
       AuthorID string            // Contains an author ID
       Author   Author     `structgen:"AuthorID"` // Will be populated from AuthorID
   }
   ```

2. **String Slice to Struct Slice**: A slice of strings referencing a slice of structs
   ```go
   type Post struct {
       TagSlugs []string       // Contains slugs of tags
       Tags     []Tag     `structgen:"TagSlugs"` // Will be populated as a slice of Tags
   }
   ```

## How to Use References

To use struct references:

1. Define your data structures with appropriate reference fields
2. Add `structgen` tags to fields that should be populated from references
3. Pass all required datasets to `NewGenerator`

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
    Tags     []Tag     `structgen:"TagSlugs"` // Will be populated from TagSlugs
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

// Generate code with both datasets
generator := genstruct.NewGenerator(config, posts, tags)
err := generator.Generate()
```

## How genstruct Finds Matching References

When looking for matching references, `genstruct` tries each of the identifier fields in order:

1. For each value in the source field (e.g., each string in `TagSlugs`)
2. It looks through each struct in the reference dataset (e.g., each `Tag`)
3. It tries each identifier field (`ID`, `Name`, `Slug`, etc.) to find a match
4. When a match is found, that struct is added to the result

The identifier fields are specified in `Config.IdentifierFields`, with a default of `["ID", "Name", "Slug", "Title", "Key", "Code"]`.

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