.PHONY: all build clean install test fmt lint run

BINARY_NAME=node_packages_extension
GO=go
GOFLAGS=-v

all: build

build:
	@echo "Building $(BINARY_NAME)..."
	$(GO) build $(GOFLAGS) -o $(BINARY_NAME) .

clean:
	@echo "Cleaning..."
	@rm -f $(BINARY_NAME)
	@$(GO) clean

test:
	@echo "Running tests..."
	$(GO) test -v ./...

test-verbose:
	@echo "Running tests with verbose output..."
	$(GO) test -v -race -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html

fmt:
	@echo "Formatting code..."
	$(GO) fmt ./...

lint:
	@echo "Running linter..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	golangci-lint run ./...

deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	$(GO) mod tidy

run:
	@echo "Note: This extension requires osquery to be running."
	@echo "Usage: osqueryi --extension ./$(BINARY_NAME)"
	@echo ""
	@echo "Or run directly with:"
	@echo "./$(BINARY_NAME) --socket <path_to_osquery_socket>"

install: build
	@echo "Installing $(BINARY_NAME)..."
	@sudo cp $(BINARY_NAME) /usr/local/bin/

uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	@sudo rm -f /usr/local/bin/$(BINARY_NAME)

example-queries:
	@echo "Example queries for node_packages table:"
	@echo ""
	@echo "1. List all packages:"
	@echo "   SELECT * FROM node_packages;"
	@echo ""
	@echo "2. Count packages by manager:"
	@echo "   SELECT manager, COUNT(*) as count FROM node_packages GROUP BY manager;"
	@echo ""
	@echo "3. Find React packages:"
	@echo "   SELECT * FROM node_packages WHERE name LIKE '%react%';"
	@echo ""
	@echo "4. List npm packages:"
	@echo "   SELECT name, version FROM node_packages WHERE manager = 'npm';"
	@echo ""
	@echo "5. Show unique packages:"
	@echo "   SELECT DISTINCT name, COUNT(DISTINCT version) as versions FROM node_packages GROUP BY name;"

help:
	@echo "Available targets:"
	@echo "  build            - Build the extension"
	@echo "  clean            - Remove build artifacts"
	@echo "  test             - Run tests"
	@echo "  test-verbose     - Run tests with coverage"
	@echo "  fmt              - Format Go code"
	@echo "  lint             - Run linter"
	@echo "  deps             - Download and tidy dependencies"
	@echo "  install          - Install the extension (requires sudo)"
	@echo "  uninstall        - Uninstall the extension (requires sudo)"
	@echo "  run              - Show run instructions"
	@echo "  example-queries  - Show example SQL queries"
	@echo "  help             - Show this help message"
