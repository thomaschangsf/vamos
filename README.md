# Vamos

A Go project template with support for AWS, LLM, and web services.

## Available Binaries

The project includes several command-line tools:

- `vamosMidas` - Midas application
- `vamosSF` - SF application
- `vamosAWS` - AWS client
- `vamosWeb` - Web server
- `vamosGitWF` - Git workflow management tool

## Project Structure

```
.
├── cmd/                # Application entry points
│   ├── midas/         # Midas application
│   │   └── main.go    # Midas entry point
│   └── sf/            # SF application
│       └── main.go    # SF entry point
├── internal/          # Private application and library code
│   ├── aws/          # AWS client implementation
│   │   ├── client.go # AWS client code
│   │   └── client_test.go # AWS tests
│   ├── llm/          # LLM client implementation
│   │   ├── client.go # LLM client code
│   │   └── client_test.go # LLM tests
│   └── web/          # Web server implementation
│       ├── server.go # Web server code
│       └── server_test.go # Web tests
├── pkg/               # Public library code
│   └── config/       # Configuration management
│       └── config.go # Configuration code
├── bin/              # Binary output directory
│   ├── midas        # Midas binary
│   ├── sf           # SF binary
│   ├── aws          # AWS binary
│   └── web          # Web server binary
├── go.mod            # Go module definition
├── go.sum            # Go module checksums
└── README.md         # Project documentation
```

## Features

- AWS Integration
  - S3 bucket operations
  - AWS SDK v2 support
- LLM Integration
  - OpenAI API support
  - Text generation capabilities
- Web Server
  - RESTful API endpoints
  - Health check endpoints
  - Graceful shutdown

## Prerequisites

- Go 1.22 or later
- AWS account and credentials
- OpenAI API key
- Git

## Setup

1. Clone the repository:
   ```bash
   git clone https://github.com/chang/vamos.git
   cd vamos
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Set up environment variables:
   ```bash
   # AWS Configuration
   export AWS_REGION=us-west-2
   export AWS_ACCESS_KEY=your_access_key
   export AWS_SECRET_KEY=your_secret_key

   # LLM Configuration
   export LLM_API_KEY=your_llm_api_key
   export LLM_MODEL_NAME=gpt-3.5-turbo

   # Web Configuration
   export WEB_PORT=8080
   export WEB_BASE_URL=http://localhost:8080
   ```

   Alternatively, you can create a `.env` file in the project root:
   ```bash
   AWS_REGION=us-west-2
   AWS_ACCESS_KEY=your_access_key
   AWS_SECRET_KEY=your_secret_key
   LLM_API_KEY=your_llm_api_key
   LLM_MODEL_NAME=gpt-3.5-turbo
   WEB_PORT=8080
   WEB_BASE_URL=http://localhost:8080
   ```

## Testing

The project includes comprehensive tests for all internal packages. Tests use mocking to avoid external service dependencies.

### Running Tests

1. Install test dependencies:
   ```bash
   make setup
   ```

2. Run all tests:
   ```bash
   make test
   ```

3. Run tests with coverage:
   ```bash
   make coverage
   ```

4. Run specific package tests:
   ```bash
   # AWS tests
   go test ./internal/aws

   # LLM tests
   go test ./internal/llm

   # Web tests
   go test ./internal/web
   ```

### Test Structure

- **AWS Tests**:
  - Tests S3 client functionality
  - Mocks AWS SDK responses
  - Tests bucket listing operations

- **LLM Tests**:
  - Tests OpenAI client functionality
  - Mocks API responses
  - Tests text generation

- **Web Tests**:
  - Tests HTTP endpoints
  - Tests server startup and shutdown
  - Tests health check endpoints

### Test Coverage

To generate and view test coverage:

```bash
# Generate coverage report
make coverage

# View coverage in browser
open coverage.html
```

## Running the Applications

### Running Web Server

1. Build the web server binary:
   ```bash
   make build-web
   ```

2. Start the web server:
   ```bash
   make run-web
   ```
   Or run the binary directly:
   ```bash
   ./bin/web
   ```

3. Test the weather endpoint:
   ```bash
   curl "http://localhost:8080/api/weather?city=San%20Francisco&state=CA" | python -m json.tool
   ```

   The response will look like:
   ```json
   {
     "city": "San Francisco",
     "state": "CA",
     "forecast": [
       {
         "date": "2024-03-21",
         "temperature": {
           "min": 15.0,
           "max": 25.0
         },
         "description": "Sunny",
         "humidity": 65,
         "wind_speed": 5.5
       }
     ]
   }
   ```

### Running Midas

1. Start the Midas application:
   ```bash
   make run-midas
   ```

2. The application will:
   - Start a web server on port 8080
   - List S3 buckets
   - Generate text using the LLM
   - Wait for an interrupt signal (Ctrl+C) to shut down

### Running SF

1. Start the SF application:
   ```bash
   make run-sf
   ```

2. The application will:
   - Start a web server on port 8080
   - List S3 buckets
   - Generate text using the LLM
   - Wait for an interrupt signal (Ctrl+C) to shut down

### Building Binaries

To build standalone binaries:

```bash
# Build Midas
make build-midas

# Build SF
make build-sf
```

## Development

### Adding New Dependencies

```bash
go get github.com/example/new-package
```

### Code Quality

```bash
# Format code
make fmt

# Run linter
make lint

