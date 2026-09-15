# Development Environment

## Local toolchain

The repository uses a project-local Go toolchain so development does not
depend on system packages or global PATH changes.

```bash
export PATH="$PWD/.toolchain/go-1.27.1/bin:$PATH"
go version
```

Expected version: `go1.27.1` on Linux amd64.

The toolchain archive was retrieved from the official Go distribution site and
verified against its published SHA-256 before extraction. The `.toolchain/`
directory is local-only and must not be committed.

## Validation commands

After the first Go module is added, run:

```bash
gofmt -w .
go test ./...
go vet ./...
git diff --check
```

Provider tests must use mocks. Do not run commands that write to HID, SMBus,
I2C, firmware, or physical RGB devices.

## Reproducibility

The selected toolchain version is recorded in the language decision. Go module
dependencies must be pinned through `go.mod` and `go.sum` when implementation
begins. Do not install dependencies globally for this project.
