# Pointer Slices Example

This example demonstrates how to use pointers and slices of pointers with genstruct.

## Features Demonstrated

1. Single pointer references (`*Author`)
2. Slices of pointer references (`[]*Comment`)
3. Bi-directional references between different types

## Data Model

This example uses three related types:

- **Author**: Represents a content creator
- **Comment**: User feedback with a pointer to its Author
- **Article**: A blog article with a pointer to its Author and a slice of pointers to Comments

## Relationships

- Each Article has one Author (single pointer: `*Author`)
- Each Article has multiple Comments (slice of pointers: `[]*Comment`)
- Each Comment has one Author (single pointer: `*Author`)

## Running the Example

```bash
go run main.go
```

This will generate `articles_generated.go` with constants, variables, and slices for all the data, including proper pointer references.

## Generated Code Usage

The generated code enables you to:

1. Access specific articles via their constants: `ArticleUnderstandingPointerSlices`
2. Navigate directly to their author: `ArticleUnderstandingPointerSlices.Author.Name`
3. Access all comments: `ArticleUnderstandingPointerSlices.Comments`
4. Navigate from comments back to their authors: `ArticleUnderstandingPointerSlices.Comments[0].Author.Name`

## Testing

Run the tests with:

```bash
go test -v
```