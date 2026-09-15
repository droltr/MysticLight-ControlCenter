# Architecture RFC 0001: Safe Local Coordination Plane

- Status: Proposed
- Date: 2026-09-15
- Scope: Initial architecture and ownership contract

## Context

The project coordinates several existing providers with different APIs,
failure modes, and ownership boundaries. Directly coupling the UI to those
providers would make failures difficult to isolate and could create competing
writes to the same physical resource.

## Decision

Use a local-first layered architecture:

```text
Desktop UI
    |
Local application service / API
    |
Domain model + profile validator + ownership registry
    |
Event bus + health/recovery coordinator
    |
Provider adapters (OpenRGB, cooling, LCD, sensors, focus)
    |
Existing provider services and OS interfaces
```

The domain layer owns normalized concepts and policy. Adapters own protocol
translation and provider-specific error handling. The UI consumes normalized
service contracts and never calls a provider directly.

## Ownership contract

| Resource | Owner | Initial operation policy |
| --- | --- | --- |
| Fan and pump PWM | CoolerControl/CoolerDash | Read-only integration |
| RGB lighting | OpenRGB | Read-only integration |
| LCD display | LCD provider | Read-only integration |
| Sensor telemetry | Existing sensor provider | Read-only integration |
| Game-focus state | Focus adapter | Read-only integration |
| Profile and safety policy | Control Center | Local validation and coordination |

The control center may request a provider operation only after the provider
contract, ownership, safety limits, and rollback path are explicitly reviewed.

## Core contracts

Every adapter should expose:

- `capabilities`: supported operations and read/write classification;
- `health`: availability, last successful observation, and redacted error;
- `observe()`: normalized telemetry or state;
- `apply()`: disabled by default and separately gated for future writes;
- `close()`: deterministic resource cleanup.

The service layer must treat provider responses as untrusted input, apply
timeouts, isolate failures, and preserve the last known safe state without
inventing telemetry.

## Event model

Events are local, typed, timestamped, and privacy-safe. Initial event types
are `ProviderHealthChanged`, `TelemetryObserved`, `ProfileValidated`,
`ProfileActivated`, `FocusChanged`, and `RecoveryRequested`. Events must not
contain credentials, serials, network addresses, or raw provider payloads.

## Failure and recovery policy

1. Start with all physical writes disabled.
2. Discover providers and capabilities.
3. Mark unavailable or conflicting providers as degraded.
4. Continue read-only telemetry where safe.
5. On restart, restore observation mode before considering any requested write.
6. Record a redacted diagnostic and provide a deterministic rollback action.

## Technology decision

No implementation language is selected by this RFC. A separate decision record
must compare suitable local-desktop candidates using provider interoperability,
concurrency, packaging, testing, dependency provenance, and long-term
maintenance. Implementation should not begin until that comparison is
reviewed.

## Alternatives considered

- UI-first integration: rejected because it couples presentation to provider
  failures and makes ownership enforcement inconsistent.
- Direct hardware access: rejected because it violates the provider ownership
  model and increases hardware risk.
- Cloud-first orchestration: rejected because local control and privacy are
  primary requirements.

## Open questions

- Which provider APIs and versions are supported in the first compatibility
  matrix?
- Which service transport best fits the selected implementation language?
- What profile fields are stable across all supported providers?
- Which diagnostics are useful without exposing unique device identity?
