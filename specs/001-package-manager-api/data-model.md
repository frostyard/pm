# Data Model: Common Package Manager API

**Feature**: `001-package-manager-api`
**Date**: 2026-01-22

> This describes conceptual entities and validation rules. It is not a Go type dump.

## Entities

### Manager (Backend)

Represents a concrete package manager implementation.

- Fields
  - `Kind`: enum {`brew`, `flatpak`, `snap`}
  - `DisplayName`: string
  - `Available`: bool (installed/reachable)
  - `Capabilities`: set of `Capability`

### Capability

Declares what the backend can do.

- Fields
  - `Operation`: enum (UpdateMetadata, UpgradePackages, Install, Uninstall, Search, ListInstalled, ListAvailable, etc.)
  - `Notes`: optional string (why unsupported / constraints)

### PackageRef

Identifies a package in a backend.

- Fields
  - `Name`: string (required)
  - `Namespace`: optional string (e.g., flatpak remote/app id scope)
  - `Channel`: optional string (snap)
  - `Kind`: optional enum (app/runtime; cask/formula)

Validation:

- `Name` must be non-empty.

### InstalledPackage

Represents an installed package.

- Fields
  - `Ref`: PackageRef
  - `Version`: string
  - `InstalledAt`: optional time
  - `Status`: enum (installed, held, disabled, etc.)

### OperationRequest

A request to perform an operation.

- Fields
  - `Op`: enum (Update, Upgrade, Install, Uninstall, Search, ListInstalled)
  - `Packages`: []PackageRef (for Install/Uninstall/Upgrade subsets)
  - `Options`: key/value options (backend-agnostic; backend-specific via extensions)
  - `Progress`: optional ProgressReporter

Validation:

- `Op` required.
- For Install/Uninstall, `Packages` must be non-empty.

### OperationResult

Outcome of an operation.

- Fields
  - `Changed`: bool
  - `PackagesChanged`: []PackageRef
  - `Messages`: []ProgressMessage (summary)

### ProgressAction / ProgressTask / ProgressStep

Hierarchy for long-running work.

- Action
  - `ID`: string
  - `Name`: string
  - `StartedAt`/`EndedAt`: time (optional)
- Task
  - `ID`: string
  - `ActionID`: string
  - `Name`: string
  - `StartedAt`/`EndedAt`: time (optional)
- Step
  - `ID`: string
  - `TaskID`: string
  - `Name`: string
  - `StartedAt`/`EndedAt`: time (optional)

### ProgressMessage

- Fields
  - `Severity`: enum {Informational, Warning, Error}
  - `Text`: string
  - `Timestamp`: time
  - `ActionID`/`TaskID`/`StepID`: optional association
