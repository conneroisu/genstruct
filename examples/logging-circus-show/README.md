# Logging Circus Show Example

This example demonstrates how to use genstruct with embedded struct references and logging in a circus show context.

## Features

- Uses Go's `log/slog` package for structured logging
- Demonstrates struct embedding and references with the `structgen` tag
- Shows how to use a custom logger with the generator
- Configurable verbosity level via command-line flags

## Domain Model

The example models a circus show with:

1. **Performers** - The circus artists who perform tricks
   - Each performer has a list of trick IDs they can perform

2. **Tricks** - The stunts and performances done by performers
   - Referenced by performers using the `structgen` tag

## Generated Code

The generated code will include:
- Constants for each performer and trick
- Variables for each performer and trick
- A slice of all performers
- A slice of all tricks
- References between performers and their tricks

## Usage

Run the example with the default settings:

```bash
go run main.go
```

Run with verbose logging:

```bash
go run main.go -v debug
```

## Code Structure

- `Trick` struct: Represents a circus trick with ID, name, description, difficulty, and required equipment
- `Performer` struct: Represents a circus performer with ID, name, role, experience, and trick references
- `createCustomLogger()`: Creates a custom logger with specific attributes for the circus show

## How Struct References Work

The `TrickIDs` field in the `Performer` struct is tagged with `structgen:"Trick"`. This tells genstruct:

1. The `TrickIDs` field will contain references to `Trick` structs
2. To look up the tricks by their IDs in the `Trick` slice
3. To link the performers to their tricks in the generated code

This creates a relationship between performers and their tricks, similar to a database foreign key.