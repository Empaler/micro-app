package main

import (
	"log"

	_ "movie-api/docs"
	"movie-api/internal/config"
	"movie-api/internal/db"
	"movie-api/internal/redisclient"
	"movie-api/internal/server"
)

func main() {
	cfg := config.Load()

	database, err := db.New(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	redisClient, err := redisclient.NewClientFromConfig(cfg.RedisHost, cfg.RedisPort, cfg.RedisDB)
	if err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}

	server := server.New(cfg.ServerPort, database, redisClient)

	if err := server.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
