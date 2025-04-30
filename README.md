# Vamos

A Go project template with support for AWS, LLM, and web services.

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
│   │   └── client.go # AWS client code
│   ├── llm/          # LLM client implementation
│   │   └── client.go # LLM client code
│   └── web/          # Web server implementation
│       └── server.go # Web server code
├── pkg/               # Public library code
│   └── config/       # Configuration management
│       └── config.go # Configuration code
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

## Running the Applications

### Running Midas

1. Start the Midas application:
   ```bash
   go run cmd/midas/main.go
   ```

2. The application will:
   - Start a web server on port 8080
   - List S3 buckets
   - Generate text using the LLM
   - Wait for an interrupt signal (Ctrl+C) to shut down

### Running SF

1. Start the SF application:
   ```bash
   go run cmd/sf/main.go
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
go build -o bin/midas cmd/midas/main.go

# Build SF
go build -o bin/sf cmd/sf/main.go
```

## Development

### Adding New Dependencies

```bash
go get github.com/example/new-package
```

### Running Tests

```bash
go test ./...
```

### Code Formatting

```bash
go fmt ./...
```

### Linting

```bash
go vet ./...
```

## API Endpoints

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
