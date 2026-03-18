# Movie Database API

A domain-based REST API for managing movie and book collections with a React frontend.

## Features

- **Movies**: CRUD operations with title, year, type (movie/series), resolution (SD/HD/FHD/4K), actors, rating
- **Books**: CRUD operations with title, author, release year, rating
- **Business Rules**:
  - Adult content not allowed for movies
  - Year must be valid (1888-current for movies, 1000-current for books)
  - Rating must be between 0 and 10

## Architecture

Domain-based structure (simplified hexagonal architecture):

- Each domain (movies, books) is self-contained
- Contains entity, service, adapter, and router in one package
- Shared interfaces in `domain/` package
- No strict layer separation - simpler but functional

```
internal/
  domain/              # Shared interfaces (Repository[T])
  movies/              # Movie domain (entity, service, adapter, router, tests)
  books/               # Book domain (entity, service, adapter, router, tests)
  config/              # Configuration
  db/                  # Database connection and migrations
  server/              # HTTP server with graceful shutdown
```

## Tech Stack

- **Backend**: Go, Gin, PostgreSQL, golang-migrate
- **Frontend**: React, Vite, Axios
- **Testing**: Go testing, testcontainers-go, testify

## Quick Start

### Docker Compose

```bash
docker-compose up --build
```

Services:
- Frontend: http://localhost
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- PostgreSQL: localhost:5432

## Configuration

Create a `.env` file based on `.env.example`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=movies
SERVER_PORT=8080
```

Environment variables take precedence over .env file values.

## API Endpoints

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

### Example Request

```bash
curl -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "year": 1999,
    "type": "movie",
    "resolution": "FHD",
    "actors": "Keanu Reeves",
    "rating": 8.7
  }'
```

## Testing

```bash
# Unit tests (fast)
go test ./... -short

# All tests including integration (requires Docker)
go test ./... -v
```

## Commands

```bash
# Run API
go run cmd/api/main.go

# Build
go build -o bin/api.exe ./cmd/api

# Run migrations
migrate -path migrations -database "postgres://..." up

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs

# Docker
docker-compose up --build
```
