# ADR 0002: Implementation Language for the Local Service

- Status: Accepted for the service and adapter layers
- Date: 2026-09-15

## Decision

Use Go for the local daemon, provider adapters, health coordinator, and event
pipeline. Keep the desktop UI technology as a separate decision after the
service contracts stabilize.

## Rationale

- Go produces a straightforward native executable suitable for a local
  systemd-managed service.
- The language has built-in concurrency primitives appropriate for polling
  independent providers and isolating timeouts.
- Go modules provide versioned dependency management, and the standard tool
  includes build and test workflows.
- A small statically compiled service reduces runtime packaging requirements
  compared with a Python runtime deployment.
- Interfaces and explicit error handling fit the adapter and ownership
  boundaries defined in RFC 0001.

## Alternatives considered

### Rust

Rust provides stronger compile-time memory and thread-safety guarantees and a
strong Cargo workflow. It was not selected for the first slice because the
project's main risk is provider-contract and ownership design, while Rust would
add more implementation complexity before those contracts are proven.

### Python

Python's `asyncio` is well suited to I/O-bound provider coordination and its
packaging ecosystem is productive. It was not selected for the daemon because
runtime distribution and dependency isolation add operational work for a
desktop service that should be easy to install and roll back.

## Constraints

- The Go toolchain must be pinned in the repository's compatibility
  documentation before implementation is merged.
- Dependencies must be minimal, reviewed, and recorded through `go.mod` and
  `go.sum`.
- The first code slice remains read-only and mock-driven.
- No implementation may access HID, SMBus, I2C, firmware, or physical RGB
  devices without a separate approved design and test gate.

## Evidence

- Go documentation describes a statically typed compiled language with
  concurrency mechanisms and module-based development:
  <https://go.dev/doc/>.
- Go's module reference documents versioned modules and immutable version
  selection:
  <https://go.dev/ref/mod>.
- Python's official documentation describes `asyncio` as concurrent I/O for
  high-level network and IPC code:
  <https://docs.python.org/3/library/asyncio.html>.
- The Rust Book and Cargo documentation describe Rust's ownership-oriented
  language and Cargo package workflow:
  <https://doc.rust-lang.org/book/> and <https://doc.rust-lang.org/cargo/>.

## Revisit triggers

Reconsider this decision if provider SDKs require a different supported
runtime, if measured latency or memory violates the target platform, or if the
chosen UI architecture requires a shared language for a material maintenance
benefit.
