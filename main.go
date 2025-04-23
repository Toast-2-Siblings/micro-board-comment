package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Toast-2-Siblings/micro-board-comment/config"
	"github.com/Toast-2-Siblings/micro-board-comment/database"
	"github.com/Toast-2-Siblings/micro-board-comment/server"
)

func main() {
	ctx, cancle := context.WithCancel(context.Background())

	// Initialize the logger
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// Load the configuration
	log.Println("Loading configuration...")
	if _, err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	// Initialize the database
	log.Println("Initializing database...")
	if err := database.InitDatabase(); err != nil {
		log.Fatalf("Failed to initialize database: %v\n", err)
	}

	// Migrate the database
	log.Println("Initializing database migration...")
	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to migrate database: %v\n", err)
	}

	// Config the Server
	server := server.NewServer(&server.ServerConfig{
		Port: "8080",
	}, ctx)

	if err := server.Init(); err != nil {
		log.Fatalf("Failed to initialize server: %v\n", err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Run(); err != nil {
			log.Fatalf("Failed to run server: %v\n", err)
		}
	}()

	<- c
	server.Shutdown(ctx)
	cancle()

	log.Println("Server shutdown gracefully")
}
