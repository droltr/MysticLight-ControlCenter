# Operations

## Current status

The repository contains a systemd user-service template for the daemon. It is
not installed or enabled automatically. Packaging must place the compiled
`mysticlightd` binary at the path declared by the unit or generate a unit with
the selected installation path.

## Installation sequence

1. Build the daemon with the pinned project-local Go toolchain.
2. Verify `go test ./...` and `go vet ./...` in a clean checkout.
3. Install the binary and unit file through the selected package format.
4. Reload the user systemd manager.
5. Start the service in read-only mode and inspect health logs.
6. Enable automatic startup only after restart behavior is verified.

## Rollback sequence

1. Stop the user service.
2. Disable the unit if it was enabled.
3. Restore the previous binary and configuration/profile files.
4. Start the previous service version.
5. Confirm that provider ownership and safe fallback behavior remain intact.

## Safety requirements

- Do not grant the service access to raw HID, SMBus, I2C, or firmware devices
  during the read-only milestones.
- Keep `NoNewPrivileges`, `PrivateTmp`, and `ProtectSystem` enabled unless a
  reviewed requirement proves otherwise.
- Store only validated profiles and redacted diagnostics in the service state
  directory.
- Do not run installation or service commands against physical hardware as
  part of automated tests.

## Release gate

Operations work is not complete until a clean-environment install, restart,
upgrade, and rollback test has been recorded. A unit-file template alone does
not satisfy that gate.
