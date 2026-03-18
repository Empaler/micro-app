# Project Tasks & Commands

## Development

### Backend (Go)
```bash
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

### Database
```bash
# Run migrations (automatic on startup)
# Manual migration command:
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" up

# Rollback migrations
migrate -path migrations -database "postgres://user:pass@localhost:5432/db?sslmode=disable" down
```

## Docker

```bash
# Build and run all services
docker-compose up --build

# Run in detached mode
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down

# Rebuild a specific service
docker-compose build backend
docker-compose build frontend
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

1. Create entity in `internal/domain/entity/{name}.go`
2. Create service in `internal/domain/service/{name}_service.go`
3. Create repository port in `internal/port/storage/{name}_repo.go`
4. Create adapter in `internal/adapter/storage/postgresadapter/{name}_adapter.go`
5. Create HTTP handler in `internal/adapter/http/ginhandler/{name}_router.go`
6. Add route registration in `main.go`
7. Add migration file in `migrations/`
8. Add Swagger annotations to handler
9. Regenerate docs: `swag init -g cmd/api/main.go -o docs`
10. Add frontend components if needed
