# Logging Example

This example demonstrates how to use the built-in logging functionality in genstruct.

## Features

- Uses Go's `log/slog` package for structured logging
- Configurable verbosity level via command-line flags
- Supports different output formats (text or JSON)
- Configurable output destination (stderr, stdout, or file)

## Usage

Run the example with the default settings:

```bash
go run main.go
```

Run with verbose logging:

```bash
go run main.go -v debug
```

Run with JSON output format:

```bash
go run main.go -log-format json
```

Change the output destination:

```bash
go run main.go -log-output stdout
```

Or log to a file:

```bash
go run main.go -log-output generation.log
```

## Available Flags

- `-v`: Log verbosity level (debug, info, warn, error). Default: info
- `-log-format`: Log format (text, json). Default: text
- `-log-output`: Log output destination (stderr, stdout, or a file path). Default: stderr