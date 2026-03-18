package books

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func setupTestDB(t *testing.T) (*sqlx.DB, func()) {
	ctx := context.Background()
	dbName := "books_test"
	dbUser := "postgres"
	dbPassword := "postgres"

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
	)
	require.NoError(t, err)

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	time.Sleep(2 * time.Second)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port.Port(), dbUser, dbPassword, dbName)

	db, err := sqlx.Connect("pgx", dsn)
	require.NoError(t, err)

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			author VARCHAR(255) NOT NULL,
			release_year INTEGER NOT NULL,
			rating REAL NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`)
	require.NoError(t, err)

	cleanup := func() {
		db.Close()
		container.Terminate(context.Background())
	}

	return db, cleanup
}

func setupRouter(db *sqlx.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	adapter := NewPostgresAdapter(db)
	service := NewBookService(adapter)
	router := &Router{service: service}

	engine := gin.New()
	router.registerRoutes(engine)
	return engine
}

func TestBookIntegration_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	t.Run("create book", func(t *testing.T) {
		body := map[string]any{
			"title":       "The Go Programming Language",
			"author":      "Donovan & Kernighan",
			"releaseYear": 2015,
			"rating":      9.5,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		data := resp["data"].(map[string]any)
		assert.Equal(t, "The Go Programming Language", data["title"])
		assert.NotZero(t, data["id"])
	})

	t.Run("get book by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/books/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		data := resp["data"].(map[string]any)
		assert.Equal(t, "The Go Programming Language", data["title"])
	})

	t.Run("get all books", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/books", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))
	})
}

func TestBookIntegration_UpdateAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	createBook := func(t *testing.T) int64 {
		body := map[string]any{
			"title":       "Original Title",
			"author":      "Original Author",
			"releaseYear": 2020,
			"rating":      7.0,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]any)
		return int64(data["id"].(float64))
	}

	t.Run("update book", func(t *testing.T) {
		id := createBook(t)

		body := map[string]any{
			"title":       "Updated Title",
			"author":      "Updated Author",
			"releaseYear": 2021,
			"rating":      8.5,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/books/%d", id), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]any)
		assert.Equal(t, "Updated Title", data["title"])
		assert.Equal(t, "Updated Author", data["author"])
	})

	t.Run("delete book", func(t *testing.T) {
		id := createBook(t)

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/books/%d", id), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/books/%d", id), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestBookIntegration_Validation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	tests := []struct {
		name       string
		body       map[string]any
		wantStatus int
	}{
		{
			name: "missing title",
			body: map[string]any{
				"author":      "Some Author",
				"releaseYear": 2020,
				"rating":      7.0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "missing author",
			body: map[string]any{
				"title":       "Some Book",
				"releaseYear": 2020,
				"rating":      7.0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid year too early",
			body: map[string]any{
				"title":       "Test",
				"author":      "Author",
				"releaseYear": 500,
				"rating":      7.0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid rating too high",
			body: map[string]any{
				"title":       "Test",
				"author":      "Author",
				"releaseYear": 2020,
				"rating":      15.0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/books", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
