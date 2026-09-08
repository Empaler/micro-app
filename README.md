# Movie Database API

A domain-based REST API for managing movie and book collections with a React frontend.

## Features

- **Movies**: CRUD operations with title, year, type (movie/series), resolution (SD/HD/FHD/4K), actors, rating
- **Books**: CRUD operations with title, author, release year, rating
- **Most looked up items**: Redis-backed popularity rankings for movies and books based on detail-page views
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
backend/
  cmd/
    api/                 # API entry point
  internal/
  domain/              # Shared interfaces (Repository[T])
  movies/              # Movie domain (entity, service, adapter, router, tests)
  books/               # Book domain (entity, service, adapter, router, tests)
  redisclient/         # Redis sorted-set popularity client
  config/              # Configuration
  db/                  # Database connection and migrations
  server/              # HTTP server with graceful shutdown
  migrations/           # Database migrations
  docs/                 # Generated Swagger docs
  Dockerfile            # Backend image definition
frontend/               # React + Vite application
```

## Tech Stack

- **Backend**: Go, Gin, PostgreSQL, Redis, golang-migrate
- **Frontend**: React, Vite, Axios
- **Testing**: Go testing, testcontainers-go, testify, miniredis

## Quick Start

### Docker Compose

```bash
docker-compose up --build
```

### Development Services Only

Start PostgreSQL and Redis in Docker, then run the backend and frontend manually or from the IDE:

```bash
docker compose -f docker-compose.dev.yml up -d
```

Run the backend from `backend/` with the `Launch API` VS Code configuration, or:

```bash
cd backend
go run cmd/api/main.go
```

Run the frontend separately:

```bash
cd frontend
npm install
npm run dev
```

Stop only the development database services with:

```bash
docker compose -f docker-compose.dev.yml down
```

Services:
- Frontend: http://localhost
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- PostgreSQL: localhost:5432
- Redis: localhost:6379

## Configuration

Create a `.env` file based on `.env.example`:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=movies
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DB=0
SERVER_PORT=8080
```

Environment variables take precedence over .env file values.

## API Endpoints

### Movies
| Method | Endpoint                    | Description                    |
|--------|-----------------------------|--------------------------------|
| GET    | /api/movies                 | List all movies                |
| GET    | /api/movies/:id             | Get movie by ID                |
| GET    | /api/movies/most-looked-up | Get most looked up movies     |
| POST   | /api/movies                 | Create a new movie             |
| PUT    | /api/movies/:id             | Update an existing movie       |
| DELETE | /api/movies/:id             | Delete a movie                 |

### Books
| Method | Endpoint                    | Description                    |
|--------|-----------------------------|--------------------------------|
| GET    | /api/books                  | List all books                 |
| GET    | /api/books/:id              | Get book by ID                 |
| GET    | /api/books/most-looked-up   | Get most looked up books      |
| POST   | /api/books                  | Create a new book              |
| PUT    | /api/books/:id              | Update an existing book        |
| DELETE | /api/books/:id              | Delete a book                  |

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
cd backend

# Unit tests (fast)
go test ./... -short

# All tests including integration (requires Docker)
go test ./... -v
```

## Exercises

Exercises are maintained as a separate Go project and can be run independently of the backend:

```bash
cd exercises
go test ./...
```

## Commands

```bash
cd backend

# Run API
go run cmd/api/main.go

# Build
go build -o bin/api.exe ./cmd/api

# Run migrations
migrate -path migrations -database "postgres://..." up

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs

# Full Docker stack
docker compose up --build

# Development database and Redis only
docker compose -f ../docker-compose.dev.yml up -d
```
