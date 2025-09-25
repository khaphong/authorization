# Go Authorization Service

A comprehensive backend authentication service built with Go 1.25, featuring JWT-based authentication, PostgreSQL database, and Docker deployment.

## 📋 Table of Contents

- [Architecture](#architecture)
- [Project Structure](#project-structure)
- [Technology Stack](#technology-stack)
- [UUID v7 vs Autoincrement](#uuid-v7-vs-autoincrement)
- [Soft Delete Implementation](#soft-delete-implementation)
- [Getting Started](#getting-started)
- [API Endpoints](#api-endpoints)
- [Configuration](#configuration)
- [Development](#development)
- [Testing](#testing)
- [Deployment](#deployment)

## 🏗️ Architecture

This service follows **Clean Architecture** principles with clear separation of concerns:

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Layer    │    │  Business Layer │    │   Data Layer    │
│                 │    │                 │    │                 │
│ • Handlers      │───▶│ • Services      │───▶│ • Repositories  │
│ • Middleware    │    │ • DTOs          │    │ • Database      │
│ • Router        │    │ • Validation    │    │ • Models        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

**Data Flow:**
HTTP Request → Middleware → Handler → Service → Repository → Database

## 📁 Project Structure

```
authorization/
├── cmd/
│   └── server/main.go              # Application entry point
│
├── internal/
│   ├── config/config.go            # Configuration management
│   ├── server/router.go            # HTTP routing setup
│   │
│   ├── model/                      # Database models (one per file)
│   │   ├── base.go                 # Common base model with UUID, timestamps
│   │   ├── user.go                 # User model with soft delete
│   │   └── refresh_token.go        # Refresh token model
│   │
│   ├── handler/                    # HTTP request handlers
│   │   ├── auth_handler.go         # Authentication endpoints
│   │   └── user_handler.go         # User management endpoints
│   │
│   ├── service/                    # Business logic layer
│   │   ├── auth_service.go         # Authentication business logic
│   │   └── user_service.go         # User management business logic
│   │
│   ├── store/                      # Data access layer
│   │   ├── user_repo.go            # User repository with soft delete
│   │   └── token_repo.go           # Refresh token repository
│   │
│   ├── middleware/                 # HTTP middleware
│   │   ├── auth_middleware.go      # JWT authentication middleware
│   │   └── logging_middleware.go   # Request logging middleware
│   │
│   ├── dto/                        # Data Transfer Objects
│   │   ├── auth_request.go         # Authentication request DTOs
│   │   ├── auth_response.go        # Authentication response DTOs
│   │   └── error_response.go       # Error response DTOs
│   │
│   ├── utils/                      # Reusable utilities
│   │   ├── hash.go                 # Password hashing (argon2id)
│   │   ├── jwt.go                  # JWT token management
│   │   └── uuid.go                 # UUID v7 generation
│   │
│   ├── constants/                  # Application constants
│   │   ├── errors.go               # Error codes and messages
│   │   └── roles.go                # User roles (for future extension)
│   │
│   └── pkg/                        # Shared packages
│       ├── logger/logger.go        # Structured logging with Zap
│       └── response/response.go    # Centralized HTTP response helpers
│
├── migrations/                     # Database migrations
│   ├── 0001_create_users.up.sql
│   ├── 0001_create_users.down.sql
│   ├── 0002_create_refresh_tokens.up.sql
│   └── 0002_create_refresh_tokens.down.sql
│
├── docker/
│   └── Dockerfile                  # Multi-stage Docker build
│
├── docker-compose.yml              # Local development setup
├── .env.example                    # Environment variables template
├── Makefile                        # Build and development commands
└── README.md                       # This file
```

## 🛠️ Technology Stack

### Backend Framework & Libraries
- **HTTP Router**: [Chi](https://github.com/go-chi/chi) - Lightweight, fast router with middleware support
- **ORM**: [GORM](https://gorm.io/) - Developer-friendly ORM with soft delete support
- **JWT**: [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - Secure JWT implementation
- **Password Hashing**: [golang.org/x/crypto/argon2](https://pkg.go.dev/golang.org/x/crypto/argon2) - Industry-standard argon2id
- **Logging**: [Zap](https://github.com/uber-go/zap) - High-performance structured logging
- **Configuration**: [godotenv](https://github.com/joho/godotenv) - Environment variable management

### Database & Infrastructure
- **Database**: PostgreSQL 15 - Reliable, ACID-compliant relational database
- **Containerization**: Docker with multi-stage builds for minimal production images
- **Orchestration**: Docker Compose for local development

### Why Chi over Gin?
- **Lightweight**: Minimal dependencies and smaller binary size
- **Standard Library Compatible**: Built on net/http standards
- **Middleware Ecosystem**: Rich middleware ecosystem with clean composition
- **Performance**: Excellent performance with lower memory footprint

## 🆔 UUID v7 vs Autoincrement

### UUID v7 Advantages
- **Time-Ordered**: Natural sorting by creation time
- **Database Performance**: Better B-tree performance than random UUIDs
- **Distributed Systems**: No collision risk across multiple instances
- **Security**: Non-sequential, harder to guess than autoincrement
- **Future-Proof**: Scalable across microservices

### Implementation
```go
// UUID v7 generation with millisecond timestamp
func GenerateUUIDv7() string {
    timestamp := time.Now().UnixMilli()
    // ... timestamp + version + random data
}
```


## 🗑️ Soft Delete Implementation

### Two-Field Approach
```sql
deleted_at TIMESTAMP WITH TIME ZONE NULL,
is_deleted BOOLEAN DEFAULT FALSE NOT NULL
```

### Benefits
- **Performance**: Boolean index (`is_deleted`) faster than timestamp queries
- **Clarity**: Explicit boolean state for application logic
- **Audit Trail**: Timestamp preserves deletion time for compliance
- **Recovery**: Easy to restore soft-deleted records

### Database Indexes
```sql
-- Unique constraints only for active records
CREATE UNIQUE INDEX idx_users_username_active ON users(username) WHERE is_deleted = false;
CREATE UNIQUE INDEX idx_users_email_active ON users(email) WHERE is_deleted = false;

-- Performance indexes
CREATE INDEX idx_users_is_deleted ON users(is_deleted);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NOT NULL;
```

### GORM Integration
```go
type BaseModel struct {
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
    IsDeleted bool       `json:"is_deleted" gorm:"default:false;index"`
}

// Soft delete method
func (b *BaseModel) SoftDelete() {
    now := time.Now()
    b.DeletedAt = &now
    b.IsDeleted = true
}
```

## 🚀 Getting Started

### Prerequisites
- **Go 1.25+**
- **Docker & Docker Compose**
- **Make** (optional, for convenience commands)

### Quick Start
```bash
# Clone the repository
git clone <repository-url>
cd authorization

# Copy environment configuration
cp .env.example .env

# Start services with Docker Compose
make docker-up
# OR
docker-compose up -d

# Check service health
curl http://localhost:8080/health
```

### Manual Setup
```bash
# Install dependencies
make deps

# Build the application
make build

# Run PostgreSQL separately
docker run -d \
  --name postgres \
  -e POSTGRES_DB=go_login \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -p 8888:5432 \
  postgres:15-alpine

# Update .env with local database URL
echo "DATABASE_URL=postgres://postgres:postgres@localhost:8888/go_login?sslmode=disable" >> .env

# Run the application
make run
```

## 📡 API Endpoints

### Authentication Endpoints

#### Register User
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

**Response:**
```json
{
  "message": "User registered successfully",
  "user": {
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-09-25T10:30:00Z"
  }
}
```

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "f47ac10b58cc4372a567...",
  "expires_at": "2025-09-25T11:45:00Z",
  "user": {
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-09-25T10:30:00Z"
  }
}
```

#### Get Current User (Protected)
```bash
curl -X GET http://localhost:8080/api/v1/me \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIs..."
```

**Response:**
```json
{
  "id": "01234567-89ab-cdef-0123-456789abcdef",
  "username": "johndoe",
  "email": "john@example.com",
  "created_at": "2025-09-25T10:30:00Z"
}
```

#### Refresh Token
```bash
curl -X POST http://localhost:8080/api/v1/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "f47ac10b58cc4372a567..."
  }'
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "b58cc4372a567f47ac10...",
  "expires_at": "2025-09-25T12:00:00Z",
  "user": {
    "id": "01234567-89ab-cdef-0123-456789abcdef",
    "username": "johndoe",
    "email": "john@example.com",
    "created_at": "2025-09-25T10:30:00Z"
  }
}
```

#### Logout
```bash
curl -X POST http://localhost:8080/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "f47ac10b58cc4372a567..."
  }'
```

**Response:**
```json
{
  "message": "Logged out successfully"
}
```

### Health Check
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "healthy",
  "service": "authorization",
  "timestamp": 1696489200
}
```

## ⚙️ Configuration

### Environment Variables (.env)
```env
# Application
APP_ENV=development          # development, production
PORT=8080                   # HTTP server port

# Database
DB_HOST=postgres            # Database host
DB_PORT=5432               # Database port  
DB_USER=postgres           # Database username
DB_PASS=postgres           # Database password
DB_NAME=go_login           # Database name
DATABASE_URL=postgres://postgres:postgres@postgres:5432/go_login?sslmode=disable

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
ACCESS_TOKEN_EXP=15m       # Access token expiration (Go duration format)
REFRESH_TOKEN_EXP=168h     # Refresh token expiration (7 days)
```
## 💻 Development

### Available Make Commands
```bash
make help           # Show all available commands
make build          # Build the application
make run           # Build and run the application
make run-dev       # Run in development mode with hot reload
make test          # Run all tests
make test-coverage # Run tests with coverage report
make clean         # Clean build artifacts

# Docker commands
make docker-build  # Build Docker image
make docker-up     # Start services
make docker-down   # Stop services
make docker-logs   # View logs

# Code quality
make fmt           # Format code
make vet          # Run go vet
make lint         # Run linters (requires golangci-lint)

# Database
make db-connect   # Connect to PostgreSQL
make db-reset     # Reset database
```

### Development Workflow
```bash
# Full development cycle
make dev  # Runs: clean → fmt → vet → test → build
```

### Adding New Endpoints
1. **Define DTOs** in `internal/dto/`
2. **Add Business Logic** in `internal/service/`
3. **Create Handler** in `internal/handler/`
4. **Register Routes** in `internal/server/router.go`
5. **Add Tests** in `*_test.go` files

### Database Schema Changes
1. **Update Models** in `internal/model/`
2. **Create Migration Files** in `migrations/`
3. **Update Repository Methods** if needed
4. **Test with Clean Database**

## 🧪 Testing

### Run Tests
```bash
# All tests
make test

# With coverage
make test-coverage

# Specific package
go test ./internal/service/...

# Verbose output
go test -v ./...
```

### Test Structure
- **Unit Tests**: Service layer business logic
- **Integration Tests**: Database operations
- **Handler Tests**: HTTP endpoint testing

### Example Test Coverage
```
authorization/internal/service     coverage: 85.7% of statements
authorization/internal/utils       coverage: 92.3% of statements
```

## 🚀 Deployment

### Docker Production Build
```bash
# Build production image
make docker-build

# Or manually
docker build -f docker/Dockerfile -t authorization-app .
```

### Health Checks
The service provides health check endpoints for monitoring:
- **Kubernetes**: Use `/health` for readiness/liveness probes
- **Docker**: `HEALTHCHECK --interval=30s CMD curl -f http://localhost:8080/health`
- **Load Balancers**: Configure health checks on `/health`


## 📞 Support

For questions, issues, or contributions, please:

1. **Check Documentation** - Review this README and code comments
2. **Search Issues** - Look for existing GitHub issues
3. **Create Issue** - Submit detailed bug reports or feature requests
4. **Submit PR** - Contribute improvements with tests

---

**Built with ❤️ using Go 1.25, PostgreSQL, and Docker**
