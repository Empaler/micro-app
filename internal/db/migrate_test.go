package db

import (
	"context"
	"testing"
	"time"

	"movie-api/internal/config"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestMigrate(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()
	dbName := "test_db"
	dbUser := "postgres"
	dbPassword := "postgres"

	container, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
	)
	require.NoError(t, err)
	t.Cleanup(func() { container.Terminate(ctx) })

	host, err := container.Host(ctx)
	require.NoError(t, err)
	port, err := container.MappedPort(ctx, "5432")
	require.NoError(t, err)

	cfg := &config.Config{
		DBHost:     host,
		DBPort:     port.Port(),
		DBUser:     dbUser,
		DBPassword: dbPassword,
		DBName:     dbName,
	}

	time.Sleep(3 * time.Second)

	db, err := connect(cfg)
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	migrationsDir := "../../migrations"

	t.Run("migrate up creates all tables", func(t *testing.T) {
		err := migrateUp(cfg, migrationsDir)
		require.NoError(t, err)

		var count int
		err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations'")
		require.NoError(t, err)
		require.Equal(t, 2, count, "should have 2 tables (movies and books)")
	})

	t.Run("migrate down drops all tables", func(t *testing.T) {
		err := MigrateDown(cfg, migrationsDir)
		require.NoError(t, err)

		var count int
		err = db.GetContext(ctx, &count, "SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public' AND tablename != 'schema_migrations'")
		require.NoError(t, err)
		require.Equal(t, 0, count, "all tables should be dropped")
	})
}