# Run security checks
make security
```

### API Endpoints

Both applications expose the following endpoints:

- `GET /health` - Health check endpoint
- `GET /api/status` - Status endpoint with current time

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Git Workflows

The project includes a comprehensive git workflow management system to help maintain clean and consistent git practices.

### Story Workflow

The story workflow helps manage feature development with proper branch naming and commit conventions.

1. **Start a new story**:
```bash
# Create a new story branch
make story-start STORY_ID=456 DESCRIPTION="Chat UI"
# Creates branch: feat/story-456-chat-ui
```

2. **Make changes and commit**:
```bash
# Commit changes with proper formatting
make story-commit SCOPE=chat DESCRIPTION="add user message bubble"
# Creates commit: feat(chat): add user message bubble
```

3. **Push the branch**:
```bash
make story-push
```

### Version Management

Manage version tags and stable points in your codebase.

```bash
# Create and push a version tag
make tag VERSION=v1.0.3 MESSAGE="Stable snapshot before auth refactor" PUSH=true
```

### Safe Reverts

Safely undo changes or revert to previous states.

```bash
# Undo last commit (keep changes in working directory)
make undo HARD=false

# Undo last commit and discard changes
make undo HARD=true

# Revert a specific commit
make revert COMMIT=abc123
```

### Remote Synchronization

Keep your local repository in sync with remote and handle conflicts.

```bash
# Sync current branch with remote
# This will:
# 1. Check for local changes (will fail if there are uncommitted changes)
# 2. Check if your branch is ahead of remote (will fail if you need to push)
# 3. Pull and rebase if your branch is behind remote
make sync MAIN=false

# Sync main branch with remote
# This checks out the main branch and pulls the latest changes from origin/main
make sync MAIN=true

# Resolve conflicts by rebasing (default)
make resolve REBASE=true

# Resolve conflicts by merging
make resolve REBASE=false
```

### Complete Workflow Example

Here's a complete example of a typical development workflow:

```bash
# Start a new feature
make story-start STORY_ID=456 DESCRIPTION="Chat UI"

# Make changes and commit
make story-commit SCOPE=chat DESCRIPTION="add user message bubble"

# If conflicts arise
make sync MAIN=true
make resolve REBASE=true

# Push changes
make story-push

# Create a stable point
make tag VERSION=v1.0.3 MESSAGE="Stable snapshot" PUSH=true
```

### Using the Binary Directly

You can also use the git workflow binary directly:

```bash
# Start a story
vamosGitWF story-start --id 456 --description "Chat UI"

# Commit changes
vamosGitWF story-commit --scope chat --description "add user message bubble"

# Push branch
vamosGitWF story-push

# Sync with remote
vamosGitWF sync --main=false

# Resolve conflicts
vamosGitWF resolve --rebase=true

# Create and push a tag
vamosGitWF tag --version v1.0.3 --message "Stable release" --push=true
```

## Git Workflow Tool

The `vamosGitWF` tool helps manage git workflow with a standardized approach to branch naming and commit messages.

### Branch Naming Convention

Branches follow the format: `W-STORY_ID[-description]`
- `W-` prefix is required
- `STORY_ID` is required (e.g., "123")
- `description` is optional (e.g., "add-login")

Examples:
- `W-123` - Basic story branch
- `W-123-add-login` - Story branch with description

### Common Commands

1. Start a new story:
   ```bash
   # With description
   vamosGitWF story-start --id "123" --description "add-login"

   # Without description
   vamosGitWF story-start --id "123"
   ```

2. Commit changes:
   ```bash
   vamosGitWF story-commit --scope "auth" --description "implement login flow"
   ```

3. Sync with remote:
   ```bash
   vamosGitWF sync
   ```

4. Push changes:
   ```bash
   vamosGitWF story-push
   ```

5. Create a version tag:
   ```bash
   vamosGitWF tag --version "v1.0.0" --message "Initial release" --push
   ```

### Best Practices

- Always start from the main branch
- Commit frequently with clear, descriptive messages
- Push your changes regularly
- Use sync before starting new work
- Resolve conflicts as soon as they appear
- Tag releases when features are complete
- Use tags as stable points to revert to if needed

## Building and Installing

### Quick Start
```bash
# Build all binaries
make build

# Add binaries to your PATH (for current session)
export PATH="/Users/thomaschang/Documents/dev/git/thomaschangsf/vamos/bin:$PATH"

# To make this permanent, add this line to your ~/.bashrc or ~/.zshrc:
export PATH="/Users/thomaschang/Documents/dev/git/thomaschangsf/vamos/bin:$PATH"

# Now you can run any tool from anywhere
vamosGitWF story-start --id 123 --description "New feature"
vamosMidas
vamosSF
vamosAWS
vamosWeb
```

### Detailed Build Instructions

1. Build specific binaries:
```bash
make build-midas    # Build vamosMidas
make build-sf       # Build vamosSF
make build-aws      # Build vamosAWS
make build-web      # Build vamosWeb
make build-git      # Build vamosGitWF
```

2. Build all binaries at once:
```bash
make build
```

The binaries will be created in the `bin/` directory:
- `bin/vamosMidas`
- `bin/vamosSF`
- `bin/vamosAWS`
- `bin/vamosWeb`
- `bin/vamosGitWF`

### Git Workflow Tool (vamosGitWF)

The `vamosGitWF` tool provides several commands for managing git workflows:

```bash
# Start a new story
vamosGitWF story-start --id 123 --description "New feature"

# Commit changes
vamosGitWF story-commit --scope feature --description "Add login button"

# Push story branch
vamosGitWF story-push

# Undo last commit
vamosGitWF undo --hard=false

# Revert a specific commit
vamosGitWF revert --commit abc123

# Create and push a version tag
vamosGitWF tag --version v1.0.3 --message "Stable release" --push=true
```
  