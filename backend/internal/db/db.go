package db

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"movie-api/internal/config"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func New(cfg *config.Config) (*sqlx.DB, error) {
	return NewWithMigrations(cfg, defaultMigrationsPath())
}

func defaultMigrationsPath() string {
	if _, err := os.Stat("migrations"); err == nil {
		return "migrations"
	}

	backendMigrationsPath := filepath.Join("backend", "migrations")
	if _, err := os.Stat(backendMigrationsPath); err == nil {
		return backendMigrationsPath
	}

	return "migrations"
}

func NewWithMigrations(cfg *config.Config, migrationsPath string) (*sqlx.DB, error) {
	if err := migrateUp(cfg, migrationsPath); err != nil {
		return nil, err
	}
	return connect(cfg)
}

func connect(cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)
	db, err := sqlx.Connect("pgx", dsn)
	if err != nil {
		return nil, err
	}
	db.MapperFunc(func(s string) string { return s })
	return db, nil
}

func migrateUp(cfg *config.Config, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName),
	)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}
	log.Println("Migrations completed successfully")
	return nil
}

func MigrateDown(cfg *config.Config, migrationsPath string) error {
	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsPath),
		fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName),
	)
	if err != nil {
		return err
	}

	version, _, _ := m.Version()
	log.Printf("Current migration version: %d", version)

	count := 0
	for {
		count++
		if err := m.Steps(-1); err != nil {
			if err == migrate.ErrNoChange {
				log.Printf("Ran %d down migrations", count-1)
				break
			}
			return fmt.Errorf("down migration %d failed: %w", count, err)
		}
		version, _, _ = m.Version()
		log.Printf("After down migration %d, version: %d", count, version)
		if version == 0 {
			break
		}
	}
	log.Println("Rollback completed successfully")
	return nil
}
