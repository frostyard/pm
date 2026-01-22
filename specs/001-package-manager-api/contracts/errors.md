# Contracts: Errors

**Feature**: `001-package-manager-api`

## Not Supported

- Unsupported operations MUST return a programmatically-detectable error.
- The error MUST be stable across backends and must not rely on string matching.

## Not Available

- If a backend is not installed/reachable, operations MUST return a programmatically-detectable error.

## External Failures

- CLI backends MUST surface stderr (sanitized as needed) to aid debugging.
- REST backends MUST surface structured error payloads where available.
