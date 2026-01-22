# flatpaktest - Flatpak Backend Test Harness

A simple CLI tool to test and demonstrate the `github.com/frostyard/pm` library using the Flatpak backend.

## Building

```bash
make build-flatpaktest
# or
go build -o bin/flatpaktest ./cmd/flatpaktest
```

## Usage

```bash
./bin/flatpaktest <command> [args]
```

### Commands

#### Search for packages

```bash
./bin/flatpaktest search <query>
```

Example: `./bin/flatpaktest search firefox`

#### List installed packages

```bash
./bin/flatpaktest list
```

#### Install packages

```bash
./bin/flatpaktest install <package>...
```

Example: `./bin/flatpaktest install org.mozilla.firefox`

#### Uninstall packages

```bash
./bin/flatpaktest uninstall <package>...
```

Example: `./bin/flatpaktest uninstall org.mozilla.firefox`

#### Update package metadata

```bash
./bin/flatpaktest update
```

#### Upgrade installed packages

```bash
./bin/flatpaktest upgrade
```

#### Show backend capabilities

```bash
./bin/flatpaktest capabilities
```

## Features

- **Progress Reporting**: Shows real-time progress of operations
- **Error Handling**: Demonstrates proper error detection (NotAvailable, NotSupported, ExternalFailure)
- **Interface Usage**: Shows how to use different backend interfaces (Searcher, Installer, etc.)

## Example Output

```text
$ ./bin/flatpaktest capabilities
Backend Capabilities:
  ✓ Search (via flatpak search CLI)
  ✓ UpdateMetadata (via flatpak update CLI)
  ✓ UpgradePackages (via flatpak update CLI)
  ✓ Install (via flatpak install CLI)
  ✓ Uninstall (via flatpak uninstall CLI)
  ✓ ListInstalled (via flatpak list CLI)

$ ./bin/flatpaktest list
→ ListInstalled
  • Running flatpak list
Installed packages (12):
  org.mozilla.firefox 133.0.3
  org.gimp.GIMP 2.10.38
  ...
```

## Notes

- This tool requires Flatpak to be installed on your Linux system
- All operations use the real Flatpak backend - exercise caution with install/uninstall/upgrade commands
- Progress output shows the operation hierarchy: Actions → Tasks → Steps with messages
- Package names should be full application IDs (e.g., `org.mozilla.firefox`)
