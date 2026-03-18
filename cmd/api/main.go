package main

import (
	"log"

	_ "movie-api/docs"

	"movie-api/internal/config"
	"movie-api/internal/db"
	"movie-api/internal/server"
)

// @title Movie API
// @version 1.0
// @description Movie and Book collection management API
// @host localhost:8080
// @BasePath /

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	server := server.New(cfg.ServerPort, database)

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
