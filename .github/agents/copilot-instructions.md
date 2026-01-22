# pm Development Guidelines

Auto-generated from all feature plans. Last updated: 2026-01-22

## Active Technologies

- Go (see `go.mod` / toolchain) + Go standard library (core); tooling via Makefile (`golangci-lint`) (001-package-manager-api)

## Project Structure

```text
src/
tests/
```

## Commands

# Add commands for Go (see `go.mod` / toolchain)

## Code Style

Go (see `go.mod` / toolchain): Follow standard conventions

## Recent Changes

- 001-package-manager-api: Added Go (see `go.mod` / toolchain) + Go standard library (core); tooling via Makefile (`golangci-lint`)

<!-- MANUAL ADDITIONS START -->

## Canonical Workflow (pm)

- Use the Makefile as the primary interface: `make tools`, `make fmt`, `make lint`, `make build`, `make test`, `make check`.
- After any code changes, run `make check`.

## Repo Structure (Go library)

```text
.
├── Makefile
├── go.mod
├── *.go              # public package: github.com/frostyard/pm
└── internal/
	├── runner/       # command runner abstraction + fakes
	└── backend/
		├── brew/
		├── flatpak/
		└── snap/
```

<!-- MANUAL ADDITIONS END -->
