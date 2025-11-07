# Makefile for Mark2PDF

# Variables
BINARY_NAME=mark2pdf
BINARY_PATH=./bin/$(BINARY_NAME)
CLI_PATH=./cmd/mark2pdf
EXAMPLES_PATH=./examples
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# Colors for output
CYAN=\033[0;36m
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: all build test clean install run-example fmt vet lint help

# Default target
all: clean fmt vet test build

# Build the CLI binary
build:
	@echo "$(CYAN)Building $(BINARY_NAME)...$(NC)"
	@mkdir -p bin
	@$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BINARY_PATH) $(CLI_PATH)
	@echo "$(GREEN)✓ Build complete: $(BINARY_PATH)$(NC)"

# Build for multiple platforms
build-all:
	@echo "$(CYAN)Building for multiple platforms...$(NC)"
	@mkdir -p bin
	GOOS=linux GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(CLI_PATH)
	GOOS=linux GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-arm64 $(CLI_PATH)
	GOOS=darwin GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(CLI_PATH)
	GOOS=darwin GOARCH=arm64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-arm64 $(CLI_PATH)
	GOOS=windows GOARCH=amd64 $(GO) build $(LDFLAGS) -o bin/$(BINARY_NAME)-windows-amd64.exe $(CLI_PATH)
	@echo "$(GREEN)✓ Cross-platform builds complete$(NC)"

# Run tests
test:
	@echo "$(CYAN)Running tests...$(NC)"
	@$(GO) test -v -race -coverprofile=coverage.out ./...
	@echo "$(GREEN)✓ Tests complete$(NC)"

# Run tests with coverage report
test-coverage: test
	@echo "$(CYAN)Generating coverage report...$(NC)"
	@$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)✓ Coverage report generated: coverage.html$(NC)"

# Run benchmarks
bench:
	@echo "$(CYAN)Running benchmarks...$(NC)"
	@$(GO) test -bench=. -benchmem ./...

# Format code
fmt:
	@echo "$(CYAN)Formatting code...$(NC)"
	@$(GO) fmt ./...
	@echo "$(GREEN)✓ Code formatted$(NC)"

# Run go vet
vet:
	@echo "$(CYAN)Running go vet...$(NC)"
	@$(GO) vet ./...
	@echo "$(GREEN)✓ Vet complete$(NC)"

# Run static analysis
lint:
	@echo "$(CYAN)Running golangci-lint...$(NC)"
	@which golangci-lint > /dev/null || (echo "$(RED)golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest$(NC)" && exit 1)
	@golangci-lint run ./...
	@echo "$(GREEN)✓ Lint complete$(NC)"

# Install the CLI tool
install:
	@echo "$(CYAN)Installing $(BINARY_NAME)...$(NC)"
	@$(GO) install $(CLI_PATH)
	@echo "$(GREEN)✓ Installed to $(shell go env GOPATH)/bin/$(BINARY_NAME)$(NC)"

# Run the example
run-example:
	@echo "$(CYAN)Running example...$(NC)"
	@cd $(EXAMPLES_PATH) && $(GO) run basic.go
	@echo "$(GREEN)✓ Example complete$(NC)"

# Clean build artifacts
clean:
	@echo "$(CYAN)Cleaning...$(NC)"
	@rm -rf bin/
	@rm -f coverage.out coverage.html
	@rm -f $(EXAMPLES_PATH)/*.pdf
	@rm -f *.pdf
	@echo "$(GREEN)✓ Clean complete$(NC)"

# Download dependencies
deps:
	@echo "$(CYAN)Downloading dependencies...$(NC)"
	@$(GO) mod download
	@$(GO) mod tidy
	@echo "$(GREEN)✓ Dependencies updated$(NC)"

# Verify dependencies
verify:
	@echo "$(CYAN)Verifying dependencies...$(NC)"
	@$(GO) mod verify
	@echo "$(GREEN)✓ Dependencies verified$(NC)"

# Run the CLI tool with example
demo: build
	@echo "$(CYAN)Running demo...$(NC)"
	@echo "# Demo Document\n\nThis is a **demo** of Mark2PDF.\n\n## Features\n\n- Easy to use\n- Pure Go\n- No dependencies" > demo.md
	@$(BINARY_PATH) -input demo.md -output demo.pdf
	@rm -f demo.md
	@echo "$(GREEN)✓ Demo PDF created: demo.pdf$(NC)"

# Development: watch for changes and rebuild (requires entr)
watch:
	@which entr > /dev/null || (echo "$(RED)entr not installed. Install with: apt-get install entr (Linux) or brew install entr (macOS)$(NC)" && exit 1)
	@echo "$(CYAN)Watching for changes...$(NC)"
	@find . -name '*.go' | entr -c make build

# Check code quality
check: fmt vet lint test
	@echo "$(GREEN)✓ All checks passed!$(NC)"

# Show project statistics
stats:
	@echo "$(CYAN)Project Statistics:$(NC)"
	@echo "Lines of code:"
	@find . -name '*.go' -not -path './vendor/*' | xargs wc -l | tail -1
	@echo "\nFiles:"
	@find . -name '*.go' -not -path './vendor/*' | wc -l
	@echo "\nPackages:"
	@go list ./... | wc -l

# Help target
help:
	@echo "$(CYAN)Mark2PDF Makefile Commands:$(NC)"
	@echo ""
	@echo "$(YELLOW)Build Commands:$(NC)"
	@echo "  make build          - Build the CLI binary"
	@echo "  make build-all      - Build for multiple platforms"
	@echo "  make install        - Install the CLI tool to GOPATH"
	@echo ""
	@echo "$(YELLOW)Test Commands:$(NC)"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make bench          - Run benchmarks"
	@echo ""
	@echo "$(YELLOW)Code Quality:$(NC)"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint"
	@echo "  make check          - Run all quality checks"
	@echo ""
	@echo "$(YELLOW)Dependencies:$(NC)"
	@echo "  make deps           - Download and tidy dependencies"
	@echo "  make verify         - Verify dependencies"
	@echo ""
	@echo "$(YELLOW)Utilities:$(NC)"
	@echo "  make run-example    - Run the example program"
	@echo "  make demo           - Build and run a quick demo"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make watch          - Watch for changes and rebuild"
	@echo "  make stats          - Show project statistics"
	@echo ""
	@echo "$(YELLOW)Combined:$(NC)"
	@echo "  make all            - Clean, format, vet, test, and build"
	@echo "  make help           - Show this help message"
