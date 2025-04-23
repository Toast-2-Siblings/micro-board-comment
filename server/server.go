package server

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Toast-2-Siblings/micro-board-comment/config"
	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Port string
}

type Server interface {
	Init() error
	Run() error
	Shutdown(ctx context.Context)
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

func (s *server) setConfig() {
	cfg := config.GetConfig()
	if cfg.Mode == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
}

func (s *server) Init() error {
	s.router = gin.Default()
	s.setConfig()
	
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
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	
	return nil
}

func (s *server) Shutdown(ctx context.Context) {
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	if err := s.server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
