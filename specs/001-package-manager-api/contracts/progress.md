# Contracts: Progress Reporting

**Feature**: `001-package-manager-api`

## Model

- Actions have Tasks
- Tasks have 0..N Steps
- Actions/Tasks/Steps may emit Messages
- Messages have severity: Informational | Warning | Error

## ProgressReporter Contract

- Provided by calling code.
- MUST be safe for concurrent calls.
- Backends MUST emit progress in a backend-agnostic way.

### Required behaviors

- An operation SHOULD emit at least one Action.
- An Action SHOULD emit at least one Task.
- Tasks MAY omit Steps (0 steps is valid).
- Errors MAY be reported as Progress Messages and/or returned as operation errors.
- Warning messages MUST NOT automatically fail the operation.

## Long-running Operation Integration

- Operations that can take noticeable time (install/uninstall/upgrade/update) MUST accept
  an optional ProgressReporter via options.
- If ProgressReporter is nil/not provided, operations MUST still work.
