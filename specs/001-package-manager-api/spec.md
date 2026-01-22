# Feature Specification: Common Package Manager API

**Feature Branch**: `001-package-manager-api`
**Created**: 2026-01-22
**Status**: Draft
**Input**: User description: "Create an api that wraps various linux alternative package managers (brew, flatpak, snap) in a common api. Use common interfaces, and make sure that each package manager implementation can have an \"empty\" implementation of the methods in the interfaces. Create common definitions for the interface methods so there is no ambiguity about \"Update\" vs \"Upgrade\", where one might mean just refreshing metadata, and the other might mean making changes to installed packages. For long running operations, use a ProgressReporter interface that the client can provide to get status of the steps and messages about the progress. Prefer native go SDKs to talk to package managers, REST next if available, and finally fall back to exec/cli wrapping as a last resort."

## User Scenarios & Testing _(mandatory)_

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.

  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Use one API across managers (Priority: P1)

As a developer integrating package management into an application, I want a single,
consistent API that can target Brew, Flatpak, or Snap, so my application logic is not
coupled to any specific package manager.

**Why this priority**: This is the core value: a stable, common contract that enables
swapping package managers without rewriting business logic.

**Independent Test**: Can be fully tested with a stubbed/empty backend by verifying the
same client code can call the common API, receive consistent results/errors, and use
capability detection for unsupported operations.

**Acceptance Scenarios**:

1. **Given** a client using the common API, **When** it selects the Brew backend,
   **Then** the same set of common operations are available with consistent inputs and
   outputs.
2. **Given** a client using the common API, **When** it switches from Brew to Flatpak,
   **Then** the client code does not change beyond selecting a different backend.

---

### User Story 2 - Run long operations with progress (Priority: P2)

As a developer, I want long-running operations (install, uninstall, upgrade) to
report structured progress (Actions → Tasks → Steps) including informational, warning,
and error messages, so my application can show a meaningful status UI and logs.

**Why this priority**: Package operations can take minutes; users need confidence and
visibility into what is happening.

**Independent Test**: Can be tested by using a stub backend that emits a known progress
sequence and verifying the client receives the right hierarchy and message severities.

**Acceptance Scenarios**:

1. **Given** a long-running operation with a progress reporter provided, **When** the
   operation starts, **Then** it emits at least one Action with one or more Tasks and
   optional Steps.
2. **Given** the operation is running, **When** warnings occur, **Then** the progress
   stream includes Warning messages without necessarily failing the operation.

---

### User Story 3 - Predictable semantics and safe fallbacks (Priority: P3)

As a developer, I want consistent definitions for common operations (especially
"Update" vs "Upgrade") and predictable behavior when operations are unsupported, so
my application avoids ambiguous outcomes and can handle “not supported” cases cleanly.

**Why this priority**: Differences in naming across package managers cause confusion,
bugs, and unexpected system changes.

**Independent Test**: Can be tested by verifying that “Update” never changes installed
packages, and that unsupported operations in the empty implementation return a
consistent “not supported” result.

**Acceptance Scenarios**:

1. **Given** the common API, **When** a client calls Update, **Then** only metadata is
   refreshed and no installed packages are changed.
2. **Given** the common API, **When** a client calls Upgrade, **Then** installed
   packages may be changed according to the package manager’s upgrade behavior.
3. **Given** a backend uses an empty implementation for an operation, **When** the
   client calls that operation, **Then** it receives a consistent “not supported”
   result that can be programmatically detected.

---

## MVP Scope

The MVP for this feature includes:

- A stable common Go API and empty implementations for unsupported operations.
- At least one real backend integration (Brew) implementing one or more operations with deterministic unit tests.

For Brew, the preferred integration order is:

- REST (Formulae API) for metadata/search where applicable, then
- CLI with JSON output for local operations where REST isn’t applicable.

Non-MVP (deferred):

- Full read/write parity across all backends (Brew/Flatpak/Snap).
- Real integrations for Flatpak and Snap beyond empty implementations.

### Edge Cases

