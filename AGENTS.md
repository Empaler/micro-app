# Project Tasks & Commands

## Development

### Backend (Go)
```bash
cd backend

# Run the API locally (requires PostgreSQL running)
go run cmd/api/main.go

# Run unit tests
go test ./internal/... -v

# Run all tests
go test ./... -v

# Build binary
go build -o bin/api.exe cmd/api/main.go

# Add new dependency
go get github.com/package/name

# Generate Swagger docs
swag init -g cmd/api/main.go -o docs
```

### Frontend (React + Vite)
```bash
cd frontend

# Install dependencies
npm install

# Run dev server
npm run dev

# Build for production
npm run build
```

### Exercises (Go)
```bash
cd exercises
go test ./...
```

### Database
```bash
# Run migrations (automatic on startup)
# Manual migration command:
cd backend
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" down
```

## Docker

```bash
# Start database and Redis for local development
docker compose -f docker-compose.dev.yml up -d

# Build and run all services in Docker
docker compose up --build

# Run in detached mode
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down

# Stop development database and Redis
docker compose -f docker-compose.dev.yml down

# Rebuild a specific service
docker compose build backend
docker compose build frontend
```

## Ports
- Frontend: http://localhost:80
- API: http://localhost:8080
- Swagger: http://localhost:8080/swagger/index.html
- PostgreSQL: localhost:5432

## Database Credentials
- User: `postgres`
- Password: `postgres`
- Database: `movies`

## VS Code Debug Configs
- **Launch API** - Run the backend server
- **Test: Unit Tests** - Run unit tests in `internal/...`
- **Test: Integration Tests** - Run testcontainer tests
- **Test: All Tests** - Run all tests

## Adding New Entities

1. Create the entity, service, adapter, router, and tests in `backend/internal/{name}/`
2. Add route registration in `backend/internal/server/server.go`
3. Add migration file in `backend/migrations/`
8. Add Swagger annotations to handler
9. Regenerate docs from `backend/`: `swag init -g cmd/api/main.go -o docs`
10. Add frontend components if needed
