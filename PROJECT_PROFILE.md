# Project Profile

## Purpose

MysticLight Control Center is a local-first coordination layer for RGB
lighting, cooling telemetry, profiles, and LCD status. It gives users one
interface while preserving ownership boundaries between established device
providers.

## Scope

- Coordinate OpenRGB, CoolerControl/CoolerDash, LCD, sensor, and game-focus
  providers through explicit adapters.
- Provide shared profiles, health and diagnostics, recovery, and safe fallback
  behavior.
- Keep fan/PWM, RGB, and LCD ownership isolated to prevent conflicting writes.

## Non-goals

- Reimplement hardware drivers or replace the existing device providers.
- Perform experimental raw HID, SMBus, or I2C operations.
- Store credentials, serial numbers, tokens, or raw host identifiers.
- Claim physical hardware verification from mocked or automated tests.

## Operating context

The project targets a local desktop environment that may coordinate RGB
controllers, cooling services, sensors, LCD devices, and game-focus events.
Provider availability is expected to vary; missing capabilities must remain
warnings with safe fallback behavior.

## Ownership and safety

- CoolerControl remains the fan and pump PWM owner.
- OpenRGB remains the RGB owner.
- The LCD provider remains the display owner.
- Automated tests must not write to physical hardware.
- Any future physical test must be explicitly authorized, limited to the
  selected device and zone, begin with a static low-brightness color, and have
  a documented rollback path.

## Platform and technology decision

The repository currently contains architecture documentation only, so no
implementation language is selected yet. The first implementation decision
must compare candidates against local desktop integration, provider protocol
support, testability, concurrency, packaging, dependency provenance, and
hardware-safety constraints. The decision record should be added before
implementation files are introduced.

## Validation

Validation will combine unit and integration tests with provider mocks,
configuration/schema checks, `git diff --check`, secret and identifier scans,
and review of hardware-affecting diffs. Physical verification, when separately
authorized, must be reported distinctly from simulated verification.

## Privacy and rollback

Diagnostics and documentation must use neutral placeholders such as `<device>`
and `<host>` and must exclude usernames, hostnames, addresses, serials,
credentials, and stable hardware identifiers. Provider changes must preserve
the existing provider services and include a tested rollback path before
replacement or activation.

## Current risks and review triggers

- Provider API and capability differences may affect adapter design.
- Ownership conflicts could cause unsafe or surprising device behavior.
- Hardware-specific behavior requires explicit compatibility evidence.
- A language and architecture decision is required before implementation.
