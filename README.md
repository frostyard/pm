# pm - Package Manager Abstraction Library

[![Go Reference](https://pkg.go.dev/badge/github.com/frostyard/pm.svg)](https://pkg.go.dev/github.com/frostyard/pm)
[![Test](https://github.com/frostyard/pm/actions/workflows/test.yml/badge.svg)](https://github.com/frostyard/pm/actions/workflows/test.yml)
[![Lint](https://github.com/frostyard/pm/actions/workflows/lint.yml/badge.svg)](https://github.com/frostyard/pm/actions/workflows/lint.yml)

A Go library that provides a unified interface for interacting with multiple package managers (Homebrew, Flatpak, and Snap). The library abstracts package management operations like search, install, uninstall, update, and upgrade across different backends with consistent error handling and progress reporting.

## Features

- **Multi-Backend Support**: Unified API for Homebrew, Flatpak, and Snap
- **Consistent Interface**: Same methods work across all supported package managers
- **Progress Reporting**: Built-in progress reporting with hierarchical action/task/step tracking
- **Error Handling**: Structured error types with detailed context
- **CLI & API Operations**: Brew uses REST API for search, all backends use CLI for operations
- **Type Safety**: Strongly typed package references and operation results

## Supported Package Managers

| Backend      | Search | Update | Upgrade | Install | Uninstall | List   |
| ------------ | ------ | ------ | ------- | ------- | --------- | ------ |
| **Homebrew** | ✅ API | ✅ CLI | ✅ CLI  | ✅ CLI  | ✅ CLI    | ✅ CLI |
| **Flatpak**  | ✅ CLI | ✅ CLI | ✅ CLI  | ✅ CLI  | ✅ CLI    | ✅ CLI |
| **Snap**     | ✅ CLI | ✅ CLI | ✅ CLI  | ✅ CLI  | ✅ CLI    | ✅ CLI |

## Installation

```bash
go get github.com/frostyard/pm
```

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/frostyard/pm"
)

func main() {
    // Create a Homebrew backend
    mgr := pm.NewBrew()

    ctx := context.Background()

    // Check if the backend is available
    available, err := mgr.Available(ctx)
    if err != nil {
        log.Fatalf("Failed to check availability: %v", err)
    }
    if !available {
        log.Fatal("Homebrew is not available")
    }

    // Search for packages
    packages, err := mgr.Search(ctx, "wget", pm.SearchOptions{})
    if err != nil {
        log.Fatalf("Search failed: %v", err)
    }

    fmt.Printf("Found %d packages:\n", len(packages))
    for _, pkg := range packages {
        fmt.Printf("  - %s\n", pkg.Name)
    }

    // Install a package
    result, err := mgr.Install(ctx, []pm.PackageRef{
        {Name: "wget", Kind: "formula"},
    }, pm.InstallOptions{})
    if err != nil {
        log.Fatalf("Install failed: %v", err)
    }

    if result.Changed {
        fmt.Println("Package installed successfully")
    }
}
```

### With Progress Reporting

```go
package main

import (
    "context"
    "fmt"

    "github.com/frostyard/pm"
)

// Simple progress reporter that prints to stdout
type StdoutReporter struct{}

func (r *StdoutReporter) BeginAction(name string) {
    fmt.Printf("==> %s\n", name)
}

func (r *StdoutReporter) EndAction() {
    fmt.Println()
}

func (r *StdoutReporter) BeginTask(name string) {
    fmt.Printf("  → %s\n", name)
}

func (r *StdoutReporter) EndTask() {}

func (r *StdoutReporter) BeginStep(name string) {
    fmt.Printf("    • %s\n", name)
}

func (r *StdoutReporter) EndStep() {}

func (r *StdoutReporter) Info(message string) {
    fmt.Printf("    ℹ %s\n", message)
}

func (r *StdoutReporter) Error(message string) {
    fmt.Printf("    ✗ %s\n", message)
}

func main() {
    // Create backend with progress reporter
    mgr := pm.NewFlatpak(pm.WithProgress(&StdoutReporter{}))

    ctx := context.Background()

    // Update metadata with progress reporting
    result, err := mgr.Update(ctx, pm.UpdateOptions{
        Progress: &StdoutReporter{},
    })
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }

    if result.Changed {
        fmt.Println("Metadata updated")
    } else {
        fmt.Println("Metadata already up to date")
    }
}
```

### Backend Capabilities

Check what operations a backend supports:

```go
capabilities, err := mgr.Capabilities(ctx)
if err != nil {
    log.Fatalf("Failed to get capabilities: %v", err)
}

for _, cap := range capabilities {
    fmt.Printf("%s: %v (%s)\n", cap.Operation, cap.Supported, cap.Notes)
}
```

## API Overview

### Core Interfaces

- `Manager`: Main interface combining all package management operations
- `Searcher`: Search for packages
- `Updater`: Update package metadata/indices
- `Upgrader`: Upgrade installed packages
- `Installer`: Install packages
- `Uninstaller`: Remove packages
- `Lister`: List installed packages

### Creating Backends

```go
// Homebrew
brew := pm.NewBrew(opts...)

// Flatpak
flatpak := pm.NewFlatpak(opts...)

// Snap
snap := pm.NewSnap(opts...)
```

### Constructor Options

```go
// Add progress reporting
mgr := pm.NewBrew(pm.WithProgress(reporter))
```

### Error Handling

The library provides structured error types:

```go
result, err := mgr.Install(ctx, packages, opts)
if err != nil {
    switch {
    case pm.IsNotSupported(err):
        fmt.Println("Operation not supported")
    case pm.IsNotAvailable(err):
        fmt.Println("Backend not available")
    case pm.IsExternalFailure(err):
        // Get detailed error information
        extErr := err.(*pm.ExternalFailureError)
        fmt.Printf("Command failed:\nStdout: %s\nStderr: %s\n",
            extErr.Stdout, extErr.Stderr)
    default:
        fmt.Printf("Error: %v\n", err)
    }
}
```

## Test Harnesses

The repository includes three CLI test harnesses demonstrating library usage:

- **brewtest**: Homebrew operations demo
- **flatpaktest**: Flatpak operations demo
- **snaptest**: Snap operations demo

Build them with:

```bash
make build-cli
```

Run them:

```bash
./bin/brewtest search wget
./bin/flatpaktest list
./bin/snaptest capabilities
```

See individual README files in `cmd/*/README.md` for detailed usage.

## Development

### Prerequisites

- Go 1.22 or later
- Make
- One or more of: Homebrew, Flatpak, or Snap (for testing)

### Building

```bash
# Install development tools
make tools

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Run all checks (format, lint, test)
make check
```

### Running Tests

```bash
# Unit tests
make test

# Full CI suite
make ci
```

## Architecture

The library is organized as:

- **`pm` package**: Public API with `Manager` interface and error types
- **`internal/types`**: Shared internal types for operations and results
- **`internal/backend/*`**: Backend implementations (brew, flatpak, snap)
- **`internal/runner`**: Command execution wrapper with structured error handling
- **`cmd/*`**: Example CLI tools demonstrating library usage

### Backend Design

Each backend implements the same set of interfaces:

1. **Available check**: Verify the package manager is installed and accessible
2. **Capabilities**: Report which operations are supported
3. **Operations**: Search, update, upgrade, install, uninstall, list

### Progress Reporting

Progress reporting follows a three-level hierarchy:

1. **Action**: Top-level operation (e.g., "Install")
2. **Task**: Major steps within an action (e.g., "Downloading packages")
3. **Step**: Fine-grained progress updates (e.g., "Fetching package metadata")

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Acknowledgments

Built for the [Frostyard](https://github.com/frostyard) project to provide unified package management across different Linux package managers and Homebrew.
