# MysticLight Control Center

Unified, safety-first control plane for RGB lighting, cooling telemetry,
profiles, and LCD status. Existing projects remain the device providers; this
application coordinates them through explicit adapters and ownership rules.

## Goals

- One installation and one user interface.
- Coordinate OpenRGB, CoolerControl/CoolerDash, LCD telemetry, sensors, and game focus.
- Keep fan/PWM, LCD, and RGB ownership isolated to prevent write conflicts.
- Detect missing capabilities as warnings and preserve safe fallback behavior.
- Provide health, diagnostics, rollback, and restart recovery.

## Ordered roadmap

1. Architecture RFC and ownership contract.
2. Shared configuration/profile schema.
3. Provider adapters (OpenRGB, CoolerControl, LCD, sensors, focus).
4. Unified health/diagnostics service.
5. Local daemon and event pipeline.
6. Desktop UI for status, profiles, cooling, RGB, and LCD.
7. Single installer, systemd integration, upgrades, and rollback.
8. Compatibility matrix and hardware-specific GPU capability probe.
9. End-to-end reboot, suspend/resume, and failure-recovery tests.

## Safety

The control center does not reimplement hardware drivers. CoolerControl remains
the fan/pump PWM owner, OpenRGB the RGB owner, and the LCD provider the display
owner. Secrets, serials, tokens, and raw host identifiers must never be stored.
