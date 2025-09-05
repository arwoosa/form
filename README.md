# Form Service

A Go-based form management microservice built with gRPC and gRPC-Gateway, providing APIs for form template and form operations.

## Architecture

This microservice follows Clean Architecture principles with:
- **gRPC + HTTP REST APIs**: Dual protocol support via gRPC-Gateway
- **MongoDB**: Document storage with geospatial indexing
- **Merchant Isolation**: Multi-tenant architecture with merchant-level data separation
- **Vulpes Framework**: Shared utilities and middleware

## Environment Requirements

### System Requirements
- **Go**: 1.23+ (toolchain go1.24.1)
- **MongoDB**: 4.4+ (for geospatial queries support)
- **Protocol Buffers**: Latest version for gRPC code generation

### Development Dependencies
```bash
# Install Protocol Buffer compiler
brew install protobuf

# Install Go tools for gRPC
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
```

## Configuration Parameters

### Main Configuration (`conf/config.yaml`)

#### Application Settings
```yaml
name: "partivo_form"          # Service name
mode: "dev"                    # Environment: "dev" or "prod"
port: 8081                     # Service port
version: "1.0.0"               # Application version
time_zone: "Asia/Taipei"       # Timezone for the application
```

#### Logging Configuration
```yaml
log:
  level: "debug"               # Log level: debug, info, warn, error
```

#### MongoDB Configuration
```yaml
mongodb:
  host: "127.0.0.1"           # MongoDB host
  port: 27017                  # MongoDB port
  # user: ""                   # MongoDB username (optional)
  # password: ""               # MongoDB password (optional)
  db: "partivo_form"         # Database name
```

#### External Services
```yaml
external:
  order_service:
    endpoint: "192.168.1.134:8081"  # Order service gRPC endpoint
    timeout: "10s"                  # Request timeout
  # media_service:                  # Future media service config
  #   endpoint: ""
  #   timeout: "10s"
```

#### Pagination Settings
```yaml
pagination:
  default_page_size: 20        # Default page size for listings
  max_page_size: 100           # Maximum allowed page size
  default_location_radius: 1000 # Default radius for geo queries (meters)
```

### Docker Configuration (`conf/config_docker.yaml`)

Same as main configuration but with Docker-specific settings:
```yaml
mongodb:
  host: "host.docker.internal"  # Docker host networking
```

## Deployment Methods

### 1. Local Development

#### Prerequisites
```bash
# Start MongoDB
brew services start mongodb-community
# or
docker run -d -p 27017:27017 --name mongodb mongo:latest
```

#### Build and Run
```bash
# Build the service
make build

# Run the service
make run

# Or run with custom config
go run ./cmd/form-server --config conf/config.yaml
```

### 2. Docker Deployment

#### Build Docker Image
```bash
# Build image
docker build -t partivo_form:1.0 .

# Run container
make docker_run
# or
docker run --rm partivo_form:1.0 server --config conf/config_docker.yaml 
```

#### Docker Compose Deployment (Recommended)

Use the provided Docker Compose setup to run both Console and Public APIs:

```bash
# Build image first
docker build -t partivo_form:1.0 .

# Start all services (Console API + Public API + MongoDB)
make docker-compose-up
# or
docker-compose up -d

# Check logs
make docker-compose-logs

# Stop all services
make docker-compose-down
```

**Service URLs:**
- Form Service API: http://localhost:8081
- MongoDB: localhost:27017

**Service commands:**
```bash
# Restart service
make docker-restart
```

### 3. Production Deployment

#### Environment Variables
```bash
# Override config file location
export CONFIG_FILE=/app/conf/config_production.yaml

# MongoDB connection (if using environment variables)
export MONGODB_HOST=your-mongodb-host
export MONGODB_PORT=27017
export MONGODB_DB=partivo_form_prod
```

#### Production Checklist
- [ ] Use `mode: "prod"` in configuration
- [ ] Configure proper MongoDB replica set
- [ ] Set up log rotation and monitoring
- [ ] Configure service monitoring and alerting
- [ ] Set up reverse proxy (nginx) for load balancing
- [ ] Enable HTTPS/TLS termination

## Service Startup

### Form Service API

The service provides a unified API for form management:

#### Form Templates API
- **Endpoints**: `/form_templates/*`
- **Features**: Full CRUD operations for form templates
- **Authentication**: Requires API Gateway headers

#### Forms API
- **Endpoints**: `/forms/*`
- **Features**: Full CRUD operations for forms
- **Authentication**: Requires API Gateway headers

### API Gateway Headers

The service expects these headers from the API Gateway:
```
X-User-Id: user-uuid
X-User-Email: user@example.com  
X-User-Name: User Name
X-Merchant-Id: merchant-uuid           # Required for merchant isolation
```

### Service Monitoring

Monitor services using Docker Compose:
```bash
# View service status
docker-compose ps

# View service logs
docker-compose logs -f form-service

# Check resource usage
docker stats
```

## Development Commands

### Code Generation
```bash
# Generate gRPC code from proto files
make grpc
```

### Testing
```bash
# Run all tests
make test

# Run unit tests only (recommended)
make test-unit

# Run with coverage
make test-coverage

# Run integration tests (requires Docker)
make test-integration
```

### Code Quality
```bash
# Format and vet code
make gotool

# Individual commands
go fmt ./...
go vet ./...
```

## Troubleshooting

### Common Issues

#### 1. MongoDB Connection Failed
**Error**: `failed to connect to mongodb`
**Solution**: 
```bash
# Check MongoDB is running
brew services list | grep mongodb
# or
docker ps | grep mongo

# Test connection
mongosh --host localhost --port 27017
```

#### 2. Port Already in Use
**Error**: `bind: address already in use`
**Solution**:
```bash
# Find process using port 8081
lsof -i :8081
kill -9 <PID>

# Or change port in config
port: 8082
```

#### 3. gRPC Code Generation Failed
**Error**: `protoc: command not found`
**Solution**:
```bash
# Install Protocol Buffer compiler
brew install protobuf

# Verify installation
protoc --version
```

#### 4. Vulpes Dependency Issues
**Error**: `vulpes package not found`
**Solution**:
```bash
# Check go.mod replace directive
grep vulpes go.mod

# Re-download dependencies
go mod download
go mod tidy
```

#### 5. Docker Build Failed
**Error**: `failed to compute cache key`
**Solution**:
```bash
# Clean Docker cache
docker system prune -a

# Check Dockerfile paths
docker build --no-cache -t partivo_event:1.0 .
```

### Logging and Debugging

#### Enable Debug Logging
```yaml
# In config.yaml
log:
  level: "debug"
mode: "dev"
```

#### Check Application Logs
```bash
# Console output (development)
tail -f logs/app.log

# Docker logs
docker logs -f container_name
```

### Dependencies Check

#### Verify External Services
```bash
# Test Order Service connection
grpcurl -plaintext 192.168.1.134:8081 list

# Test MongoDB connection  
mongosh mongodb://localhost:27017/partivo_form
```

### Performance Monitoring

#### MongoDB Indexes
```bash
# Check if indexes are created
mongosh partivo_form --eval "db.form_templates.getIndexes()"
mongosh partivo_form --eval "db.forms.getIndexes()"
```

#### Memory and CPU Usage
```bash
# Monitor resource usage
docker stats partivo_form
top -p $(pgrep form-server)
```

---

## Notes

- `dao/mongodb/migration` can define collections and indexes, which will be created during initialization
- The service uses the Vulpes framework for shared utilities and middleware
- Merchant isolation is implemented at the service layer - ensure all operations include merchant filtering
- For production deployment, consider using a process manager like systemd or supervisord