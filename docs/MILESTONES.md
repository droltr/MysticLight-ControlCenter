# Project Milestones

Milestones are ordered by safety, contract stability, and operational value.
No milestone may bypass the ownership, privacy, or rollback gates.

## M1 — Contracts and governance

**Status:** Complete

- Define purpose, users, scope, non-goals, and safety boundaries.
- Establish provider ownership and privacy rules.
- Select Go 1.27.1 for the service and adapter layers.

## M2 — Core domain and provider contracts

**Status:** In progress

- Model resources, capabilities, ownership, provider health, and profiles.
- Reject duplicate ownership and unknown providers deterministically.
- Provide a registry and mock adapter for isolated tests.

**Exit criteria:** domain and provider unit tests pass without external
services or hardware.

## M3 — Service and event layer

**Status:** Next

- Add provider polling with context timeouts.
- Normalize observations and health transitions.
- Emit typed local events with redacted payloads.
- Add failure-injection and recovery tests.

**Exit criteria:** provider failures are isolated, observable, and recoverable
in deterministic mock tests.

## M4 — Read-only provider adapters

**Status:** Planned

- Implement OpenRGB, cooling, LCD, sensor, and focus adapters.
- Keep each adapter behind the provider interface.
- Document supported provider versions and capability gaps.

**Exit criteria:** contract tests pass with recorded compatibility evidence;
physical writes remain disabled.

## M5 — Profile runtime and safe fallback

**Status:** Planned

- Load and validate profiles.
- Select safe fallback behavior when capabilities are missing.
- Restore observation mode after restart.

**Exit criteria:** invalid, conflicting, unavailable, and recovery scenarios
are covered by automated tests.

## M6 — Desktop UI

**Status:** Planned

- Show provider health, capabilities, telemetry, profiles, and diagnostics.
- Consume service contracts only; never call providers directly.

**Exit criteria:** UI behavior remains correct with degraded and unavailable
providers.

## M7 — Operations and release

**Status:** Planned

- Package the daemon and UI.
- Add systemd integration, upgrade checks, and rollback procedures.
- Publish a compatibility matrix and operator documentation.

**Exit criteria:** installation, restart, upgrade, and rollback are tested in a
clean environment.

## M8 — Controlled hardware enablement

**Status:** Deferred

- Review any requested write operation separately.
- Define device, zone, limits, authorization, verification, and rollback.
- Start with a static low-brightness test only.

**Exit criteria:** explicit authorization and successful safety review. Mock or
automated tests do not satisfy this milestone.
