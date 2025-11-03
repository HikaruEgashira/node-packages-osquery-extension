# osquery Node Packages Extension

osquery extension written in Go to discover installed Node.js packages by scanning package manager caches.

## Quick Start

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from [Releases](https://github.com/HikaruEgashira/node-packages-osquery-extension/releases).

```bash
# Run with osquery
osqueryi --extension ./node_packages_extension
```

Then query your packages:
```sql
SELECT * FROM node_packages;
```

### Option 2: Build from Source

```bash
# Clone and build
git clone https://github.com/HikaruEgashira/node-packages-osquery-extension.git
cd node-packages-osquery-extension
go build -o node_packages_extension .

# Run with osquery
osqueryi --extension ./node_packages_extension
```

## Features

This extension provides a `node_packages` table that scans various package manager caches to discover installed Node.js packages:

- **npm**: Scans `~/.npm`, global node_modules
- **pnpm**: Scans `~/.pnpm-store`, `~/.local/share/pnpm/store`
- **yarn**: Scans `~/.yarn-cache`, `~/.cache/yarn` (Yarn v1)
- **bun**: Scans `~/.bun/install/cache`
- **deno**: Scans `~/.cache/deno/npm`, `$DENO_DIR`
- **jsr**: Scans JSR packages through Deno cache

## Table Schema

```sql
CREATE TABLE node_packages (
    name TEXT,
    version TEXT,
    manager TEXT,
    cache_path TEXT
);
```

### Columns

- `name`: Package name
- `version`: Package version
- `manager`: Package manager (npm, pnpm, yarn, bun, deno, jsr)
- `cache_path`: Path to the package.json in cache

## Installation

### System-wide Installation

```bash
# Build and install to /usr/local/bin
go build -o node_packages_extension .
sudo cp node_packages_extension /usr/local/bin/
sudo chown root:root /usr/local/bin/node_packages_extension
sudo chmod 755 /usr/local/bin/node_packages_extension
```

Then add to `/etc/osquery/extensions.load`:
```
/usr/local/bin/node_packages_extension
```

## Building

### Prerequisites

- Go 1.21 or higher
- osquery installed on your system

### Build Instructions

```bash
# Clone the repository
git clone https://github.com/HikaruEgashira/node-packages-osquery-extension.git
cd node-packages-osquery-extension

# Download dependencies
go mod download

# Build the extension
go build -o node_packages_extension .

# Run tests
go test -v ./...
```

## Usage

### Query all packages

```sql
SELECT * FROM node_packages;
```

### Query packages by manager

```sql
SELECT name, version FROM node_packages WHERE manager = 'npm';
```

### Count packages per manager

```sql
SELECT manager, COUNT(*) as count FROM node_packages GROUP BY manager;
```

### Find specific package

```sql
SELECT * FROM node_packages WHERE name LIKE '%react%';
```

### List unique packages (across all managers)

```sql
SELECT DISTINCT name, version FROM node_packages ORDER BY name;
```

See [example_queries.sql](example_queries.sql) for more SQL query examples.

## Testing

### Run all tests

```bash
go test -v ./...
```

### Run tests with coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

This will generate a `coverage.html` file that you can open in your browser.

### Run specific tests

```bash
go test -v ./pkg/scanner -run TestScanAllManagers
```

## Cache Locations

The extension scans the following default locations:

### npm
- `~/.npm`
- `/usr/local/lib/node_modules`
- `/usr/lib/node_modules`
- Global node_modules paths

### pnpm
- `~/.pnpm-store`
- `~/.local/share/pnpm/store`
- `~/Library/pnpm/store` (macOS)
- `$PNPM_HOME/store`

### yarn
- `~/.yarn-cache` (Yarn v1 old)
- `~/.cache/yarn` (Yarn v1 new)
- `~/Library/Caches/Yarn` (macOS)
- `$YARN_CACHE_FOLDER`

### bun
- `~/.bun/install/cache`
- `~/.bun/install/global`

### deno
- `$DENO_DIR/npm`
- `$DENO_DIR/deps/https`
- `~/.cache/deno` (default)

### jsr
- `$DENO_DIR/deps/https/jsr.io`
- `~/.cache/deno/deps/https/jsr.io` (default)

## Performance

The extension uses efficient scanning:
- **Concurrent processing**: All package managers scanned in parallel
- **Permission-safe**: Gracefully handles permission errors
- **Memory efficient**: Streams large directories
- **Fast JSON parsing**: Standard library performance

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

MIT
