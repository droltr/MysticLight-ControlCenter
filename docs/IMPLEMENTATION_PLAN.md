# Implementation Plan

## Immediate work

1. Approve RFC 0001 and record the implementation-language decision.
2. Define the profile schema and ownership/capability vocabulary.
3. Build provider interfaces with mock implementations only.
4. Add deterministic validation and failure-injection tests.
5. Implement read-only health and telemetry collection.

## First implementation slice

The first code change should include only:

- domain types for capability, health, telemetry, ownership, and profile;
- provider adapter interfaces;
- a mock provider;
- profile validation;
- tests proving unavailable providers degrade safely;
- no physical provider writes and no UI dependency.

## Definition of done for the first slice

- The selected language and build tool are documented.
- Tests run without provider services or hardware.
- Invalid ownership and profile states are rejected deterministically.
- Provider timeouts and malformed responses are represented as safe health
  failures.
- No credentials or stable device identifiers enter logs or fixtures.
- `main` remains unchanged and the work is delivered through a focused topic
  branch and pull request.
