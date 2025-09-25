# Form Service

A Go-based form management microservice built with gRPC and gRPC-Gateway, providing APIs for **form template** operations. This service focuses specifically on managing reusable form templates with JSON Schema and UI Schema support.

## Architecture

This microservice follows Clean Architecture principles with:
- **gRPC + HTTP REST APIs**: Dual protocol support via gRPC-Gateway.
- **MongoDB**: Document storage for form templates.
- **Multi-Tenancy**: Data is isolated by `merchant_id`.
- **Vulpes Framework**: A local submodule providing shared utilities and middleware.

## Environment Requirements

### System Requirements
- **Go**: 1.24.0+ (toolchain go1.24.6)
- **MongoDB**: 4.4+
- **Docker**: For running `make grpc` and for containerized deployment.
- **Protocol Buffers**: `protoc` compiler. Can be installed via `brew install protobuf`.

### Go Tools (for Local `protoc` Execution)

**Note:** If you exclusively use the `make grpc` command for code generation, you can **skip** this section. The required tools are already included in the Docker image used by the Makefile.

The following tools are only necessary if you intend to run the `protoc` compiler directly on your host machine:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
```

## How to Run

### 1. Local Development

**Prerequisites:**
Ensure MongoDB is running and accessible.

**Build and Run:**
```bash
# Format code and run go vet
make gotool

# Build the binary into ./bin/form-server
make build

# Run the server using the default configuration
make run
# This is equivalent to:
# go run ./cmd/form-server server --config conf/config.yaml

# To run with a different configuration (e.g., for Docker networking)
make run-docker
# This is equivalent to:
# go run ./cmd/form-server server --config conf/config_docker.yaml
```

### 2. Docker Deployment (Recommended)

The provided `docker-compose.yml` file orchestrates the service and a MongoDB instance.

```bash
# Build the Docker image
docker build -t partivo_form:1.0 .

# Start all services in the background
make docker-compose-up

# View logs
make docker-compose-logs

# Stop all services
make docker-compose-down

# To restart the services
make docker-restart
```

**Service URLs:**
- **Form Service API**: `http://localhost:8081`
- **MongoDB**: `localhost:27017`

## Development Commands

The `Makefile` contains helpers for common development tasks.

### Code Generation
To regenerate gRPC, gRPC-Gateway, and validation code from the `.proto` files:
```bash
# This command runs protoc via a Docker container
make grpc
```

### Linting
```bash
# Run the linter
make lint

# Run the linter and automatically fix issues
make lint-fix
```

### Testing
The project has a comprehensive test suite.

```bash
# Run all tests
make test

# Run unit tests only (models and services)
make test-unit

# Run integration tests (repository layer, requires Docker)
make test-integration

# Generate and view a test coverage report
make test-coverage

# Clean up coverage files
make test-clean
```

## API Endpoints

The service exposes the following endpoints for form template management.

- `POST /form_templates`: Create a new form template.
- `GET /form_templates`: List all form templates for a merchant (supports pagination).
- `GET /form_templates/{id}`: Get a single form template by its ID.
- `PUT /form_templates/{id}`: Update an existing form template.
- `DELETE /form_templates/{id}`: Delete a form template.
- `POST /form_templates/{id}/duplicate`: Create a copy of an existing form template.
- `GET /config`: Retrieve frontend-relevant configuration, such as business rules.

**Note**: The `Form` entity APIs (for managing form instances) are defined in the `.proto` file but are currently commented out and not served by the application.

## Configuration

Configuration is managed via `conf/config.yaml` and can be overridden by environment variables.

### Key Parameters (`conf/config.yaml`)

```yaml
name: "partivo_form_service"
mode: "dev"                    # "dev" or "prod"
port: 8081
time_zone: "UTC"

log:
  level: "debug"               # debug, info, warn, error

mongodb:
  host: "127.0.0.1"
  port: 27017
  db: "partivo"

keto:
  write_addr: "172.20.0.22:4467"
  read_addr: "172.20.0.22:4466"

pagination:
  default_page_size: 1000
  max_page_size: 2000

business_rules:
  max_templates_per_merchant: 3
```

## Troubleshooting

### MongoDB Connection Failed
Ensure your MongoDB instance is running and accessible at the host/port specified in your active `config.yaml`. For local development, `127.0.0.1:27017` is the default. For Docker Compose, the service connects to the `mongodb` service name.

### Port Already in Use
If you see an `address already in use` error, another process is occupying port `8081`. Find and stop the process or change the `port` in your `config.yaml`.

### `make grpc` Fails
This command depends on Docker. Ensure the Docker daemon is running. It uses the `94peter/grpc-gateway-builder` image to execute `protoc` with all necessary plugins.
