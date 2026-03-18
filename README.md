# Movie Database API

A hexagonal architecture REST API for managing a movie collection with a React frontend.

## Features

- **CRUD Operations**: Create, read, update, and delete movies
- **Movie Fields**: title, year, type (movie/series), resolution (SD/HD/FHD/4K), actors, rating
- **Business Rules**:
  - Adult content is not allowed
  - Year must be between 1888 and current year
  - Rating must be between 0 and 10
  - Valid type and resolution values

## Architecture

```
├── cmd/api/                 # Application entry point
├── internal/
│   ├── domain/             # Business logic
│   │   ├── entity/        # Movie entity with validation
│   │   └── service/       # Movie service (use cases)
│   ├── port/              # Interfaces
│   │   ├── http/          # HTTP handler interface
│   │   └── storage/       # Repository interface
│   ├── adapter/           # Implementations
│   │   ├── http/gin/     # Gin HTTP router
│   │   └── storage/postgres/  # PostgreSQL adapter
│   └── config/            # Configuration
├── tests/
│   ├── integration/       # Testcontainer integration tests
│   └── mocks/             # Mock implementations
└── frontend/             # React + Vite frontend
```

## Tech Stack

- **Backend**: Go 1.21+, Gin, GORM, PostgreSQL
- **Frontend**: React 18, Vite, Axios
- **Testing**: Go testing, testcontainers-go, testify

## Quick Start

### With Docker Compose

```bash
docker-compose up --build
```

Services:
- Frontend: http://localhost:80
- API: http://localhost:8080/api/movies
- Swagger: http://localhost:8080/swagger/index.html
- PostgreSQL: localhost:5432

### Local Development

**Prerequisites:**
- Go 1.21+
- Node.js 18+
- PostgreSQL 15

**Backend:**
```bash
# Create database
createdb movies

# Run API
go run cmd/api/main.go
```

**Frontend:**
```bash
cd frontend
npm install
npm run dev
```

## API Endpoints

| Method | Endpoint      | Description           |
|--------|---------------|-----------------------|
| GET    | /api/movies   | List all movies       |
| GET    | /api/movies/:id | Get movie by ID    |
| POST   | /api/movies   | Create new movie     |
| PUT    | /api/movies/:id | Update movie       |
| DELETE | /api/movies/:id | Delete movie       |

### Example Request

```bash
curl -X POST http://localhost:8080/api/movies \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Matrix",
    "year": 1999,
    "type": "movie",
    "resolution": "FHD",
    "actors": "Keanu Reeves, Laurence Fishburne",
    "rating": 8.7,
    "isAdult": false
  }'
```

## Testing

```bash
# Unit tests
go test ./internal/... -v

# Integration tests (requires Docker)
go test ./tests/integration/... -v

# All tests
go test ./... -v
```

## VS Code

Open the project in VS Code and use the debug configurations in `.vscode/launch.json`:

- **Launch API** - Run the backend server
- **Test: Unit Tests** - Run unit tests
- **Test: Integration Tests** - Run testcontainer tests
