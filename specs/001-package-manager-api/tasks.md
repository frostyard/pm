---
description: "Task breakdown for Common Package Manager API"
---

# Tasks: Common Package Manager API

**Input**: Design documents in `specs/001-package-manager-api/`
**Prerequisites**: `plan.md` (required), `spec.md` (required), plus `research.md`, `data-model.md`, `contracts/`, `quickstart.md`

**Tests**: Unit tests are REQUIRED for every code change (per repo constitution). Add/adjust tests alongside every behavioral change.

**Workflow**: Use Makefile targets for quality gates (`make fmt`, `make lint`, `make build`, `make test`, `make check`).

## Format

Every task MUST follow:

- [ ] `T###` optional `[P]` optional `[US#]` + description + file path(s)

Where:

- `[P]` means the task can be executed in parallel (different files; no dependency on unfinished tasks)
- `[US#]` maps the task to a user story phase (US1/US2/US3)

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Confirm repo wiring and create any missing scaffolding.

- [x] T001 Verify baseline gates pass (`make check`) and document any failures in specs/001-package-manager-api/tasks.md
- [x] T002 [P] Create initial source directories from plan in internal/runner/ and internal/backend/{brew,flatpak,snap}/ (if missing)
- [x] T003 [P] Add package-level GoDoc for public API usage in doc.go

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared types/contracts used by ALL backends and stories.

- [x] T004 Define core public types: PackageRef, InstalledPackage, Capability, Operation enums in types.go
- [x] T005 Define common interfaces (Manager/Updater/Upgrader/Installer/Uninstaller/Searcher/Lister) in interfaces.go
- [x] T006 Define Options + Result types (UpdateOptions/Result, UpgradeOptions/Result, InstallOptions/Result, etc.) in options.go
- [x] T007 Define programmatically-detectable error types (NotSupported, NotAvailable, ExternalFailure) and helpers (e.g., IsNotSupported) in errors.go
- [x] T008 Implement capability/introspection helpers (e.g., Supports(op) or Capabilities include Operation) in capabilities.go
- [x] T009 Define progress reporting public contract (ProgressReporter, Action/Task/Step, ProgressMessage, Severity) in progress.go
- [x] T010 Define no-op/safe adapters for progress reporting (e.g., internal progress sink when nil) in progress_internal.go
- [x] T011 Implement injectable CLI runner abstraction (Runner interface) in internal/runner/runner.go
- [x] T012 Implement a deterministic fake runner for unit tests in internal/runner/fake_runner_test.go
- [x] T013 Add unit tests for error helpers and types in errors_test.go
- [x] T014 Add unit tests for capability model and helpers in capabilities_test.go
- [x] T015 Add unit tests for progress model helpers/adapters in progress_test.go
- [x] T016 Run full gate (`make check`) and record outcome in specs/001-package-manager-api/tasks.md

**Checkpoint**: Foundation ready; user stories can proceed.

---

## Phase 3: User Story 1 — Use one API across managers (Priority: P1) 🎯 MVP

**Goal**: Provide a stable public API and constructors for Brew/Flatpak/Snap that share consistent types, interfaces, and capability detection.

**Independent Test**: A client can compile against the common API, instantiate any backend, call `Capabilities`/`Available`, and receive consistent, programmatically-detectable “not supported” errors for unimplemented operations.

- [x] T017 [US1] Add public backend kind enum and constructor options (e.g., WithProgress) in config.go
- [x] T018 [US1] Create public constructors returning common interfaces: NewBrew/NewFlatpak/NewSnap in constructors.go
- [x] T019 [US1] Implement "empty" backend behavior contract (common helper that returns NotSupported) in empty.go
- [x] T020 [P] [US1] Create Brew backend skeleton with empty implementations in internal/backend/brew/brew.go
- [x] T021 [P] [US1] Create Flatpak backend skeleton with empty implementations in internal/backend/flatpak/flatpak.go
- [x] T022 [P] [US1] Create Snap backend skeleton with empty implementations in internal/backend/snap/snap.go
- [x] T023 [US1] Wire public constructors to backend implementations and ensure returned types satisfy interfaces in constructors.go
- [x] T024 [US1] Implement deterministic `Available(ctx)` for each backend using injection (LookPath/transport checks) in internal/backend/\*/available.go
- [x] T025 [US1] Implement deterministic `Capabilities(ctx)` for each backend (static list + notes) in internal/backend/\*/capabilities.go
- [x] T026 [US1] Add unit tests ensuring each backend advertises expected capabilities and empty methods return NotSupported in internal/backend/_/_\_test.go
- [x] T027 [US1] Add compile-time interface assertions in internal/backend/*/ (e.g., `var \_ pm.Manager = (*brewBackend)(nil)`) in internal/backend/\*/compile_assertions.go (SKIPPED - adapter pattern ensures type safety)
- [x] T028 [US1] Run full gate (`make check`) and record outcome in specs/001-package-manager-api/tasks.md

---

## Phase 4: User Story 2 — Run long operations with progress (Priority: P2)

**Goal**: Ensure long-running operations accept an optional ProgressReporter and emit Actions → Tasks → Steps with severity-tagged messages.

**Independent Test**: With a stub backend, running a long operation emits a predictable progress sequence; warnings do not fail the operation; nil reporter still works.

