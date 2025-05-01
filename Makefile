   # Build and run Midas
   # make build-midas
   # make run-midas

   # Or build and run SF
   # make build-sf
   # make run-sf

   # Check code quality
   # make test
   # make lint
   # make security

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME_MIDAS=bin/midas
BINARY_NAME_SF=bin/sf
BINARY_NAME_AWS=bin/aws
MAIN_MIDAS=cmd/midas/main.go
MAIN_SF=cmd/sf/main.go
MAIN_AWS=cmd/aws/main.go

# Build flags
LDFLAGS=-ldflags "-w -s"

# Test coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all build clean test run-midas run-sf run-aws deps tidy vet fmt lint help setup coverage test-aws test-llm test-web

all: clean deps build

help:
	@echo "Available targets:"
	@echo "  all         - Clean, install deps, and build both applications"
	@echo "  build       - Build both applications"
	@echo "  clean       - Remove build artifacts"
	@echo "  deps        - Install dependencies"
	@echo "  tidy        - Clean up dependencies"
	@echo "  test        - Run all tests"
	@echo "  test-aws    - Run AWS tests"
	@echo "  test-llm    - Run LLM tests"
	@echo "  test-web    - Run Web tests"
	@echo "  coverage    - Run tests with coverage"
	@echo "  vet         - Run go vet"
	@echo "  fmt         - Run go fmt"
	@echo "  lint        - Run golangci-lint"
	@echo "  run-midas   - Run Midas application"
	@echo "  run-sf      - Run SF application"
	@echo "  run-aws     - Run AWS client"
	@echo "  setup       - Setup development environment"
	@echo "  security    - Run security checks"

build: build-midas build-sf build-aws

build-midas:
	@echo "Building Midas..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_MIDAS) $(MAIN_MIDAS)

build-sf:
	@echo "Building SF..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_SF) $(MAIN_SF)

build-aws:
	@echo "Building AWS client..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_AWS) $(MAIN_AWS)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME_MIDAS)
	rm -f $(BINARY_NAME_SF)
	rm -f $(BINARY_NAME_AWS)
	rm -f $(COVERAGE_FILE)
	rm -f $(COVERAGE_HTML)

deps:
	@echo "Installing dependencies..."
	$(GOGET) -v ./...

tidy:
	@echo "Tidying dependencies..."
	$(GOMOD) tidy

test: test-aws test-llm test-web

test-aws:
	@echo "Running AWS tests..."
	$(GOTEST) -v ./internal/aws/...

test-llm:
	@echo "Running LLM tests..."
	$(GOTEST) -v ./internal/llm/...

test-web:
	@echo "Running Web tests..."
	$(GOTEST) -v ./internal/web/...

coverage:
	@echo "Running tests with coverage..."
	$(GOTEST) -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated at $(COVERAGE_HTML)"

vet:
	@echo "Running go vet..."
	$(GOCMD) vet ./...

fmt:
	@echo "Running go fmt..."
	$(GOCMD) fmt ./...

lint:
	@echo "Running golangci-lint..."
	golangci-lint run

run-midas:
	@echo "Running Midas..."
	$(GOCMD) run $(MAIN_MIDAS)

run-sf:
	@echo "Running SF..."
	$(GOCMD) run $(MAIN_SF)

run-aws:
	@echo "Running AWS client..."
	$(GOCMD) run $(MAIN_AWS)

# Development environment setup
setup:
	@echo "Setting up development environment..."
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint
	$(GOGET) -u github.com/securego/gosec/v2/cmd/gosec
	$(GOGET) -u golang.org/x/tools/cmd/goimports
	$(GOGET) -u github.com/stretchr/testify

# Security checks
security:
	@echo "Running security checks..."
	gosec ./...

# Generate documentation
docs:
	@echo "Generating documentation..."
	$(GOCMD) doc -all ./...

# Clean up all generated files
distclean: clean
	rm -f coverage.out coverage.html
	rm -rf bin/ 