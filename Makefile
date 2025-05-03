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

# Git workflow
# make story-start STORY_ID=123 DESCRIPTION="Feature description"
# make story-commit SCOPE=feature DESCRIPTION="Commit description"
# make story-push
# make undo HARD=false
# make revert COMMIT=abc123
# make tag VERSION=v1.0.3 MESSAGE="Stable snapshot" PUSH=true
# make sync MAIN=false
# make resolve REBASE=true

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
BINARY_NAME_WEB=bin/web
BINARY_NAME_GIT=bin/gitworkflow
MAIN_MIDAS=cmd/midas/main.go
MAIN_SF=cmd/sf/main.go
MAIN_AWS=cmd/aws/main.go
MAIN_WEB=main.go
MAIN_GIT=cmd/gitworkflow/main.go

# Build flags
LDFLAGS=-ldflags "-w -s"

# Test coverage
COVERAGE_FILE=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: all build clean test run-midas run-sf run-aws run-web deps tidy vet fmt lint help setup coverage test-aws test-llm test-web story-start story-commit story-push build-git undo revert tag sync resolve

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
	@echo "  run-web     - Run Web server"
	@echo "  setup       - Setup development environment"
	@echo "  security    - Run security checks"
	@echo "  build-git   - Build git workflow binary"
	@echo "  story-start - Start a new story branch (requires STORY_ID and DESCRIPTION)"
	@echo "  story-commit - Commit changes (requires SCOPE and DESCRIPTION)"
	@echo "  story-push  - Push current story branch"
	@echo "  undo        - Undo last commit (HARD=true to discard changes)"
	@echo "  revert      - Revert a specific commit (requires COMMIT)"
	@echo "  tag         - Create a version tag (requires VERSION and MESSAGE, PUSH=true to push)"
	@echo "  sync        - Sync with remote (MAIN=true to sync main branch)"
	@echo "  resolve     - Resolve conflicts (REBASE=false to use merge instead)"

build: build-midas build-sf build-aws build-web build-git

build-git:
	@echo "Building git workflow..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_GIT) $(MAIN_GIT)

# Git workflow commands
story-start:
	@if [ -z "$(STORY_ID)" ] || [ -z "$(DESCRIPTION)" ]; then \
		echo "Error: STORY_ID and DESCRIPTION are required"; \
		exit 1; \
	fi
	@echo "Starting new story branch..."
	$(BINARY_NAME_GIT) story-start --id $(STORY_ID) --description "$(DESCRIPTION)"

story-commit:
	@if [ -z "$(SCOPE)" ] || [ -z "$(DESCRIPTION)" ]; then \
		echo "Error: SCOPE and DESCRIPTION are required"; \
		exit 1; \
	fi
	@echo "Committing changes..."
	$(BINARY_NAME_GIT) story-commit --scope $(SCOPE) --description "$(DESCRIPTION)"

story-push:
	@echo "Pushing story branch..."
	$(BINARY_NAME_GIT) story-push

undo:
	@echo "Undoing last commit..."
	$(BINARY_NAME_GIT) undo --hard=$(HARD)

revert:
	@if [ -z "$(COMMIT)" ]; then \
		echo "Error: COMMIT is required"; \
		exit 1; \
	fi
	@echo "Reverting commit..."
	$(BINARY_NAME_GIT) revert --commit $(COMMIT)

tag:
	@if [ -z "$(VERSION)" ] || [ -z "$(MESSAGE)" ]; then \
		echo "Error: VERSION and MESSAGE are required"; \
		exit 1; \
	fi
	@echo "Creating tag..."
	$(BINARY_NAME_GIT) tag --version $(VERSION) --message "$(MESSAGE)" --push=$(PUSH)

sync:
	@echo "Syncing with remote..."
	$(BINARY_NAME_GIT) sync --main=$(MAIN)

resolve:
	@echo "Resolving conflicts..."
	$(BINARY_NAME_GIT) resolve --rebase=$(REBASE)

build-midas:
	@echo "Building Midas..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_MIDAS) $(MAIN_MIDAS)

build-sf:
	@echo "Building SF..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_SF) $(MAIN_SF)

build-aws:
	@echo "Building AWS client..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_AWS) $(MAIN_AWS)

build-web:
	@echo "Building Web server..."
	$(GOBUILD) $(LDFLAGS) -o $(BINARY_NAME_WEB) $(MAIN_WEB)

clean:
	@echo "Cleaning..."
	$(GOCLEAN)
	rm -f $(BINARY_NAME_MIDAS)
	rm -f $(BINARY_NAME_SF)
	rm -f $(BINARY_NAME_AWS)
	rm -f $(BINARY_NAME_WEB)
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

run-web:
	@echo "Running Web server..."
	$(GOCMD) run $(MAIN_WEB)

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