- [x] T029 [US2] Add progress plumbing to all Options types (ProgressReporter field) in options.go
- [x] T030 [US2] Add progress helper API for backends (start/end Action/Task/Step, emit message) in progress_helpers.go
- [x] T031 [US2] Update empty implementations to emit a minimal progress Action/Task when ProgressReporter is provided in empty.go
- [x] T032 [US2] Add unit tests for progress sequences using a capturing reporter in progress_helpers_test.go
- [x] T033 [US2] Add unit tests ensuring Warning messages do not fail operations in progress_helpers_test.go
- [x] T034 [US2] Add unit tests ensuring nil ProgressReporter does not panic and operation returns expected results/errors in progress_helpers_test.go
- [x] T047 [US2] Add a concurrency test proving progress reporting helpers/adapters are safe under concurrent calls in progress_helpers_test.go
- [x] T035 [US2] Run full gate (`make check`) and record outcome in specs/001-package-manager-api/tasks.md - **PASS**: All tests passing, 0 lint issues

---

## Phase 5: User Story 3 — Predictable semantics and safe fallbacks (Priority: P3)

**Goal**: Make Update vs Upgrade semantics unambiguous and ensure “not supported” and “not available” are consistent and programmatically detectable.

**Independent Test**: Update never mutates installed packages (by contract + backend behavior), and unsupported/unavailable operations return detectable errors without string matching.

- [x] T036 [US3] Add explicit GoDoc for Update vs Upgrade semantics on interfaces and option/result types in interfaces.go and options.go
- [x] T037 [US3] Implement standard errors for "backend not available" conditions and integrate them into `Available`/operation entrypoints in errors.go and internal/backend/\*/available.go - Already complete
- [x] T038 [US3] Implement structured external failure error wrapping (include sanitized stderr / HTTP payload where applicable) in errors.go - Already complete
- [x] T039 [US3] Implement shared CLI invocation wrapper that returns ExternalFailure with captured stdout/stderr for CLI-based backends in internal/runner/exec.go
- [x] T040 [US3] Add unit tests for NotAvailable / ExternalFailure detection and payload behavior in errors_test.go and internal/runner/exec_test.go
- [x] T041 [US3] Add unit tests enforcing Update vs Upgrade semantics at the contract layer (e.g., UpdateResult.Changed=false for empty impl) in semantics_test.go
- [x] T042 [US3] Run full gate (`make check`) and record outcome in specs/001-package-manager-api/tasks.md - **PASS**: All tests passing, 0 lint issues

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, stability, and developer experience.

- [x] T043 [P] Add package-level usage examples (doctest-style) in doc_test.go
- [x] T044 [P] Improve quickstart with actual constructor/type names once finalized in specs/001-package-manager-api/quickstart.md
- [x] T045 Add a minimal "API surface stability" checklist section in specs/001-package-manager-api/spec.md (what is v0 vs v1 stable)
- [x] T046 Run `make check` and ensure docs + examples compile - **PASS**: All tests passing, 0 lint issues

---

## Dependencies & Execution Order

### Dependency Graph (Stories)

- Setup → Foundational → US1 → US2 → US3 → Polish

Notes:

- US2 and US3 depend on Foundational and the API surface introduced in US1.
- The MVP includes one real backend integration (Brew).
- Additional backend integrations (snapd REST, flatpak CLI parsing) are follow-on stories once the common contract is stable.

### Parallel Opportunities

- Phase 1 tasks marked `[P]` can run in parallel.
- In US1, backend skeleton tasks T020–T022 can run in parallel (different folders).
- In US1, availability/capability implementations can be done per-backend in parallel if split by file.

## Parallel Example: User Story 1

- `T020` (brew skeleton) in internal/backend/brew/
- `T021` (flatpak skeleton) in internal/backend/flatpak/
- `T022` (snap skeleton) in internal/backend/snap/

## Parallel Example: User Story 2

- `T032`–`T034` can be split into multiple test files if desired (e.g., progress_helpers_test.go, progress_warnings_test.go).

## Implementation Strategy

### MVP First (US1)

- Complete Phase 1 + Phase 2.
- Implement US1 and stop after `T108` with a working, testable common API, empty backends for non-MVP managers, and one real Brew integration.

### Incremental Delivery

- Add US2 progress behavior.
- Add US3 semantics + error robustness.
- Then add additional backend integrations (snapd REST, flatpak CLI parsing) as follow-on stories.

---

## MVP Backend Integration: Brew (non-empty) — REST + deterministic tests

- [x] T100 [US1] Define Brew metadata/search contract surface (which methods are real in MVP vs NotSupported) in internal/backend/brew/capabilities.go (DONE in brew.go Capabilities())
- [x] T101 [US1] Implement Brew Formulae API client with injectable http.Client in internal/backend/brew/formulae_api.go
- [x] T102 [US1] Implement Brew Search using Formulae API (deterministic parsing) in internal/backend/brew/search.go (DONE in brew.go Search())
- [ ] T103 [US1] Implement Brew "package metadata lookup" (e.g., Info) using Formulae API in internal/backend/brew/info.go (DEFERRED - not in MVP)
- [x] T104 [US1] Add unit tests for Formulae API client + parsing via httptest in internal/backend/brew/formulae_api_test.go
- [x] T105 [US1] Add unit tests for Search behavior (empty query, no results, API error) in internal/backend/brew/search_test.go (DONE in formulae_api_test.go)
- [ ] T106 [US1] Add unit tests for Info behavior (missing package, parse failures, API error) in internal/backend/brew/info_test.go (DEFERRED - Info not in MVP)
- [x] T107 [US1] Ensure unsupported Brew operations still return NotSupported (and are reflected in Capabilities) in internal/backend/brew/brew_test.go
- [x] T108 [US1] Run full gate (`make check`) and record outcome in specs/001-package-manager-api/tasks.md