- Package manager not installed or not available on the system.
- Operation requires privileges (e.g., admin/root) and the caller does not have them.
- Network unavailable or repository metadata endpoints unavailable.
- Partial failures (some packages succeed, some fail) during batch operations.
- Cancellation mid-operation (client cancels; operation stops safely).
- Backend supports a concept only partially (e.g., can list but cannot upgrade).
- Progress reporting is not provided by the client (operation still works).

## Requirements _(mandatory)_

### Functional Requirements

- **FR-001**: System MUST provide a single common API surface that supports multiple
  package manager backends (initially: Brew, Flatpak, Snap).
- **FR-002**: Each backend MUST implement the common contract, enabling a client to
  switch backends without changing application logic beyond backend selection.
- **FR-003**: System MUST include an explicit “empty implementation” option for each
  backend that can satisfy the common contract while returning a consistent “not
  supported” result for operations that are intentionally unimplemented.
- **FR-004**: System MUST define unambiguous operation semantics, including:
  - **Update**: refresh metadata/indexes only; MUST NOT change installed packages.
  - **Upgrade**: may change installed packages (e.g., update versions, apply upgrades).
  - **Install/Uninstall**: changes installed packages.
- **FR-005**: System MUST provide a capability/introspection mechanism so callers can
  determine whether an operation is supported by a chosen backend before invoking it.
- **FR-006**: Long-running operations MUST support an optional progress reporting
  interface provided by the caller.
- **FR-007**: The progress reporting model MUST use a common definition:
  - Actions contain Tasks
  - Tasks contain zero or more Steps
  - Actions/Tasks/Steps may emit messages with severity Informational, Warning, or Error
- **FR-008**: Progress reporting MUST be consistent across backends and safe for
  concurrent use.
- **FR-009**: When integrating with a package manager, the system MUST prefer
  integration methods in this order:
  1. official/native SDK/API where available
  2. REST interface where available
  3. CLI execution/wrapping only as a last resort
- **FR-010**: When CLI execution is used, command invocation and output parsing MUST be
  abstracted to enable deterministic unit testing.

- **FR-011 (MVP)**: System MUST implement at least one real backend integration.
  For the MVP, Brew MUST implement at least one non-empty operation (e.g., Search and/or
  package metadata lookup) using REST (Formulae API) when applicable, and MUST have
  deterministic unit tests for parsing and error handling.

### Key Entities _(include if feature involves data)_

- **Package Manager Backend**: A selectable provider (Brew/Flatpak/Snap) implementing
  the shared contract and exposing supported capabilities.
- **Package Identifier**: A stable reference to a package, including name and optional
  source/namespace information.
- **Installed Package**: A package currently installed on the system with its version
  and status.
- **Operation**: A request to perform an action (Update, Upgrade, Install, Uninstall,
  Search, List, etc.) and its resulting outcome.
- **Capability**: A declaration that a backend supports a specific operation or feature.
- **Progress Action / Task / Step**: A structured progress hierarchy used to report
  long-running work.
- **Progress Message**: A message emitted during progress with severity Informational,
  Warning, or Error.

## Success Criteria _(mandatory)_

### Measurable Outcomes

- **SC-001**: A client can switch between Brew, Flatpak, and Snap backends without
  changing their application logic beyond selecting a different backend.
- **SC-002**: The common API defines Update vs Upgrade semantics such that Update does
  not modify installed packages, and this is verifiable via acceptance scenarios.
- **SC-003**: Long-running operations can emit progress with Actions/Tasks/Steps and
  severity-tagged messages, and a client can render a meaningful progress view.
- **SC-004**: Unsupported operations are programmatically detectable and do not result
  in ambiguous behavior (e.g., silent no-ops without a clear result).

## Assumptions

- Target environment is Linux where Brew/Flatpak/Snap may be present.
- Not all package managers will be installed; the system must handle “not available”.
- Some backends may lack SDK/API/REST options; CLI may be required as fallback.

## Out of Scope (initial release)

- Support for additional package managers beyond Brew/Flatpak/Snap.
- A UI application; this feature is an API intended to be used by other programs.

