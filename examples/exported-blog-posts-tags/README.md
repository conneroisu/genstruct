# Exported Blog Posts and Tags Example

This example demonstrates how to generate code for a blog with posts and tags, where:

1. The type definitions are in a separate package (`pkg`)
2. The generated code is exported to another package (`out`)
3. References between posts and tags are maintained using pointer slices

## Structure

- `pkg/pkg.go` - Contains type definitions and sample data
- `main.go` - Runs the code generation process
- `main_test.go` - Tests that verify the exported code generation works correctly
- `out/` - Directory where generated code will be written

## Key Features Demonstrated

1. **Cross-Package Type References** - How to reference types from another package
2. **Pointer-Based References** - Using `[]*Tag` for efficient memory usage
3. **External Package Generation** - Generating code in a separate package from the types
4. **Tag-Based Relationship Resolution** - Using the `structgen` tag to define relationships

## Running the Example

```bash
go run main.go
```

This will generate the code and display the contents of the generated file.

## Testing

```bash
go test -v
```

This will run the tests to verify that the code generation works correctly.