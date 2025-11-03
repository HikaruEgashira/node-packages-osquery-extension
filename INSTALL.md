# Installation Guide

## Prerequisites

### Required

- **Go**: Version 1.21 or higher
  - Ubuntu/Debian: `sudo apt install golang-go`
  - macOS: `brew install go`
  - Or download from: https://go.dev/dl/

- **osquery**: Version 5.0.0 or higher
  - Ubuntu/Debian: `sudo apt install osquery`
  - macOS: `brew install osquery`
  - CentOS/RHEL: `sudo yum install osquery`

### Optional

- **golangci-lint**: For linting (development)
  - `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest`

## Building from Source

### Step 1: Clone the repository

```bash
git clone <repository-url>
cd osquery-extensions-node-packages
```

### Step 2: Download dependencies

```bash
make deps
```

Or manually:

```bash
go mod download
go mod tidy
```

### Step 3: Build the extension

```bash
make build
```

Or manually:

```bash
go build -o node_packages_extension .
```

### Step 4: Verify the build

```bash
ls -lh node_packages_extension
./node_packages_extension --help
```

## Running the Extension

### Method 1: Interactive mode with osqueryi

```bash
osqueryi --extension ./node_packages_extension
```

Then run queries:

```sql
osquery> SELECT * FROM node_packages LIMIT 5;
osquery> .schema node_packages
```

### Method 2: Daemon mode with osqueryd

1. Install the extension system-wide:

```bash
make install
```

Or manually:

```bash
sudo cp node_packages_extension /usr/local/bin/
```

2. Add to osquery extensions configuration:

Create or edit `/etc/osquery/extensions.load`:

```
/usr/local/bin/node_packages_extension
```

3. Restart osqueryd:

```bash
sudo systemctl restart osqueryd
```

### Method 3: Autoload configuration

Edit `/etc/osquery/osquery.conf`:

```json
{
  "options": {
    "extensions_autoload": "/etc/osquery/extensions.load",
    "extensions_timeout": "3",
    "extensions_interval": "3"
  }
}
```

## Testing

### Run unit tests

```bash
make test
```

### Run tests with coverage

```bash
make test-verbose
```

This generates `coverage.html` that you can open in a browser.

### Run specific tests

```bash
go test -v ./pkg/scanner -run TestScanAllManagers
```

## Verification

### Check if extension is loaded

```bash
osqueryi --extension ./node_packages_extension
```

```sql
SELECT * FROM osquery_extensions;
```

You should see `node_packages` extension listed.

### Test the table

```sql
SELECT * FROM node_packages LIMIT 5;
```

### Check what packages were found

```bash
# Run the scanner tests to see detected packages
go test -v ./pkg/scanner -run TestScanAllManagers
```

## Cross-Platform Building

### Build for Linux (from any platform)

```bash
GOOS=linux GOARCH=amd64 go build -o node_packages_extension_linux .
```

### Build for macOS (from any platform)

```bash
GOOS=darwin GOARCH=amd64 go build -o node_packages_extension_macos .
```

### Build for Windows (from any platform)

```bash
GOOS=windows GOARCH=amd64 go build -o node_packages_extension.exe .
```

### Build for ARM (Raspberry Pi, M1 Mac)

```bash
# ARM64 Linux
GOOS=linux GOARCH=arm64 go build -o node_packages_extension_arm64 .

# ARM64 macOS (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o node_packages_extension_m1 .
```

## Troubleshooting

### Extension fails to load

**Error**: `Failed to start extension`

**Solution**:
- Ensure osquery is installed: `osqueryi --version`
- Check extension permissions: `chmod +x node_packages_extension`
- Verify socket path: check osquery is running and socket exists

### No packages found

**Error**: Query returns no results

**Solution**:
```bash
# Run tests to see what's detected
make test

# Verify package manager caches exist
ls -la ~/.npm
ls -la ~/.pnpm-store
ls -la ~/.yarn-cache
ls -la ~/.cache/yarn
ls -la ~/.bun
ls -la ~/.cache/deno

# Install some packages to populate cache
npm install -g express
```

### Build errors

**Error**: `go: cannot find main module`

**Solution**:
```bash
# Ensure you're in the project directory
cd osquery-extensions-node-packages

# Regenerate go.mod if needed
go mod init github.com/HikaruEgashira/osquery-extensions-node-packages
make deps
```

**Error**: `package github.com/osquery/osquery-go: cannot find package`

**Solution**:
```bash
# Download dependencies
make deps

# Or manually
go get github.com/osquery/osquery-go
go mod tidy
```

### Test failures

**Error**: Tests fail or timeout

**Solution**:
```bash
# Run tests with verbose output
go test -v -timeout 30s ./...

# Run specific failing test
go test -v ./pkg/scanner -run TestName
```

### Runtime errors

**Error**: Extension crashes

**Solution**:
```bash
# Build with race detector
go build -race -o node_packages_extension .

# Run with verbose logging
./node_packages_extension --socket /path/to/osquery.sock -verbose
```

## Platform-Specific Notes

### macOS

- Default cache locations may differ
- pnpm store: `~/Library/pnpm/store`
- Ensure osquery is installed via Homebrew for best compatibility

### Linux

- Ensure file permissions allow reading cache directories
- Check environment variables: `PNPM_HOME`, `DENO_DIR`, `BUN_INSTALL`
- For systemd: enable osqueryd service

### Windows

- Use PowerShell or Command Prompt
- Cache paths differ: `%APPDATA%\npm`, etc.
- Extension socket path format differs

## Development Setup

### Install development tools

```bash
# Install linter
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Install test coverage tools
go install golang.org/x/tools/cmd/cover@latest
```

### Run linter

```bash
make lint
```

### Format code

```bash
make fmt
```

### Run all quality checks

```bash
make fmt
make lint
make test-verbose
```

## Uninstallation

### Remove system-wide installation

```bash
make uninstall
```

Or manually:

```bash
sudo rm /usr/local/bin/node_packages_extension
```

### Remove from osquery configuration

```bash
# Edit extensions.load
sudo nano /etc/osquery/extensions.load
# Remove the line with node_packages_extension

# Restart osqueryd
sudo systemctl restart osqueryd
```

### Clean build artifacts

```bash
make clean
```

## Docker Deployment

### Create Dockerfile

```dockerfile
FROM ubuntu:latest

RUN apt-get update && apt-get install -y \
    osquery \
    golang-go

WORKDIR /app
COPY . .

RUN make build

CMD ["osqueryi", "--extension", "./node_packages_extension"]
```

### Build and run

```bash
docker build -t osquery-node-packages .
docker run -it osquery-node-packages
```

## Support

For issues and questions:
- Check the main [README.md](README.md)
- Review [example_queries.sql](example_queries.sql) for usage examples
- Run `make test` to verify installation
- Submit issues to the repository
