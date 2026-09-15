# Project Scope

## Product outcome

MysticLight Control Center will provide one local control plane for monitoring
and coordinating RGB lighting, cooling telemetry, LCD status, profiles, and
game-focus events without taking ownership away from the existing providers.

## Intended users

Desktop users who run multiple hardware-control services and need a single,
observable, recoverable interface with safe behavior when a provider is absent
or unhealthy.

## In scope

1. A provider-neutral domain model for capabilities, devices, zones, profiles,
   health, and diagnostics.
2. Read-only discovery and telemetry adapters for OpenRGB, cooling, LCD, and
   sensor providers.
3. Explicit ownership rules for RGB, PWM, LCD, and telemetry operations.
4. A local event pipeline for provider state, game focus, profile changes, and
   recovery events.
5. Profile validation, safe fallback selection, health reporting, and restart
   recovery.
6. A desktop UI after the domain and service contracts are stable.
7. Packaging, service integration, upgrade, rollback, and compatibility
   documentation.

## Out of scope

- New hardware drivers or firmware changes.
- Experimental raw HID, SMBus, or I2C operations.
- Automatic provider replacement or simultaneous ownership of one resource.
- Cloud accounts, remote control, telemetry upload, or mandatory internet use.
- Storing credentials, serial numbers, stable hardware identifiers, or raw
  host identity data.

## MVP boundary

The first usable milestone is a local service that can discover provider
availability, expose normalized read-only health and telemetry, validate one
profile format, and report ownership conflicts without performing physical
hardware writes. RGB, PWM, and LCD writes remain disabled until their adapter
contracts and rollback behavior have been separately reviewed.

## Success criteria

- Provider failure is visible and does not crash the control plane.
- No resource has more than one declared write owner.
- Invalid profiles are rejected before reaching a provider.
- Restart recovery returns to a documented safe state.
- Automated tests use mocks and perform no physical hardware writes.
- Diagnostics contain no credentials, personal identifiers, addresses, or
  stable device identifiers.

## Delivery phases

| Phase | Outcome | Exit criteria |
| --- | --- | --- |
| 0. Contracts | Architecture, ownership, threat, and language decisions | Reviewed RFCs and test strategy |
| 1. Core model | Configuration/profile schema and validation | Deterministic schema tests pass |
| 2. Read-only providers | Provider capability and telemetry adapters | Mock contract tests pass |
| 3. Service/event layer | Health, diagnostics, events, and recovery | Failure-injection tests pass |
| 4. UI | Status, profiles, and diagnostics screens | UI consumes service contracts only |
| 5. Operations | Packaging, service integration, rollback | Reproducible install and rollback checks |
| 6. Hardware enablement | Reviewed, limited provider writes | Explicit authorization and physical test plan |

## Review gates

Each phase requires a focused change, validation evidence, a privacy review,
and an updated rollback note. Hardware-affecting work requires an additional
review and must remain isolated from read-only development until approved.
