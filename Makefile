.PHONY: build build-all test lint install release clean help

# Build variables
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE ?= $(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS := -X main.Version=$(VERSION) -X main.Commit=$(COMMIT) -X main.Date=$(DATE)

# Binary name
BINARY_NAME := qspin

# Build directory
BUILD_DIR := bin

# Go commands
GOCMD := go
GOBUILD := $(GOCMD) build
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt
GOVET := $(GOCMD) vet

help: ## Display this help message
	@echo "QuickSpin CLI - Makefile commands:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

build: ## Build the binary for current platform
	@echo "Building $(BINARY_NAME) $(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/qspin
	@echo "Binary built: $(BUILD_DIR)/$(BINARY_NAME)"

build-all: ## Build binaries for all platforms
	@echo "Building for all platforms..."
	@mkdir -p $(BUILD_DIR)
	# Linux amd64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/qspin
	# Linux arm64
	GOOS=linux GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/qspin
	# macOS amd64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/qspin
	# macOS arm64 (Apple Silicon)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/qspin
	# Windows amd64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/qspin
	@echo "All binaries built in $(BUILD_DIR)/"

test: ## Run tests
	@echo "Running tests..."
	$(GOTEST) -v -race -coverprofile=coverage.txt -covermode=atomic ./...

test-coverage: test ## Run tests with coverage report
	@echo "Generating coverage report..."
	$(GOCMD) tool cover -html=coverage.txt -o coverage.html
	@echo "Coverage report: coverage.html"

lint: ## Run linters
	@echo "Running linters..."
	@if command -v golangci-lint > /dev/null; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Running go vet..."; \
		$(GOVET) ./...; \
	fi

fmt: ## Format code
	@echo "Formatting code..."
	$(GOFMT) ./...

install: build ## Install binary to GOPATH/bin
	@echo "Installing $(BINARY_NAME)..."
	@mkdir -p $(GOPATH)/bin
	@cp $(BUILD_DIR)/$(BINARY_NAME) $(GOPATH)/bin/
	@echo "Installed to $(GOPATH)/bin/$(BINARY_NAME)"

release: ## Create release with goreleaser
	@echo "Creating release..."
	@if command -v goreleaser > /dev/null; then \
		goreleaser release --clean; \
	else \
		echo "goreleaser not installed. Install with: go install github.com/goreleaser/goreleaser@latest"; \
		exit 1; \
	fi

clean: ## Clean build artifacts
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.txt coverage.html
	@echo "Cleaned"

deps: ## Download dependencies
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy

run: ## Run the CLI (dev mode)
	@$(GOCMD) run -ldflags "$(LDFLAGS)" ./cmd/qspin

docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t quickspin/qspin:$(VERSION) .
	docker tag quickspin/qspin:$(VERSION) quickspin/qspin:latest
	@echo "Docker image built: quickspin/qspin:$(VERSION)"

docker-run: ## Run Docker container
	docker run -it --rm quickspin/qspin:latest

.DEFAULT_GOAL := help
