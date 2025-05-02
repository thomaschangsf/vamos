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
