# Implementation Plan: Common Package Manager API

**Branch**: `001-package-manager-api` | **Date**: 2026-01-22 | **Spec**: specs/001-package-manager-api/spec.md
**Input**: Feature specification from specs/001-package-manager-api/spec.md

**Note**: This template is filled in by the `/speckit.plan` command. See `.github/prompts/speckit.plan.prompt.md` and `.github/agents/speckit.plan.agent.md` for the execution workflow.

## Summary

Provide a single Go library API (`github.com/frostyard/pm`) that abstracts Brew, Flatpak,
and Snap behind common interfaces with unambiguous Update vs Upgrade semantics, optional
empty implementations, and structured progress reporting for long-running operations.

Integration preference order (per spec + constitution): SDK/API → REST → CLI/exec (last resort).

Design decisions and backend research live in:

- [specs/001-package-manager-api/research.md](specs/001-package-manager-api/research.md)

## Technical Context

**Language/Version**: Go (see `go.mod` / toolchain)
**Primary Dependencies**: Go standard library (core); tooling via Makefile (`golangci-lint`)
**Storage**: N/A
**Testing**: `make test` (unit tests required for every change; deterministic)
**Target Platform**: Linux hosts where Brew/Flatpak/Snap may be present
**Project Type**: Go library
**Performance Goals**: Library overhead negligible vs package-manager I/O; parsing linear in output size
**Constraints**: Unit tests must not require installed package managers; prefer SDK/API → REST → CLI; avoid GPL dependencies in core library
**Scale/Scope**: Single-host operations; scope limited to Brew/Flatpak/Snap backends and common operations

## Constitution Check

_GATE: Must pass before Phase 0 research. Re-check after Phase 1 design._

[Gates determined based on constitution file]

Pass/Fail (pre-Phase 0): PASS

- I. Go API Best Practices: planned (context-first APIs, GoDoc, small interfaces).
- II. Common Interfaces per Package Manager: planned (shared interfaces + empty implementations + ProgressReporter).
- III. Unit Tests Required: enforced via tasks and `make check`.
- IV. Prefer SDK/API → REST → CLI: planned per backend decisions.
- V. Quality Gates via Makefile: enforced via `make check`.

Minimum required gates for this repo:

- `make tools`
- `make fmt`
- `make lint`
- `make build`
- `make test`
- `make check`

## Project Structure

### Documentation (this feature)

```text
specs/001-package-manager-api/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
Makefile
go.mod
doc.go

internal/
  runner/                 # injected command runner interface + fakes
  backend/
    brew/                 # brew implementation (JSON CLI + optional API)
    flatpak/              # flatpak implementation (CLI)
    snap/                 # snapd REST implementation

# Public API lives at repo root package: github.com/frostyard/pm
# Tests live alongside code as *_test.go
```

**Structure Decision**: Go library with internal backends. Public constructors return
common interfaces so callers are not coupled to backend packages.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

None.

## Phased Execution

### Phase 0: Research

- Completed: [specs/001-package-manager-api/research.md](specs/001-package-manager-api/research.md)

### Phase 1: Design

- Common interfaces and capability model: [specs/001-package-manager-api/contracts/go-interfaces.md](specs/001-package-manager-api/contracts/go-interfaces.md)
- Progress model contract: [specs/001-package-manager-api/contracts/progress.md](specs/001-package-manager-api/contracts/progress.md)
- Integration preference contract: [specs/001-package-manager-api/contracts/integration-preference.md](specs/001-package-manager-api/contracts/integration-preference.md)
- Error contract: [specs/001-package-manager-api/contracts/errors.md](specs/001-package-manager-api/contracts/errors.md)
- Entity model: [specs/001-package-manager-api/data-model.md](specs/001-package-manager-api/data-model.md)
- Usage overview: [specs/001-package-manager-api/quickstart.md](specs/001-package-manager-api/quickstart.md)

### Phase 2: Task Breakdown (handled by /speckit.tasks)

- Create tasks grouped by user story (US1/US2/US3), each independently testable.
- Include explicit tasks to run `make check` for each story.
