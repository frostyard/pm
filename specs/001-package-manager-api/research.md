# Research: Common Package Manager API

**Feature**: `001-package-manager-api`
**Date**: 2026-01-22

## Integration Methods (SDK/API → REST → CLI)

### Homebrew (brew)

- Decision: Use `brew` CLI JSON output as the primary machine-readable interface (e.g., `brew info --json`, `brew outdated --json`). Optionally use the public Formulae API for metadata without requiring a local Homebrew install.
- Rationale:
  - Homebrew does not provide a stable external SDK for third-party consumers; the CLI is canonical.
  - JSON output significantly reduces brittleness vs parsing human-formatted text.
  - Formulae API provides a supported REST-style metadata surface.
- Alternatives considered:
  - Calling internal Ruby APIs: rejected (not stable as an external contract).
  - Parsing plain text CLI output: rejected unless no JSON output is available.

References:

- https://docs.brew.sh/Manpage
- https://docs.brew.sh/Querying-Brew
- https://formulae.brew.sh/docs/api/

### Snap (snapd)

- Decision: Prefer snapd’s local REST API (`/v2/...`, HTTP over unix socket) for list/search/install/remove/refresh and for long-running progress via “changes”.
- Rationale:
  - snapd exposes an official, structured API used by the `snap` CLI.
  - Long-running operations map naturally to a progress model (polling changes and surfacing tasks/messages).
  - Deterministic unit tests are straightforward with injected HTTP transport and fixtures.
- Alternatives considered:
  - Using `github.com/snapcore/snapd/client`: rejected by default due to GPL-3.0 licensing risk.
  - Wrapping `snap` CLI via `exec`: kept as fallback.

Notes:

- Many mutating actions require elevated privileges; unit tests MUST NOT depend on system state.

### Flatpak

- Decision: Default to wrapping the Flatpak CLI using documented “machine-readable-ish” controls (`--columns`, `--show-details`) and noninteractive flags.
- Rationale:
  - libflatpak exists (official C/GObject API) but requires cgo + system dev dependencies; it is a complexity/cross-compile burden for an initial release.
  - Standard manpages document stable field selection via `--columns`, but do not advertise JSON output.
  - Deterministic unit tests are achievable with an injected command runner and fixture parsing.
- Alternatives considered:
  - libflatpak via cgo: deferred (complexity); possible future build-tagged implementation.
  - D-Bus helper interfaces: deferred (not intended as a stable third-party integration contract).

Key references:

- libflatpak API reference: https://docs.flatpak.org/en/latest/libflatpak-api-reference.html
- `flatpak list` columns: https://manpages.debian.org/bookworm/flatpak/flatpak-list.1.en.html
- `flatpak remote-ls` columns: https://manpages.debian.org/bookworm/flatpak/flatpak-remote-ls.1.en.html
- `flatpak remotes` columns: https://manpages.debian.org/bookworm/flatpak/flatpak-remotes.1.en.html
- `flatpak search` columns: https://manpages.debian.org/bookworm/flatpak/flatpak-search.1.en.html
- `flatpak install` noninteractive: https://manpages.debian.org/bookworm/flatpak/flatpak-install.1.en.html
- `flatpak update` noninteractive: https://manpages.debian.org/bookworm/flatpak/flatpak-update.1.en.html
- `flatpak uninstall` noninteractive: https://manpages.debian.org/bookworm/flatpak/flatpak-uninstall.1.en.html

## Progress Reporting Model

- Decision: Standardize progress as Actions → Tasks → Steps, each optionally emitting Messages with severity Informational | Warning | Error.
- Rationale:
  - Meets the constitution requirement for a backend-agnostic progress contract.
  - Aligns with snapd’s “changes” model and can be layered over CLI backends.
- Alternatives considered:
  - Only streaming logs: rejected (hard to build consistent UI/UX).
  - Only percent-complete: rejected (insufficient for multi-phase tasks).

## Update vs Upgrade Semantics

- Decision: Define `Update` as metadata refresh only (no changes to installed packages). Define `Upgrade` as potentially changing installed packages.
- Rationale:
  - Removes ambiguity across package managers.
  - Enables safe refresh flows in client applications.

## Testing & Determinism

- Decision: Make all system interactions injectable.
- Rationale:
  - Unit tests MUST NOT require package managers installed or network access.
  - Parsing logic can be tested with fixed fixtures.

Concrete approach:

- CLI backends (brew/flatpak fallback): inject a `Runner` interface and parse fixture stdout.
- Snap backend: inject HTTP client/transport and use `httptest` for deterministic tests.
