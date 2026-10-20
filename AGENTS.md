# AGENTS.md

Guidance for automated agents and contributors working in this repository.

## Project

osquery extension (Go) that discovers installed Node.js packages by scanning
package manager caches. Entry point: `main.go`. Core logic: `pkg/scanner`.

## Requirements

- Go 1.25+ (see `go.mod`)
- osquery (only needed for manual end-to-end testing)

## Commands

```bash
# Download dependencies
go mod download

# Verify module checksums
go mod verify

# Format
go fmt ./...

# Vet / static checks
go vet ./...

# Build the extension
go build ./...

# Run all tests
go test -v -race ./...
```

## Validation

CI (`.github/workflows/ci.yml`) runs on every push to `main` and on every pull
request. It verifies dependencies, checks that `go.mod`/`go.sum` are tidy,
vets, builds, and runs the test suite with the race detector. Keep `go.mod`
and `go.sum` tidy (`go mod tidy`) so the check stays green.

Do not weaken tests or disable checks to make CI pass.
