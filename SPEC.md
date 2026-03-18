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
    main.go                 # Entry point - only calls initialization
internal/
  config/
    config.go               # Configuration
  db/
    db.go                    # Database connection and migrations
  server/
    server.go                # HTTP server with graceful shutdown
  movies/
    entity.go               # Movie entity
    service.go              # Business logic
    repository.go           # Repository interface
    postgres_adapter.go     # PostgreSQL adapter
    router.go               # HTTP handlers
  books/
    entity.go               # Book entity
    service.go              # Business logic
    repository.go           # Repository interface
    postgres_adapter.go     # PostgreSQL adapter
    router.go               # HTTP handlers
frontend/
  # React + Vite application
migrations/
  # SQL migration files
```

## 3. API Endpoints

### Movies
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | /api/movies      | List all movies            |
| GET    | /api/movies/:id  | Get movie by ID            |
| POST   | /api/movies      | Create a new movie        |
| PUT    | /api/movies/:id  | Update an existing movie  |
| DELETE | /api/movies/:id  | Delete a movie            |

### Books
| Method | Endpoint         | Description                |
|--------|------------------|----------------------------|
| GET    | /api/books       | List all books            |
| GET    | /api/books/:id   | Get book by ID            |
| POST   | /api/books       | Create a new book        |
| PUT    | /api/books/:id   | Update an existing book  |
| DELETE | /api/books/:id   | Delete a book            |

## 4. Technology Stack

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

## 5. Commands

```bash
# Run API
go run cmd/api/main.go

# Build
go build -o bin/api.exe cmd/api/main.go

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs

# Docker
docker-compose up --build
```
