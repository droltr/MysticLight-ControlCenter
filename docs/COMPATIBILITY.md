# Compatibility Matrix

This matrix records the planned provider boundary. A provider is not marked
supported until its API/version behavior is verified with read-only contract
tests. No entry in this document authorizes physical writes.

| Provider | Resource | Initial mode | Write operations | Status |
| --- | --- | --- | --- | --- |
| OpenRGB | RGB | Adapter boundary only | Disabled | Contract scaffold |
| CoolerControl/CoolerDash | PWM and cooling telemetry | Adapter boundary only | Disabled | Contract scaffold |
| LCD provider | LCD | Adapter boundary only | Disabled | Contract scaffold |
| Sensor provider | Telemetry | Adapter boundary only | Disabled | Contract scaffold |
| Game-focus adapter | Focus events | Adapter boundary only | Not applicable | Contract scaffold |

## Compatibility evidence required

Before an adapter is marked supported, record:

- Provider name and documented API or protocol version.
- Supported platform and transport.
- Capability mapping and timeout behavior.
- Malformed-response and unavailable-provider behavior.
- Mock contract tests and redacted diagnostics review.
- Rollback and disable procedure.

## Safety boundary

The current adapter package exposes observation capability only. Provider
constructors return a safe not-configured implementation until a provider
contract is reviewed and implemented. Tests must not start provider services,
open hardware handles, or write to devices.
