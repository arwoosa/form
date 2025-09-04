# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Build and Run
- `make build` - Build the binary for Linux/amd64
- `make run` - Run the service locally with development config
- `go run ./cmd/server -c internal/conf/config.yaml` - Alternative run command

### Code Quality
- `make gotool` - Format code and run vet
- `make lint` - Run golangci-lint checks
- `make lint-fix` - Run golangci-lint with automatic fixes
- `go fmt ./...` - Format Go code
- `go vet ./...` - Run Go vet

### Protocol Buffers
- `make grpc` - Generate gRPC code from .proto files
- Generates Go files, gRPC-Gateway, and validation code

### Testing
- `make test` - Run all tests (includes unit and integration tests)
- `make test-unit` - Run unit tests for models and service layers (recommended)
- `make test-models` - Run only model layer tests
- `make test-service` - Run only service layer tests  
- `make test-integration` - Run integration tests with MongoDB testcontainers (currently has issues)
- `make test-coverage` - Generate HTML test coverage report
- `make test-race` - Run tests with race condition detection
- `make test-short` - Run quick tests only
- `go test ./...` - Alternative to run all tests directly

**Current Test Status:**
- ✅ **Unit Tests**: Models and Service layers have comprehensive coverage (100% passing)
- ❌ **Integration Tests**: MongoDB serialization issues need investigation
- 📊 **Coverage**: Core business logic fully tested, gRPC and config layers need work

### Module Management
- `make mod` - Update module name and references (used for project setup)
- `go mod tidy` - Clean up module dependencies

### Docker
- `make docker_run` - Run containerized service on port 8081

## Architecture Overview

This is a **Go gRPC microservice** for form management using the **Vulpes toolkit** framework. The service provides both gRPC and HTTP REST APIs through gRPC-Gateway.

### Key Architecture Patterns
- **Clean Architecture**: Separated into transport, service, repository, and data layers
- **Microservice**: Designed for horizontal scaling and service isolation
- **Dual Protocol**: gRPC for internal services, HTTP REST for web clients
- **Merchant Isolation**: Multi-tenant architecture with merchant-level data separation

### Core Technologies
- **Go 1.23+** with standard libraries
- **gRPC + gRPC-Gateway** for dual protocol support
- **MongoDB** for document storage with geospatial indexing
- **Vulpes Framework** for shared utilities and middleware
- **Prometheus** for metrics and monitoring
- **Viper** for configuration management

## Project Structure

```
form-service/
├── cmd/form-server/      # Application entry point
├── internal/             # Private application code
│   ├── conf/            # Configuration management
│   ├── service/         # Business logic layer
│   ├── dao/             # Data access layer
│   ├── models/          # Domain models
│   └── helper/          # Utility functions
├── proto/               # Protocol buffer definitions
├── pkg/vulpes/          # Shared toolkit (submodule)
└── deployments/         # Docker and deployment configs
```

### Vulpes Framework Integration
The project uses a local submodule of the Vulpes toolkit (`pkg/vulpes/`) which provides:
- **EzGRPC**: Automated gRPC server setup with interceptors
- **MongoDB utilities**: Connection management and operations
- **Structured logging**: Zap-based JSON logging
- **Prometheus metrics**: Automatic gRPC metrics collection
- **Validation**: Request/response validation utilities

## Configuration

### Config Files
- `internal/conf/config.yaml` - Main development configuration
- `internal/conf/config_docker.yaml` - Docker environment config
- Configuration is loaded via Viper with command-line flag support

### Environment Setup
The service requires:
- **MongoDB** running on localhost:27017 (configurable)
- **Go 1.23+** for development
- **Protocol Buffer compiler** for gRPC code generation

### Key Configuration Sections
- **MongoDB**: Database connection and credentials
- **Logging**: Structured JSON logging with rotation
- **Server**: Port configuration (default 8081)
- **Timezone**: Application timezone setting

## Business Domain

This service manages **Forms and Form Templates** for a multi-merchant platform with:

### Core Entities
- **FormTemplate**: Reusable form definitions with JSON Schema and UI Schema
- **Form**: Individual form instances based on templates or custom schemas
- **Merchant**: Tenant isolation boundary

### API Endpoints
- **Form Template API** (`/form_templates/*`): Full CRUD for form template management
- **Form API** (`/forms/*`): Full CRUD for form management

### Business Rules
- **Template Limit**: Configurable maximum number of templates per merchant
- **Schema Validation**: All forms must have valid JSON Schema
- **Template Duplication**: Templates can be duplicated within the same merchant

## Development Patterns

### Error Handling
- Use custom error types in `internal/dto/error.go`
- Convert business errors to appropriate gRPC status codes
- Structured error logging with context

### Database Operations
- Use repository pattern in `dao/repository/`
- MongoDB operations with proper indexing
- **Session Updates**: Smart diff-based bulk operations using MongoDB BulkWrite API
- Geospatial queries for location-based features

### Validation
- Request validation using Vulpes validate package
- Business rule validation in service layer
- Input sanitization for security

### Headers and Authentication
The service expects these headers from API Gateway:
- `X-User-Id`, `X-User-Email`, `X-User-Name`, `X-User-Avatar`
- `X-Merchant-Id` for merchant isolation (needs to be added to header map)

## Important Notes

- **No direct authentication**: Relies on API Gateway for auth/authz
- **Merchant isolation**: All data operations must include merchant_id filtering
- **Submodule dependency**: Vulpes toolkit is included as Git submodule
- **Proto generation**: Run `make grpc` after modifying .proto files