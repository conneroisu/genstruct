# Development Environment

Genstruct provides a comprehensive development environment using Nix Flakes, making it easy to get started with consistent tooling and dependencies. This document covers how to set up and use the development environment.

## Getting Started with Nix

Genstruct uses Nix Flakes to provide a reproducible development environment with all the necessary tools pre-configured.

### Prerequisites

1. Install Nix following the [official installation instructions](https://nixos.org/download.html)
2. Enable Flakes with the following settings in your `~/.config/nix/nix.conf`:

```
experimental-features = nix-command flakes
```

### Entering the Development Environment

To enter the development environment, run:

```bash
nix develop
```

This will drop you into a shell with all the required tools and commands available.

## Available Development Tools

The Nix development environment provides a rich set of tools for Go development:

### Go Tools

- **Go 1.24**: The Go compiler and standard tools
- **air**: Live reload for Go applications
- **pprof**: Profiling tool for Go applications
- **revive**: Fast, configurable, extensible, flexible, and beautiful linter for Go
- **golangci-lint**: Fast Go linters runner
- **gopls**: The Go language server
- **templ**: HTML template language for Go
- **golines**: A Go formatter that automatically shortens long lines
- **gomarkdoc**: Generates Markdown documentation for Go code
- **gotests**: Tool to generate table-driven tests
- **reftools**: Tools for updating Go references

### Formatters and Linters

- **prettierd**: Daemon for the Prettier formatter
- **alejandra**: Nix code formatter
- **nixd**: Nix language server

## Predefined Scripts

The development environment includes several useful scripts for common tasks:

| Command | Description |
|---------|-------------|
| `dx` | Edit flake.nix |
| `tests` | Run short Go tests |
| `unit-tests` | Run all Go tests |
| `coverage-tests` | Run all Go tests with coverage |
| `lint` | Run golangci-lint |
| `generate-all` | Generate all code artifacts and format them |
| `format` | Format code files across multiple languages |

### Using the Scripts

These scripts can be run directly from the shell:

```bash
# Run the tests
tests

# Run the linter
lint

# Format all code
format
```

## Documentation Development

The documentation is built using [mdBook](https://rust-lang.github.io/mdBook/). To build and preview the documentation:

```bash
# Navigate to the doc directory
cd doc

# Build the documentation
mdbook build

# Serve the documentation locally (with live reload)
mdbook serve --open
```

## Continuous Integration

The repository is set up with GitHub Actions for continuous integration. The workflow runs:

1. Tests across multiple Go versions
2. Linting with golangci-lint
3. Documentation build checks

## Tips for Efficient Development

1. **Use Air for Live Reload**: When developing, use Air to automatically rebuild and run your code on changes:
   ```bash
   air -c .air.toml
   ```

2. **Generate Test Coverage Report**: Generate and view test coverage:
   ```bash
   coverage-tests
   go tool cover -html=coverage.out
   ```

3. **Format Code Before Commit**: Always format your code before committing:
   ```bash
   format
   ```

4. **Use gopls with Your Editor**: Configure your editor to use gopls for Go code intelligence.

5. **Follow Go Best Practices**: Use the provided linters to ensure your code follows Go best practices.

## Example Development Workflow

Here's a typical workflow for developing with genstruct:

1. Clone the repository and enter the development environment:
   ```bash
   git clone https://github.com/conneroisu/genstruct.git
   cd genstruct
   nix develop
   ```

2. Make your changes to the code

3. Run the tests to verify your changes:
   ```bash
   tests
   ```

4. Format your code:
   ```bash
   format
   ```

5. Run the linter to check for issues:
   ```bash
   lint
   ```

6. Generate any necessary code artifacts:
   ```bash
   generate-all
   ```

7. Commit your changes and create a pull request