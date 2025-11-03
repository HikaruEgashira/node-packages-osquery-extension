# Contributing to node-packages-osquery-extension

Thank you for considering contributing to this project! This document provides guidelines for contributing.

## How to Contribute

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass (`make test`)
6. Commit your changes with conventional commits
7. Push to your fork (`git push origin feature/amazing-feature`)
8. Open a Pull Request

## Development Setup

### Prerequisites

- Go 1.21 or higher
- osquery installed on your system

### Build and Test

```bash
# Install dependencies
make deps

# Build the extension
make build

# Run tests
make test

# Run tests with coverage
make test-verbose

# Format code
make fmt
```

### Project Structure

```
.
├── main.go              # Extension entry point
├── pkg/
│   └── scanner/
│       ├── scanner.go      # Core scanning logic
│       └── scanner_test.go # Unit tests
├── Makefile
├── go.mod
└── README.md
```

## Code Guidelines

### Performance Requirements

**CRITICAL: Directory Walking Restrictions**

- Avoid recursive directory traversal where possible
- `filepath.WalkDir()` should be used judiciously as it can cause performance issues
- Be mindful of scanning large directory trees

**Rationale:**
- Performance can degrade significantly (e.g., 84 seconds for large directories)
- Poor performance degrades user experience

**Alternatives:**
- Use indexed data when available
- Leverage existing package manager cache structures
- Accept feature limitations if no efficient implementation exists

### Code Quality

- Write clear, idiomatic Go code
- Add unit tests for new functionality
- Maintain test coverage above 80%
- Use meaningful variable and function names
- Add comments for complex logic

### Testing Requirements

All contributions must include:
- Unit tests for new functions
- Integration tests for new features
- Test coverage should not decrease

## Adding a New Package Manager

To add support for a new package manager:

1. Add scanner function in `pkg/scanner/scanner.go`:
```go
func ScanYourManager() ([]Package, error) {
    home, _ := os.UserHomeDir()
    if home == "" {
        return []Package{}, nil
    }
    // Implementation
    return scanPaths("yourmanager", paths)
}
```

2. Register in `ScanAllManagers()`:
```go
{"yourmanager", ScanYourManager},
```

3. Add comprehensive tests in `scanner_test.go`:
```go
func TestScanYourManager(t *testing.T) {
    // Test implementation
}
```

4. Update README.md with the new package manager information

## Code Review Process

All submissions require review. We use GitHub pull requests for this purpose. The review process includes:

- Code quality and style check
- Test coverage verification
- Performance impact assessment
- Documentation completeness

## Commit Message Guidelines

Use conventional commits format:

- `feat: add support for new package manager`
- `fix: correct path resolution on Windows`
- `docs: update installation instructions`
- `test: add tests for edge cases`
- `refactor: improve error handling`

## Questions or Issues?

- Open an issue for bugs or feature requests
- Use discussions for questions and ideas
- Check existing issues before creating new ones

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