## API Surface Stability

### Version 0.x (Current - Experimental)

**Status**: Under active development, breaking changes may occur

**Stable Elements** (unlikely to change):

- Core interfaces: `Manager`, `Updater`, `Upgrader`, `Installer`, `Uninstaller`, `Searcher`, `Lister`
- Backend constructors: `NewBrew()`, `NewFlatpak()`, `NewSnap()`
- Error types: `NotSupportedError`, `NotAvailableError`, `ExternalFailureError`
- Error detection: `IsNotSupported()`, `IsNotAvailable()`, `IsExternalFailure()`
- Progress hierarchy: `ProgressAction` → `ProgressTask` → `ProgressStep` → `ProgressMessage`
- Severity levels: `SeverityInfo`, `SeverityWarning`, `SeverityError`

**May Change**:

- Options struct fields (may add new optional fields)
- Result struct fields (may add new informational fields)
- Backend-specific behavior details
- Progress message content and formatting
- Constructor option patterns (currently no options, may add functional options)

**Update vs Upgrade Contract** (guaranteed stable):

- `Update`: MUST NOT modify installed packages (metadata refresh only)
- `Upgrade`: MAY modify installed packages (install new versions)
- `UpdateResult`: Will never have `PackagesChanged` field (enforced by struct definition)
- `UpgradeResult`: Will always have `PackagesChanged` field

### Version 1.0 (Future - Stable)

**Stability Guarantees** (planned for 1.0):

- No breaking changes to public interfaces
- Backward-compatible additions only
- Semantic versioning commitment
- Deprecation warnings before removal
- Migration guides for any necessary changes

**What qualifies for 1.0**:

- All three backends (Brew, Flatpak, Snap) have full implementation
- Comprehensive test coverage (>80%)
- Production usage validation
- Documentation completeness
- Performance benchmarks established

**Pre-1.0 Migration Strategy**:

- Import path will remain `github.com/frostyard/pm`
- Major version changes (if any) will use Go modules versioning (`/v2`, `/v3`, etc.)
- Deprecation notices will be added at least one minor version before removal

### Current Implementation Status

| Backend | Available | Capabilities | Search | Update | Upgrade | Install | Uninstall | List |
| ------- | --------- | ------------ | ------ | ------ | ------- | ------- | --------- | ---- |
| Brew    | ✅        | ✅           | ✅     | ⚠️     | ⚠️      | ⚠️      | ⚠️        | ⚠️   |
| Flatpak | ✅        | ✅           | ⚠️     | ⚠️     | ⚠️      | ⚠️      | ⚠️        | ⚠️   |
| Snap    | ✅        | ✅           | ⚠️     | ⚠️     | ⚠️      | ⚠️      | ⚠️        | ⚠️   |

Legend:

- ✅ Fully implemented
- ⚠️ Returns `NotSupported` (empty implementation)
- ❌ Not implemented

### Breaking Change Policy (0.x)

During 0.x development:

- **Minor version bumps** (0.1 → 0.2): May include breaking changes
- **Patch version bumps** (0.1.0 → 0.1.1): Bug fixes only, no breaking changes
- **Pre-release tags** (`-alpha`, `-beta`, `-rc`): Experimental, use at own risk

Users depending on 0.x versions should:

- Pin to specific versions in `go.mod`
- Review changelogs before upgrading
- Test thoroughly after version changes
- Expect API evolution based on real-world usage feedback

## Backend Integration Notes

### Flatpak

- Preferred integration order for `pm`: SDK/API → REST → CLI (FR-009).
- Flatpak has an official C SDK (**libflatpak**), but using it from Go requires cgo + GLib/GObject bindings; for the initial version, `pm` should default to wrapping the `flatpak` CLI.
- CLI supports `--columns=...` for list-style queries and `-y/--assumeyes` + `--noninteractive` for automation, but does not advertise JSON output in standard manpages.
- Details and links: see [specs/001-package-manager-api/research.md](specs/001-package-manager-api/research.md).
