package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port string
}

type Server interface {
	Init() error
	Run() error
	Shutdown(ctx context.Context) error
}

type server struct {
	cfg *ServerConfig

	router *gin.Engine
	server *http.Server
	ctx context.Context
}

func NewServer(cfg *ServerConfig, ctx context.Context) Server {
	return &server{
		cfg: cfg,
		router: nil,
		server: nil,
		ctx: ctx,
	}
}

func (s *server) Init() error {
	s.router = gin.Default()
	s.router.Use(gin.Logger())
	
	s.router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowedMethods:  []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:  []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowCredentials: true,
		ExposedHeaders: []string{"Content-Length", "X-Requested-With"},
		MaxAge: 12 * time.Hour,
	}))

	s.router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	s.server = &http.Server{
		Addr:    ":" + s.cfg.Port,
		Handler: s.router,
	}

	return nil
}

func (s *server) Run() error {
	log.Println("Starting server on Port", s.cfg.Port)
	if err := s.server.ListenAndServe(); err != nil {
		return err
	}
	
	return nil
}

func (s *server) Shutdown(ctx context.Context) error {
	log.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Println("Server forced to shutdown:", err)
		return err
	}

	log.Println("Server exiting")
	return nil
}
