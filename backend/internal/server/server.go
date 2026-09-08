package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movie-api/internal/books"
	"movie-api/internal/movies"
	"movie-api/internal/redisclient"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Server struct {
	engine *gin.Engine
	server *http.Server
}

func New(port string, db *sqlx.DB, redisClient *redisclient.Client) *Server {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())

	movies.RegisterRoutes(engine, db, redisClient)
	books.RegisterRoutes(engine, db, redisClient)

	return &Server{
		engine: engine,
		server: &http.Server{
			Addr:    ":" + port,
			Handler: engine,
		},
	}
}

func (s *Server) Start() error {
	s.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	go func() {
		log.Printf("Server starting on :%s", s.server.Addr[1:])
		log.Printf("Swagger docs available at http://localhost:%s/swagger/index.html", s.server.Addr[1:])
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("Server exited gracefully")
	return nil
}
