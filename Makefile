# Quiver Makefile
# Common tasks for development and deployment

# Variables
BINARY_NAME=quiver
BUILD_DIR=bin
MAIN_PATH=./cmd/quiver
DOCKER_IMAGE=quiver:latest
VERSION?=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-X main.version=${VERSION} -X main.buildTime=$(shell date -u +%Y-%m-%dT%H:%M:%SZ) -X main.gitCommit=$(shell git rev-parse --short HEAD)"

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

.PHONY: all build clean test coverage deps fmt lint vet help install uninstall docker-build docker-run

# Default target
all: clean deps test build

# Build the application
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@echo "Copying embedded files..."
	@cp config.json internal/config/config.json
	@cp metadata.json internal/metadata/metadata.json
	$(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Cleaning up embedded files..."
	@rm -f internal/config/config.json internal/metadata/metadata.json
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "Copying embedded files..."
	@cp config.json internal/config/config.json
	@cp metadata.json internal/metadata/metadata.json
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Cleaning up embedded files..."
	@rm -f internal/config/config.json internal/metadata/metadata.json
	@echo "Multi-platform build complete"

# Clean build artifacts
clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f coverage.out
	@echo "Clean complete"

# Run tests
test:
	@echo "Running tests..."
	$(GOTEST) -race -v ./...

# Run tests with coverage
coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -race -coverprofile=coverage.out -covermode=atomic ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

# Format code
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .

# Lint code
lint:
	@echo "Linting code..."
	@if ! command -v $(GOLINT) > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	fi
	$(GOLINT) run ./...

# Vet code
vet:
	@echo "Vetting code..."
	$(GOCMD) vet ./...

# Run all checks
check: fmt vet lint test

# Install the binary
install: build
	@echo "Installing $(BINARY_NAME)..."
	cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Install complete"

# Uninstall the binary
uninstall:
	@echo "Uninstalling $(BINARY_NAME)..."
	rm -f $(GOPATH)/bin/$(BINARY_NAME)
	@echo "Uninstall complete"

# Run the application
run: build
	@echo "Running $(BINARY_NAME)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run in development mode
dev:
	@echo "Running in development mode..."
	$(GOCMD) run $(MAIN_PATH)

# Generate documentation
docs:
	@echo "Generating documentation..."
	@if ! command -v godoc > /dev/null; then \
		echo "Installing godoc..."; \
		$(GOGET) golang.org/x/tools/cmd/godoc; \
	fi
	@echo "Starting documentation server on :6060"
	godoc -http=:6060

# Docker build
docker-build:
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .
	@echo "Docker build complete: $(DOCKER_IMAGE)"

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	docker run -p 8080:8080 $(DOCKER_IMAGE)

# Create release
release: clean deps test build-all
	@echo "Creating release $(VERSION)..."
	@mkdir -p release
	tar -czf release/$(BINARY_NAME)-$(VERSION)-linux-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-linux-amd64
	tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-amd64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-amd64
	tar -czf release/$(BINARY_NAME)-$(VERSION)-darwin-arm64.tar.gz -C $(BUILD_DIR) $(BINARY_NAME)-darwin-arm64
	zip -j release/$(BINARY_NAME)-$(VERSION)-windows-amd64.zip $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe
	@echo "Release created in release/ directory"

# Setup development environment
setup:
	@echo "Setting up development environment..."
	$(GOMOD) download
	@if ! command -v golangci-lint > /dev/null; then \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2; \
	fi
	@echo "Development environment setup complete"

# Update dependencies
update:
	@echo "Updating dependencies..."
	$(GOGET) -u ./...
	$(GOMOD) tidy
	@echo "Dependencies updated"

# Security scan
security:
	@echo "Running security scan..."
	@if ! command -v gosec > /dev/null; then \
		echo "Installing gosec..."; \
		$(GOGET) github.com/securecodewarrior/gosec/v2/cmd/gosec; \
	fi
	gosec ./...

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  build-all   - Build for multiple platforms"
	@echo "  clean       - Clean build artifacts"
	@echo "  test        - Run tests"
	@echo "  coverage    - Run tests with coverage"
	@echo "  deps        - Install dependencies"
	@echo "  fmt         - Format code"
	@echo "  lint        - Lint code"
	@echo "  vet         - Vet code"
	@echo "  check       - Run all checks (fmt, vet, lint, test)"
	@echo "  install     - Install binary"
	@echo "  uninstall   - Uninstall binary"
	@echo "  run         - Build and run"
	@echo "  dev         - Run in development mode"
	@echo "  docs        - Generate and serve documentation"
	@echo "  docker-build- Build Docker image"
	@echo "  docker-run  - Build and run Docker container"
	@echo "  release     - Create release packages"
	@echo "  setup       - Setup development environment"
	@echo "  update      - Update dependencies"
	@echo "  security    - Run security scan"
	@echo "  bench       - Run benchmark tests"
	@echo "  help        - Show this help message" 