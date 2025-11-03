# Makefile for Media CMS Backend

.PHONY: help build run dev test test-verbose clean lint

# Binary output
BINARY_NAME=server.bin
MAIN_PATH=./cmd/server

# Go commands
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
GOMOD=$(GOCMD) mod

# Default target
.DEFAULT_GOAL := help

## help: Show this help message
help:
	@echo "Available targets:"
	@echo ""
	@echo "  make build         - Compile the server binary (output: $(BINARY_NAME))"
	@echo "  make run           - Run the server directly without building"
	@echo "  make dev           - Build and run the server"
	@echo "  make test          - Run all tests"
	@echo "  make test-verbose  - Run all tests with verbose output"
	@echo "  make test-unit     - Run only unit tests (short mode)"
	@echo "  make test-config   - Run config tests"
	@echo "  make clean         - Remove all .bin files and clean build cache"
	@echo "  make deps          - Download and tidy dependencies"
	@echo "  make lint          - Run go fmt and go vet"
	@echo "  make help          - Show this help message"
	@echo ""

## build: Compile the server binary
build:
	@echo "Building $(BINARY_NAME)..."
	$(GOBUILD) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BINARY_NAME)"

## run: Run the server directly
run:
	@echo "Running server..."
	$(GORUN) $(MAIN_PATH)

## dev: Build and run the server
dev: build
	@echo "Starting $(BINARY_NAME)..."
	./$(BINARY_NAME)

## test: Run all tests
test:
	@echo "Running tests..."
	$(GOTEST) ./...

## test-verbose: Run all tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	$(GOTEST) -v ./...

## test-unit: Run only unit tests (skip integration tests)
test-unit:
	@echo "Running unit tests..."
	$(GOTEST) -short ./...

## test-config: Run config package tests
test-config:
	@echo "Running config tests..."
	$(GOTEST) -v ./configs

## clean: Remove binaries and build cache
clean:
	@echo "Cleaning..."
	@rm -f *.bin
	@rm -f ./cmd/server/*.bin
	$(GOCLEAN) -cache
	@echo "Clean complete"

## deps: Download and tidy dependencies
deps:
	@echo "Downloading dependencies..."
	$(GOMOD) download
	$(GOMOD) tidy
	@echo "Dependencies updated"

## lint: Run go fmt and go vet
lint:
	@echo "Running linters..."
	@echo "  - go fmt"
	@gofmt -w .
	@echo "  - go vet"
	@$(GOCMD) vet ./...
	@echo "Linting complete"

## production: Build for production (optimized)
production:
	@echo "Building for production..."
	CGO_ENABLED=0 $(GOBUILD) -ldflags="-w -s" -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "Production build complete: $(BINARY_NAME)"
