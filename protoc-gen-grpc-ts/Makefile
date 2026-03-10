# Makefile for protoc-gen-grpc-ts plugin

.PHONY: all build install test test-unit test-integration test-coverage test-clean clean deps fmt lint help

# Default target - build the plugin
all: build

# Build the plugin
build:
	@echo "Building protoc-gen-grpc-ts plugin..."
	@go build -o protoc-gen-grpc-ts main.go
	@chmod +x protoc-gen-grpc-ts
	@echo "Build complete: ./protoc-gen-grpc-ts"

# Install the plugin to a directory in PATH (default: /usr/local/bin)
install: build
	@echo "Installing plugin..."
	@INSTALL_DIR=${INSTALL_DIR:-/usr/local/bin}; \
	if [ -w $$INSTALL_DIR ]; then \
		cp protoc-gen-grpc-ts $$INSTALL_DIR/; \
		echo "Installed to $$INSTALL_DIR/protoc-gen-grpc-ts"; \
	else \
		echo "$$INSTALL_DIR is not writable. Run with sudo or set INSTALL_DIR:"; \
		echo "   sudo make install"; \
		echo "   or"; \
		echo "   make install INSTALL_DIR=~/bin"; \
	fi

# Run all tests
test: test-unit test-integration

# Run unit tests
test-unit:
	@echo "Running unit tests..."
	go test -v -run "TestFieldNaming|TestNamespaceGeneration|TestTypeScriptTypeMapping|TestFileGeneration|TestErrorHandling"

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	go test -v -run "TestComprehensiveIntegration|TestFieldNamingConventions|TestNamespaceGeneration|TestImportHandling|TestServiceGeneration|TestErrorHandling|TestPerformance"

# Run all tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out -covermode=atomic
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report saved to coverage.html"

# Run specific test
test-specific:
	@echo "Running specific test: $(TEST)"
	go test -v -run "$(TEST)"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	go test -v -race

# Run tests with verbose output
test-verbose:
	@echo "Running tests with verbose output..."
	go test -v -count=1

# Run benchmark tests
benchmark:
	@echo "Running benchmark tests..."
	go test -v -bench=.

# Clean test artifacts
test-clean:
	@echo "Cleaning test artifacts..."
	rm -rf testoutput
	rm -f protoc-gen-grpc-ts
	rm -f coverage.out coverage.html

# Clean everything
clean: test-clean
	@echo "Cleaning all artifacts..."
	@go clean
	@rm -f protoc-gen-grpc-ts
	@echo "Clean complete"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	go mod tidy
	go mod download

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	golangci-lint run

# Run the test script
run-tests: build
	@echo "Running comprehensive test suite..."
	./run_tests.sh

# Quick test (just build and run basic tests)
quick-test: build
	@echo "Running quick test..."
	go test -v -run "TestFieldNaming|TestNamespaceGeneration" -timeout 30s

# Full test suite (all tests with coverage and benchmarks)
full-test: deps fmt lint test-coverage benchmark

# Help
help:
	@echo "protoc-gen-grpc-ts - Makefile"
	@echo ""
	@echo "Build targets:"
	@echo "  build          - Build the plugin (default)"
	@echo "  install        - Build and install to PATH (use INSTALL_DIR=~/bin for custom location)"
	@echo ""
	@echo "Development:"
	@echo "  deps           - Install and update dependencies"
	@echo "  fmt            - Format Go code"
	@echo "  lint           - Lint Go code (requires golangci-lint)"
	@echo ""
	@echo "Testing:"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-specific  - Run specific test (use TEST=TestName)"
	@echo "  test-race      - Run tests with race detection"
	@echo "  test-verbose   - Run tests with verbose output"
	@echo "  benchmark      - Run benchmark tests"
	@echo "  quick-test     - Run quick test"
	@echo "  full-test      - Run full test suite (deps, fmt, lint, coverage, benchmark)"
	@echo ""
	@echo "Cleanup:"
	@echo "  test-clean     - Clean test artifacts"
	@echo "  clean          - Clean all artifacts (build files, test outputs)"
	@echo ""
	@echo "Examples:"
	@echo "  make              # Build the plugin"
	@echo "  make install      # Build and install to /usr/local/bin"
	@echo "  make test         # Run all tests"
	@echo "  make clean        # Clean all build artifacts"
