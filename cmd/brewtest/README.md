# brewtest - Brew Backend Test Harness

A simple CLI tool to test and demonstrate the `github.com/frostyard/pm` library using the Brew backend.

## Building

```bash
make build-brewtest
# or
go build -o bin/brewtest ./cmd/brewtest
```

## Usage

```bash
./bin/brewtest <command> [args]
```

### Commands

#### Search for packages

```bash
./bin/brewtest search <query>
```

Example: `./bin/brewtest search nodejs`

#### List installed packages

```bash
./bin/brewtest list
```

#### Install packages

```bash
./bin/brewtest install <package>...
```

Example: `./bin/brewtest install wget curl`

#### Uninstall packages

```bash
./bin/brewtest uninstall <package>...
```

Example: `./bin/brewtest uninstall wget`

#### Update package metadata

```bash
./bin/brewtest update
```

#### Upgrade installed packages

```bash
./bin/brewtest upgrade
```

#### Show backend capabilities

```bash
./bin/brewtest capabilities
```

## Features

- **Progress Reporting**: Shows real-time progress of operations
- **Error Handling**: Demonstrates proper error detection (NotAvailable, NotSupported, ExternalFailure)
- **Interface Usage**: Shows how to use different backend interfaces (Searcher, Installer, etc.)

## Example Output

```text
$ ./bin/brewtest capabilities
Backend Capabilities:
  ✓ Search (via Formulae API)
  ✓ UpdateMetadata (via brew update CLI)
  ✓ UpgradePackages (via brew upgrade CLI)
  ✓ Install (via brew install CLI)
  ✓ Uninstall (via brew uninstall CLI)
  ✓ ListInstalled (via brew list CLI)

$ ./bin/brewtest list
→ ListInstalled
  • Running brew list
Installed packages (142):
  autoconf 2.72
  automake 1.17
  ...
```

## Notes

- This tool requires Homebrew to be installed on macOS or Linux
- All operations use the real Brew backend - exercise caution with install/uninstall/upgrade commands
- Progress output shows the operation hierarchy: Actions → Tasks → Steps with messages
