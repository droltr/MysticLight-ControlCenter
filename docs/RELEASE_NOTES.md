# Release Notes

This file records repository milestones and candidate release content. A
section marked `Unreleased` is not a published release and does not imply that
the project is ready for production or physical hardware use.

## Unreleased

### Added

- Defined the project purpose, intended users, scope, non-goals, MVP boundary,
  delivery phases, success criteria, and review gates.
- Documented the safe local coordination architecture and provider ownership
  contract.
- Selected Go 1.27.1 for the local daemon, provider adapters, health
  coordinator, and event pipeline.
- Added a project-local Go development toolchain and development instructions.
- Added the initial Go module with domain ownership, capability, profile
  validation, and provider registry contracts.
- Added a mock provider and deterministic unit tests for ownership and provider
  validation.

### Safety and privacy

- Kept RGB, PWM, LCD, HID, SMBus, I2C, firmware, and physical device writes
  outside the first implementation slice.
- Preserved existing provider ownership boundaries.
- Excluded credentials, serial numbers, addresses, and stable device
  identifiers from the new artifacts and tests.
- Established mock-only validation as the first implementation gate.

### Validation

- `go test ./...` passes.
- `go vet ./...` passes.
- `git diff --check` passes.
- No physical hardware verification has been performed.

### Known limitations

- Provider adapters currently define contracts only; no real provider
  connection is implemented.
- There is no service/event runtime or desktop UI yet.
- No compatibility matrix, installer, systemd integration, or rollback
  automation has been implemented.
- CI checks have not yet reported for the initial pull request.

### Next planned milestone

Implement the local service and event contracts with provider timeouts,
redacted health diagnostics, failure injection, and mock-only recovery tests.

## Repository milestones

| Commit | Milestone |
| --- | --- |
| `6b47501` | Established the control center architecture roadmap. |
| `d023082` | Added the English project profile. |
| `35b2538` | Defined project scope, architecture RFC, and implementation plan. |
| `2c4892c` | Recorded the Go implementation-language decision. |
| `d462606` | Prepared the local Go development environment. |
| `8f806fa` | Added core provider and domain contracts. |
| `42bab29` | Added mock provider contract coverage. |
