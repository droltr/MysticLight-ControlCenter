# MysticLight Control Center

MysticLight Control Center is a local-first coordination service for RGB
lighting, cooling telemetry, LCD status, sensors, profiles, and game-focus
events. It provides one consistent control plane while leaving existing
hardware providers in charge of their devices.

## Why this project exists

Desktop hardware-control tools often expose overlapping capabilities without a
shared ownership model. That can produce conflicting writes, unclear failure
states, and difficult recovery. This project adds a provider-neutral domain
layer, explicit ownership rules, health reporting, profiles, diagnostics, and
safe fallback behavior.

## Project scope

### Included

- Provider discovery and capability reporting.
- Normalized RGB, cooling, LCD, sensor, and game-focus state.
- Explicit ownership for RGB, fan/PWM, LCD, telemetry, and profiles.
- Validated profiles with conflict detection.
- Provider health, timeouts, redacted diagnostics, and recovery events.
- A local daemon and event pipeline.
- A desktop UI consuming service contracts.
- Compatibility documentation, packaging, systemd integration, upgrades, and
  rollback.

### Excluded

- New hardware drivers or firmware changes.
- Experimental raw HID, SMBus, or I2C operations.
- Multiple writers for the same physical resource.
- Cloud accounts, remote control, or mandatory internet access.
- Credentials, serial numbers, addresses, stable device identifiers, or raw
  host identity in configuration, logs, tests, or documentation.

## Ownership and safety model

| Resource | Provider owner | Control Center responsibility |
| --- | --- | --- |
| RGB lighting | OpenRGB | Capability, profile, health, and policy coordination |
| Fan and pump PWM | CoolerControl/CoolerDash | Read-only telemetry and ownership enforcement |
| LCD display | LCD provider | Status normalization and recovery policy |
| Sensor telemetry | Existing sensor provider | Normalization and diagnostics |
| Game focus | Focus adapter | Local event production |

The first implementation stages are read-only and mock-driven. Automated tests
must never write to physical hardware. Any later physical test requires
explicit authorization, a selected device and zone, low-brightness startup,
and a tested rollback path.

## Development milestones

1. **Contracts and governance** — scope, ownership, privacy, architecture, and
   language decisions.
2. **Core domain** — capabilities, resources, profiles, validation, and
   provider registry.
3. **Service and events** — health, timeouts, normalized observations,
   redacted diagnostics, and recovery events.
4. **Read-only adapters** — OpenRGB, cooling, LCD, sensors, and focus provider
   integrations behind the contracts.
5. **Profile runtime** — validated activation, safe fallback, and restart
   recovery without physical writes.
6. **Desktop UI** — status, profiles, diagnostics, and capability visibility.
7. **Operations** — packaging, systemd service, upgrades, rollback, and
   compatibility matrix.
8. **Controlled hardware enablement** — separately reviewed and authorized
   provider writes, if required.

Each milestone requires tests, a privacy review, documented safety impact, and
rollback notes before it is considered complete.

## Current implementation

The service and adapter layers use Go 1.27.1 with a project-local toolchain.
The repository currently contains the initial domain/provider contracts,
profile validation, a mock provider, and unit tests. Real provider connections,
the daemon runtime, and the UI are not implemented yet.

## Repository map

- `internal/domain/` — provider-neutral resources, capabilities, ownership,
  health, and profile rules.
- `internal/provider/` — adapter contracts, registry, and test doubles.
- `docs/PROJECT_SCOPE.md` — detailed scope and phase exit criteria.
- `docs/ARCHITECTURE_RFC.md` — layered architecture and safety contract.
- `docs/IMPLEMENTATION_PLAN.md` — first implementation slices.
- `docs/RELEASE_NOTES.md` — recorded project milestones and unreleased work.
- `DEVELOPMENT.md` — local toolchain and validation commands.

## Validation

```bash
export PATH="$PWD/.toolchain/go-1.27.1/bin:$PATH"
gofmt -w .
go test ./...
go vet ./...
git diff --check
```

Provider tests use mocks and do not claim physical hardware verification.
