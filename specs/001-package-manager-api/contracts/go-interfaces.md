# Contracts: Go Interfaces

**Feature**: `001-package-manager-api`

## Naming and Semantics

- **Update**: Refresh metadata/indexes only. MUST NOT change installed packages.
- **Upgrade**: May change installed packages (upgrade versions, apply changes).

## Common Interfaces (conceptual)

> Exact names may change during implementation, but semantics MUST match.

### Manager

- `Capabilities(ctx) ([]Capability, error)`
- `Available(ctx) (bool, error)`

### Updater

- `Update(ctx, opts UpdateOptions) (UpdateResult, error)`

### Upgrader

- `Upgrade(ctx, opts UpgradeOptions) (UpgradeResult, error)`

### Installer

- `Install(ctx, pkgs []PackageRef, opts InstallOptions) (InstallResult, error)`

### Uninstaller

- `Uninstall(ctx, pkgs []PackageRef, opts UninstallOptions) (UninstallResult, error)`

### Searcher

- `Search(ctx, query string, opts SearchOptions) ([]PackageRef, error)`

### Lister

- `ListInstalled(ctx, opts ListOptions) ([]InstalledPackage, error)`

## Empty Implementations

- Each backend MUST provide an “empty” implementation option that satisfies the common
  interfaces but returns a consistent “not supported” result for unimplemented methods.
- The “not supported” result MUST be programmatically detectable.
