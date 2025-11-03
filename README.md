# osquery Node Packages Extension

osquery extension written in Go to discover installed Node.js packages by scanning package manager caches.

## Installation

```bash
mise use -g github:HikaruEgashira/node-packages-osquery-extension
```

For other installation methods, see the [releases page](https://github.com/HikaruEgashira/node-packages-osquery-extension/releases).

## Quick Start

```bash
# Run with osquery
osqueryi --extension ./node_packages_extension
```

Then query your packages:
```sql
SELECT * FROM node_packages;
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

### System-wide Installation

```bash
# Move binary to /usr/local/bin
sudo mv node_packages_extension /usr/local/bin/
sudo chown root:root /usr/local/bin/node_packages_extension
sudo chmod 755 /usr/local/bin/node_packages_extension
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

See [example_queries.sql](example_queries.sql) for more SQL query examples.

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

## Building from Source

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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

MIT
