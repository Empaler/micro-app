package movies

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
	dbName := "movies_test"
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
		CREATE TABLE IF NOT EXISTS movies (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			year INTEGER NOT NULL,
			type VARCHAR(20) NOT NULL,
			resolution VARCHAR(10) NOT NULL,
			actors TEXT,
			rating REAL NOT NULL,
			is_adult BOOLEAN NOT NULL DEFAULT FALSE,
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
	service := NewMovieService(adapter)
	router := &Router{service: service}

	engine := gin.New()
	router.registerRoutes(engine)
	return engine
}

func TestMovieIntegration_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	t.Run("create movie", func(t *testing.T) {
		body := map[string]any{
			"title":      "The Matrix",
			"year":       1999,
			"type":       "movie",
			"resolution": "FHD",
			"actors":     "Keanu Reeves",
			"rating":     8.7,
			"isAdult":    false,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/movies", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		data := resp["data"].(map[string]any)
		assert.Equal(t, "The Matrix", data["title"])
		assert.NotZero(t, data["id"])
	})

	t.Run("get movie by id", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/movies/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))

		data := resp["data"].(map[string]any)
		assert.Equal(t, "The Matrix", data["title"])
	})

	t.Run("get all movies", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/movies", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)
		assert.True(t, resp["success"].(bool))
	})
}

func TestMovieIntegration_UpdateAndDelete(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	db, cleanup := setupTestDB(t)
	defer cleanup()

	router := setupRouter(db)

	createMovie := func(t *testing.T) int64 {
		body := map[string]any{
			"title":      "Original Title",
			"year":       2000,
			"type":       "movie",
			"resolution": "HD",
			"rating":     7.0,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPost, "/api/movies", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
		require.Equal(t, http.StatusCreated, w.Code)
		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]any)
		return int64(data["id"].(float64))
	}

	t.Run("update movie", func(t *testing.T) {
		id := createMovie(t)

		body := map[string]any{
			"title":      "Updated Title",
			"year":       2001,
			"type":       "series",
			"resolution": "FHD",
			"rating":     8.5,
		}
		jsonBody, _ := json.Marshal(body)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/movies/%d", id), bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]any
		json.Unmarshal(w.Body.Bytes(), &resp)
		data := resp["data"].(map[string]any)
		assert.Equal(t, "Updated Title", data["title"])
		assert.Equal(t, "series", data["type"])
	})

	t.Run("delete movie", func(t *testing.T) {
		id := createMovie(t)

		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/movies/%d", id), nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)

		req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/movies/%d", id), nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestMovieIntegration_Validation(t *testing.T) {
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
				"year":       2000,
				"type":       "movie",
				"resolution": "HD",
				"rating":     7.0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "adult content rejected",
			body: map[string]any{
				"title":      "Adult Movie",
				"year":       2000,
				"type":       "movie",
				"resolution": "HD",
				"rating":     7.0,
				"isAdult":    true,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid year",
			body: map[string]any{
				"title":      "Test",
				"year":       1800,
				"type":       "movie",
				"resolution": "HD",
				"rating":     7.0,
			},
			wantStatus: http.StatusBadRequest,
		},
		{
			name: "invalid type",
			body: map[string]any{
				"title":      "Test",
				"year":       2000,
				"type":       "documentary",
				"resolution": "HD",
				"rating":     7.0,
			},
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBody, _ := json.Marshal(tt.body)
			req := httptest.NewRequest(http.MethodPost, "/api/movies", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tt.wantStatus, w.Code)
		})
	}
}
