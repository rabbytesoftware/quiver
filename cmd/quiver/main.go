package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbytesoftware/quiver/internal/config"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/metadata"
	"github.com/rabbytesoftware/quiver/internal/server"
	"github.com/rabbytesoftware/quiver/internal/ui"
	"github.com/rabbytesoftware/quiver/packages"
)

func main() {
	// Initialize metadata
	if err := metadata.Load(); err != nil {
		log.Fatalf("Failed to load metadata: %v", err)
	}

	// Initialize configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := logger.New(cfg.Logger)

	// Show welcome message
	ui.ShowWelcome()

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("Received shutdown signal, shutting down gracefully...")
		cancel()
	}()

	// Initialize package manager
	logger.Info("Initializing package manager...")
	pkgManager := packages.NewArrowsServer(cfg.Packages.Repository, logger.WithService("pkgsServer"))
	// if err := pkgManager.Initialize(ctx); err != nil {
	// 	logger.Fatal("Failed to initialize package manager: %v", err)
	// }
	
	pkgManager.Load("./template/cs2.yaml")

	// Initialize and start server
	logger.Info("Starting Quiver server...")
	srv := server.New(cfg.Server, pkgManager, logger)
	if err := srv.Start(ctx); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}

	// Wait for context cancellation
	<-ctx.Done()
	logger.Info("Quiver server stopped")
} 