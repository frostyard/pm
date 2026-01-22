<!--
Sync Impact Report
==================
Version change: 1.2.0 → 1.3.0

Modified principles:
- I. Go API Best Practices
- II. Common Interfaces per Package Manager
- III. Unit Tests Required for Every Change (NON-NEGOTIABLE)
- IV. Deterministic, Testable Package Manager Adapters
- V. Quality Gates (Format/Lint/Build/Test)

Added sections:
- API & Compatibility
- Error Handling & Observability
- Development Workflow

Amendments in this version:
- Prefer interacting with package managers via API/SDK when available; use `exec` only
	as a fallback.
- Makefile is the primary interface for formatting, linting, building, testing, and
	installing required tools.
- Add a required ProgressReporter interface model for long-running operations.

Removed sections:
- None (template was placeholders)

Templates requiring updates:
- ✅ .specify/templates/tasks-template.md
- ✅ .specify/templates/plan-template.md

Follow-up TODOs:
- None
-->

# pm Constitution

## Core Principles

### I. Go API Best Practices

All exported APIs MUST follow established Go conventions and be easy to use correctly.

- Public APIs MUST accept a `context.Context` as the first parameter when they may block,
  perform I/O, execute external commands, or otherwise be cancelable.
- Exported identifiers MUST be documented with GoDoc comments.
- Interfaces MUST be small and behavior-focused (prefer “accept interfaces, return
  concrete types” where appropriate).
- Avoid global state; configuration MUST be explicit (e.g., constructor options).
- Public types and functions MUST have stable semantics; breaking changes require a
  major version bump.

### II. Common Interfaces per Package Manager

Each supported package manager (initially: Homebrew, Flatpak, Snap) MUST implement a
shared set of Go interfaces so callers can switch backends without changing their code.

- The library MUST define a small set of common interfaces (e.g., `Manager`, `Searcher`,
  `Installer`, `Upgrader`, `Uninstaller`, `Lister`, etc.).
- Each backend MUST implement the common interfaces and adhere to the shared contracts
  (inputs, outputs, idempotency expectations, and error semantics).
- The library MUST provide a consistent way to construct/select managers (factory,
  registry, or explicit constructors) without callers needing backend-specific wiring.

Progress reporting for long-running operations:

- The library MUST define a `ProgressReporter` interface that can be used by calling
  code to observe progress for long-running operations.
- The progress model MUST use a common definition:
  - Actions have Tasks
  - Tasks have zero or more Steps
  - Actions/Tasks/Steps MAY emit Messages
  - Messages MUST have a severity: Informational, Warning, or Error
- Progress reporting MUST be backend-agnostic and consistent across package managers.
- Implementations MUST be safe to call from concurrent goroutines.
- Long-running operations SHOULD accept a `ProgressReporter` via explicit options or
  dependency injection (do not require global variables).

### III. Unit Tests Required for Every Change (NON-NEGOTIABLE)

Every code change MUST be accompanied by unit tests that meaningfully cover the change.

- New logic MUST have direct unit test coverage.
- Bug fixes MUST include a regression test.
- Tests MUST be deterministic and MUST NOT depend on the developer machine’s state
  (installed packages, network access, current OS configuration).
- External command execution MUST be abstracted behind interfaces and mocked in unit
  tests.

### IV. Deterministic, Testable Package Manager Adapters

Package manager integrations MUST be designed for correctness, portability, and
testability.

- Adapters MUST isolate system interactions (command execution, filesystem reads,
  environment variables) behind injectable dependencies.
- Adapters MUST prefer package-manager APIs/SDKs when available and stable; resort to
  command execution (`exec`) only when no viable API/SDK exists.
- Parsing of command output MUST be robust and covered by unit tests.
- Operations MUST be explicit about side effects and MUST prefer “dry-run” style
  capabilities where feasible.
- Concurrency MUST be safe-by-default; shared mutable state is forbidden unless
  protected and justified.

### V. Quality Gates (Format/Lint/Build/Test)

Any time an agent modifies code, it MUST run the project’s quality gates.

- The Makefile is the canonical entrypoint for these gates; contributors and agents
  MUST use `make` targets instead of ad-hoc commands.
- Code MUST be formatted via `make fmt`.
- Lint MUST be run via `make lint`.
- Build/compile MUST be checked via `make build`.
- Tests MUST be run via `make test`.
- The primary local gate is `make check`.

## API & Compatibility

- This project is a Go library providing an API to interact with package managers.
- The initial supported backends are: Homebrew (`brew`), Flatpak (`flatpak`), Snap (`snap`).
- The API MUST present a consistent abstraction layer; backend-specific details SHOULD be
  exposed only via optional extension interfaces or feature flags.
- Public API surface area MUST be kept small and intentional; avoid premature generality.

## Error Handling & Observability

- Errors MUST be actionable and wrap underlying causes using `fmt.Errorf("...: %w", err)`.
- Command failures MUST include useful context (command, exit status when available,
  and stderr excerpt where safe).
- Do not log by default from the library; prefer returning rich errors. If logging is
  provided, it MUST be opt-in and use an injected logger interface.

## Development Workflow

- Every feature/change MUST be specified and planned before implementation.
- Plans MUST include an explicit “Constitution Check” and list the `make` targets used
  for formatting, linting, building, and testing.
- PR review (or equivalent) MUST verify:
  - Interface contracts are upheld across all backends
  - Unit tests cover changes
  - Quality gates were run successfully

## Governance

- This constitution governs all development in this repository.
- Amendments require documenting rationale and impact.
- Versioning follows semantic versioning:
  - MAJOR: Breaking API changes or principle removal/redefinition
  - MINOR: New principles/sections or materially expanded requirements
  - PATCH: Clarifications and non-semantic refinements

**Version**: 1.3.0 | **Ratified**: 2026-01-22 | **Last Amended**: 2026-01-22
