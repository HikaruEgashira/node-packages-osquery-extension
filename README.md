# osquery Node Packages Extension

osquery extension written in Go to discover installed Node.js packages by scanning package manager caches.

## Features

This extension provides a `node_packages` table that scans various package manager caches to discover installed Node.js packages:

- **npm**: Scans `~/.npm`, global node_modules
- **pnpm**: Scans `~/.pnpm-store`, `~/.local/share/pnpm/store`
- **yarn**: Scans `~/.yarn-cache`, `~/.cache/yarn` (Yarn v1)
- **bun**: Scans `~/.bun/install/cache`
- **deno**: Scans `~/.cache/deno/npm`, `$DENO_DIR`
- **jsr**: Scans JSR packages through Deno cache

### Why Go?

- ✅ **Concurrent scanning**: Scans all package managers in parallel using goroutines
- ✅ **Cross-platform**: Easy to compile for Linux, macOS, Windows
- ✅ **Single binary**: No runtime dependencies
- ✅ **Standard JSON parsing**: Built-in encoding/json package
- ✅ **Memory safe**: Automatic garbage collection

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
make deps

# Build the extension
make build

# Run tests
make test
```

Or build manually:

```bash
go build -o node_packages_extension .
```

## Installation

### Quick Start

1. Build the extension:
```bash
make build
```

2. Run osquery with the extension:
```bash
osqueryi --extension ./node_packages_extension
```

### System-wide Installation

```bash
make install
```

Then add to `/etc/osquery/extensions.load`:
```
/usr/local/bin/node_packages_extension
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

## Testing

### Run all tests

```bash
make test
```

### Run tests with coverage

```bash
make test-verbose
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

## Troubleshooting

### No packages found

```bash
# Verify package manager caches exist
ls -la ~/.npm
ls -la ~/.pnpm-store
ls -la ~/.bun

# Run tests to see what's detected
make test
```

### Extension fails to load

```bash
# Check osquery is running
osqueryi --version

# Run with verbose logging
./node_packages_extension --socket /path/to/osquery.sock
```

### Build errors

```bash
# Ensure Go is installed
go version

# Clean and rebuild
make clean
make deps
make build
```

## Security Considerations

- ✅ Read-only operations
- ✅ Permission errors handled gracefully
- ✅ No network access
- ✅ No cache modification
- ✅ Memory-safe Go implementation

## License

MIT

## Contributing

Contributions are welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure `make test` passes
5. Submit a pull request

## Example Queries

See [example_queries.sql](example_queries.sql) for more SQL query examples.
