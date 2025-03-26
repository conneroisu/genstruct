# Blog Posts and Tags Example with Many-to-Many Relationships

This example demonstrates how to use genstruct with many-to-many relationships between entities, specifically showing self-referential relationships in tags.

## Current Implementation

The current example shows two approaches to handling many-to-many relationships:

1. **Post-to-Tag relationship** (one-to-many): This works out of the box with genstruct, where `Post.Tags` references `Tag` objects via `Post.TagSlugs`.

2. **Tag-to-Tag relationship** (many-to-many): Currently demonstrated in `setupTagRelationships()` which sets up relationships at runtime to avoid initialization cycles.

## Enhancement Needed for Direct Many-to-Many Support

To directly support many-to-many self-references in the generated code, we need to enhance the genstruct generator to handle pointer references properly. This would allow code like:

```go
var TagGoProgramming = Tag{
    ID:       "tag-001",
    Name:     "Go Programming",
    Slug:     "go-programming",
    TagSlugs: []string{"code-generation", "developer-tools"},
    Tags:     []*Tag{&TagCodeGeneration, &TagDeveloperTools},
}
```

### Required Changes in the Generator

1. Add support for pointer references in `generateReferenceSlice()` (in values.go)
2. Update type detection to handle `[]*Struct` in addition to the current `[]Struct`
3. Generate pointer references with `&` operator for self-references

This enhancement would allow proper handling of circular references without initialization cycles by using pointers to refer to other structs.

## Current Workaround

Until the generator is enhanced, the current recommended approach for handling many-to-many relationships is to:

1. Keep the Tag slugs in your data structure
2. Set up relationships after initialization using the approach shown in `setupTagRelationships()`
3. Use a wrapper type like `TagWithRefs` to handle the circular references

## Testing

```
go test ./examples/blog-posts-tags
```

## Running the Example

```
go run ./examples/blog-posts-tags/main.go
```