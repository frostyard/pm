# Quickstart: Common Package Manager API

**Feature**: `001-package-manager-api`

## Prereqs

- Go toolchain installed

## Developer workflow

- Install tools: `make tools`
- Run full local gate: `make check`

## Example usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/frostyard/pm"
)

// MyProgressReporter implements pm.ProgressReporter interface
type MyProgressReporter struct{}

func (p *MyProgressReporter) OnAction(action pm.ProgressAction) {
    if !action.EndedAt.IsZero() {
        fmt.Printf("[Action] %s completed\n", action.Name)
    }
}

func (p *MyProgressReporter) OnTask(task pm.ProgressTask) {}

func (p *MyProgressReporter) OnStep(step pm.ProgressStep) {}

func (p *MyProgressReporter) OnMessage(msg pm.ProgressMessage) {
    fmt.Printf("[%s] %s\n", msg.Severity, msg.Text)
}

func main() {
    ctx := context.Background()

    // Create a backend instance
    backend := pm.NewBrew()

    // Check if backend is available
    available, err := backend.Available(ctx)
    if err != nil {
        log.Fatalf("Error checking availability: %v", err)
    }
    if !available {
        log.Fatal("Homebrew not available")
    }

    // Check capabilities
    caps, err := backend.Capabilities(ctx)
    if err != nil {
        log.Fatalf("Error getting capabilities: %v", err)
    }

    fmt.Printf("Backend supports %d operations\n", len(caps))

    // Type-assert to specific interfaces as needed
    updater, hasUpdate := backend.(pm.Updater)
    upgrader, hasUpgrade := backend.(pm.Upgrader)
    searcher, hasSearch := backend.(pm.Searcher)

    // Update: refresh metadata only (never modifies packages)
    if hasUpdate {
        result, err := updater.Update(ctx, pm.UpdateOptions{
            Progress: &MyProgressReporter{},
        })
        if pm.IsNotSupported(err) {
            fmt.Println("Update not supported")
        } else if err != nil {
            log.Printf("Update failed: %v", err)
        } else {
            fmt.Printf("Metadata updated: %v\n", result.Changed)
        }
    }

    // Upgrade: may change installed packages
    if hasUpgrade {
        result, err := upgrader.Upgrade(ctx, pm.UpgradeOptions{
            Progress: &MyProgressReporter{},
        })
        if pm.IsNotSupported(err) {
            fmt.Println("Upgrade not supported")
        } else if err != nil {
            log.Printf("Upgrade failed: %v", err)
        } else {
            fmt.Printf("Upgraded %d packages\n", len(result.PackagesChanged))
        }
    }

    // Search for packages
    if hasSearch {
        packages, err := searcher.Search(ctx, "nodejs", pm.SearchOptions{})
        if pm.IsNotSupported(err) {
            fmt.Println("Search not supported")
        } else if err != nil {
            log.Printf("Search failed: %v", err)
        } else {
            fmt.Printf("Found %d packages\n", len(packages))
        }
    }
}
```

### Available Backends

- **Brew**: `pm.NewBrew()` - Homebrew backend (macOS/Linux)
- **Flatpak**: `pm.NewFlatpak()` - Flatpak backend (Linux)
- **Snap**: `pm.NewSnap()` - Snap backend (Linux)

### Error Handling

The API provides structured error types:

```go
_, err := backend.(pm.Installer).Install(ctx, packages, pm.InstallOptions{})

if pm.IsNotSupported(err) {
    // Operation not supported by this backend
} else if pm.IsNotAvailable(err) {
    // Backend not available (not installed or unreachable)
} else if pm.IsExternalFailure(err) {
    // External command or API failed
    // Can extract stdout/stderr for debugging
} else if err != nil {
    // Other error
}
```

### Progress Reporting

All operations accept an optional `Progress` field for real-time updates:

```go
type MyReporter struct{}

func (r *MyReporter) OnAction(action pm.ProgressAction) {
    fmt.Printf("Action: %s\n", action.Name)
}

func (r *MyReporter) OnTask(task pm.ProgressTask) {
    fmt.Printf("  Task: %s\n", task.Name)
}

func (r *MyReporter) OnStep(step pm.ProgressStep) {
    fmt.Printf("    Step: %s\n", step.Name)
}

func (r *MyReporter) OnMessage(msg pm.ProgressMessage) {
    switch msg.Severity {
    case pm.SeverityInfo:
        fmt.Printf("    ℹ %s\n", msg.Text)
    case pm.SeverityWarning:
        fmt.Printf("    ⚠ %s\n", msg.Text)
    case pm.SeverityError:
        fmt.Printf("    ✗ %s\n", msg.Text)
    }
}

// Use with any operation
backend.(pm.Searcher).Search(ctx, "git", pm.SearchOptions{
    Progress: &MyReporter{},
})
```

### Update vs Upgrade Semantics

**Update** (refreshes metadata only):

- **MUST NOT** modify installed packages
- Safe to run repeatedly
- Examples: `apt update`, `brew update`

**Upgrade** (modifies packages):

- **MAY** install newer versions
- Changes system state
- Examples: `apt upgrade`, `brew upgrade`

Notes:

- Progress reporting is optional; operations work without it (nil-safe)
- All operations respect `context.Context` for cancellation and timeouts
- Backend constructors take no required arguments (variadic options for future extensions)
