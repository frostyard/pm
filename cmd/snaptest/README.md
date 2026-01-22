# snaptest - Snap Backend Test Harness

A simple CLI tool to test and demonstrate the `github.com/frostyard/pm` library using the Snap backend.

## Building

```bash
make build-snaptest
# or
go build -o bin/snaptest ./cmd/snaptest
```

## Usage

```bash
./bin/snaptest <command> [args]
```

### Commands

#### Search for packages

```bash
./bin/snaptest search <query>
```

Example: `./bin/snaptest search firefox`

#### List installed packages

```bash
./bin/snaptest list
```

#### Install packages

```bash
./bin/snaptest install <package>...
```

Example: `./bin/snaptest install firefox vlc`

#### Uninstall packages

```bash
./bin/snaptest uninstall <package>...
```

Example: `./bin/snaptest uninstall firefox`

#### Update package metadata

```bash
./bin/snaptest update
```

#### Upgrade installed packages

```bash
./bin/snaptest upgrade
```

#### Show backend capabilities

```bash
./bin/snaptest capabilities
```

## Features

- **Progress Reporting**: Shows real-time progress of operations
- **Error Handling**: Demonstrates proper error detection (NotAvailable, NotSupported, ExternalFailure)
- **Interface Usage**: Shows how to use different backend interfaces (Searcher, Installer, etc.)

## Example Output

```
$ ./bin/snaptest capabilities
Backend Capabilities:
  ✓ Search (via snap find CLI)
  ✓ UpdateMetadata (via snap refresh CLI)
  ✓ UpgradePackages (via snap refresh CLI)
  ✓ Install (via snap install CLI)
  ✓ Uninstall (via snap remove CLI)
  ✓ ListInstalled (via snap list CLI)

$ ./bin/snaptest list
→ ListInstalled
  • Running snap list
Installed packages (28):
  firefox 133.0.3
  vlc 3.0.21
  ...
```

## Notes

- This tool requires Snap to be installed on your Linux system
- All operations use the real Snap backend - exercise caution with install/uninstall/upgrade commands
- Progress output shows the operation hierarchy: Actions → Tasks → Steps with messages
- Snap is typically pre-installed on Ubuntu and many other Linux distributions
