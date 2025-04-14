# Exported Blog Network Example

This example demonstrates how to generate exported code for a network of blog-related data types that use embedded structs.

## Data Structure

The example consists of four interconnected data types:
- `Post`: Blog posts with references to tags and projects
- `Tag`: Tags for categorizing posts and projects
- `Project`: Projects that group multiple posts
- `Employment`: Employment history with associated tags

Each of these types embeds a common `Embedded` struct that provides all the base fields and relationships.

## Running the Example

To run this example:

```bash
go run main.go
```

This will generate a file in the `out` directory containing statically defined data that can be imported and used in your application.

## Key Features

1. **Embedded Struct Pattern**: Demonstrates the use of the Go embedding feature to create a common base for all entities
2. **Cross-References**: Shows how to set up bidirectional relationships between different entity types
3. **Export Mode**: Generates code in a separate package that can be imported and used elsewhere
4. **Custom Identifiers**: Uses the slug field for naming variables instead of IDs

## Relationship Resolution

The structgen tag is used to automatically resolve relationships between entities, such as:
- `Posts` field is populated from `PostSlugs` references
- `Tags` field is populated from `TagSlugs` references
- `Projects` field is populated from `ProjectSlugs` references
- `Employments` field is populated from `EmploymentSlugs` references

This creates a fully connected graph of entities that can be traversed in all directions.