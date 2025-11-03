# osquery Node Packages Extension

osquery extension written in Go to discover installed Node.js packages by scanning package manager caches.

## Installation

```bash
mise use -g github:HikaruEgashira/node-packages-osquery-extension

# Linux
sudo chown root:root $(which node_packages_extension)
sudo chmod 755 $(which node_packages_extension)

# macOS
sudo chown root:wheel $(which node_packages_extension)
sudo chmod 755 $(which node_packages_extension)
```

For other installation methods, see the [releases page](https://github.com/HikaruEgashira/node-packages-osquery-extension/releases).

## Quick Start

```bash
osqueryi --extension $(which node_packages_extension)
> SELECT * FROM node_packages;
```

## Features

| Package Manager | Cache Locations                                      | Supported |
|-----------------|------------------------------------------------------|-----------|
| npm             | `~/.npm`, global `node_modules`                      | Yes       |
| pnpm            | `~/.pnpm-store`, `~/.local/share/pnpm/store`         | Yes       |
| yarn            | `~/.yarn-cache`, `~/.cache/yarn` (Yarn v1)           | Yes       |
| bun             | `~/.bun/install/cache`                               | Yes       |
| deno            | `~/.cache/deno/npm`, `$DENO_DIR`                     | Yes       |

## Table Schema

```sql
CREATE TABLE node_packages (
    name TEXT,
    version TEXT,
    manager TEXT,
    cache_path TEXT
);
```

- `name`: Package name
- `version`: Package version
- `manager`: Package manager (npm, pnpm, yarn, bun, deno)
- `cache_path`: Path to the package.json in cache

## Query Examples

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

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

## License

MIT
