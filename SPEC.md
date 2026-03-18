# Movie Database API - Specification

## 1. Project Overview

- **Project Name**: Movie Database API
- **Project Type**: Full-stack web application (API + Frontend)
- **Core Functionality**: Movie and Book collection management system
- **Target Users**: General audience (no authentication required)

## 2. Architecture

### Domain-Based Structure

```
cmd/
  api/
    main.go                    # Entry point - only calls initialization
internal/
  domain/
    repository.go              # Generic Repository[T] interface
  config/
    config.go                 # Configuration
  db/
    db.go                     # Database connection and migrations
    migrate_test.go            # Migration tests
  server/
    server.go                 # HTTP server with graceful shutdown
  movies/
    entity.go                 # Movie entity with validation
    service.go                 # Movie service (uses domain.Repository[Movie])
    postgres_adapter.go        # PostgreSQL implementation
    router.go                  # HTTP handlers with Swagger annotations
    *_test.go                 # Unit and integration tests
  books/
    entity.go                 # Book entity with validation
    service.go                 # Book service (uses domain.Repository[Book])
    postgres_adapter.go        # PostgreSQL implementation
    router.go                  # HTTP handlers with Swagger annotations
    *_test.go                 # Unit and integration tests
frontend/                     # React + Vite application
migrations/
  *_up.sql                   # Up migrations
  *_down.sql                 # Down migrations
docs/                         # Swagger docs
```

## 3. API Endpoints

### Movies
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | /api/movies      | List all movies            |
| GET    | /api/movies/:id  | Get movie by ID            |
| POST   | /api/movies      | Create a new movie        |
| PUT    | /api/movies/:id  | Update an existing movie   |
| DELETE | /api/movies/:id  | Delete a movie            |

### Books
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | /api/books       | List all books            |
| GET    | /api/books/:id   | Get book by ID            |
| POST   | /api/books       | Create a new book         |
| PUT    | /api/books/:id   | Update an existing book   |
| DELETE | /api/books/:id   | Delete a book            |

## 4. Domain Interface

Generic repository used by both movies and books:

```go
type Repository[T any] interface {
    Create(ctx context.Context, entity *T) error
    GetByID(ctx context.Context, id int64) (*T, error)
    GetAll(ctx context.Context) ([]T, error)
    Update(ctx context.Context, entity *T) error
    Delete(ctx context.Context, id int64) error
}
```

## 5. Technology Stack

### Backend
- **Language**: Go
- **Framework**: Gin
- **Database**: PostgreSQL
- **Migration**: golang-migrate

### Frontend
- **Framework**: React + Vite
- **HTTP Client**: Axios

### Infrastructure
- **Docker**: Multi-stage builds
- **Database**: PostgreSQL

## 6. Testing

### Test Types
- **Unit Tests**: Service and entity validation tests using mocks
- **Integration Tests**: HTTP handlers with testcontainers PostgreSQL
- **Migration Tests**: Up and down migration tests

### Commands
```bash
# Unit tests only
go test ./... -short

# All tests (requires Docker)
go test ./...

# Specific package
go test ./internal/movies/... -v
go test ./internal/books/... -v
go test ./internal/db/... -v
```

## 7. Commands

```bash
# Run API
go run cmd/api/main.go

# Build
go build -o bin/api.exe ./cmd/api

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs

# Docker
docker-compose up --build
docker-compose up -d

# Logs
docker-compose logs -f

# Stop
docker-compose down
```

## 8. Configuration

Create a `.env` file (see `.env.example`):

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=movies
SERVER_PORT=8080
```

Environment variables take precedence over .env file values.
