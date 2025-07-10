package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rabbytesoftware/quiver/internal/cli"
	"github.com/rabbytesoftware/quiver/internal/config"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/metadata"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server"
	"github.com/rabbytesoftware/quiver/internal/ui"
)

// Build-time variables (set via ldflags)
var (
	version   = "dev"
	buildTime = "unknown"
	gitCommit = "unknown"
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

	// Check if we have command line arguments
	args := os.Args[1:] // Skip program name

	if len(args) > 0 {
		// CLI mode - execute command against running server
		runCLIMode(cfg, args)
	} else {
		// Server mode - start and run the server
		runServerMode(cfg)
	}
}

// runCLIMode executes a CLI command against an already running server
func runCLIMode(cfg *config.Config, args []string) {
	// Initialize logger for CLI (minimal logging)
	logger := logger.New(cfg.Logger)
	
	// Create CLI instance
	cliInstance := cli.New(cfg, logger)

	// Execute CLI command directly against running server
	if err := cliInstance.Execute(args); err != nil {
		log.Fatalf("CLI command failed: %v", err)
	}
}

// runServerMode starts and runs the server continuously
func runServerMode(cfg *config.Config) {
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
	pkgManager := packages.NewArrowsServer(
		cfg.Packages.Repositories,
		cfg.Packages.InstallDir,
		cfg.Packages.DatabasePath,
		logger.WithService("pkgsServer"),
	)
	
	// Initialize package manager
	if err := pkgManager.Initialize(); err != nil {
		logger.Fatal("Failed to initialize package manager: %v", err)
	}

	// Initialize and start server
	logger.Info("Starting Quiver API...")
	srv := server.New(cfg.Server, pkgManager, logger)
	if err := srv.Start(ctx); err != nil {
		logger.Fatal("Failed to start server: %v", err)
	}

	// Wait for context cancellation
	<-ctx.Done()
	logger.Info("Quiver stopped")
} 