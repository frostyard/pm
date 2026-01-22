# Contracts: Integration Preference Order

**Feature**: `001-package-manager-api`

When implementing a backend, the integration methods MUST be attempted in this order:

1. Official/native SDK or API (local library, supported DBus API, etc.)
2. REST API (local daemon API, official HTTP API)
3. CLI wrapping via `exec` (last resort)

Additional rules:

- If a higher-preference method exists but is impractical due to licensing (e.g., GPL),
  the backend MAY use the next method, but MUST document the decision in the plan.
- If CLI wrapping is used, output parsing MUST be deterministic and covered by unit tests.
