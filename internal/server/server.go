package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/config"
	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/netbridge"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/ui"
)

// Server represents the HTTP server (now using Gin)
type Server struct {
	config     config.ServerConfig
	logger     *logger.Logger
	pkgManager *packages.ArrowsServer
	netbridge  *netbridge.Netbridge
	httpServer *http.Server
	gin        *gin.Engine
	handlers   *Handlers
}

// New creates a new server instance using Gin
func New(
	cfg config.ServerConfig,
	pkgManager *packages.ArrowsServer,
	logger *logger.Logger,
) *Server {
	// Set Gin to release mode by default
	gin.SetMode(gin.ReleaseMode)

	// Initialize netbridge
	netbridgeInstance, err := netbridge.NewNetbridge()
	if err != nil {
		logger.Warn("Failed to initialize netbridge: %v (port forwarding will be disabled)", err)
		netbridgeInstance = nil
	}

	s := &Server{
		config:     cfg,
		logger:     logger.WithService("server"),
		pkgManager: pkgManager,
		netbridge:  netbridgeInstance,
		gin:        gin.New(),
	}

	// Initialize handlers
	s.handlers = NewHandlers(s.pkgManager, s.netbridge, s.logger)

	// Setup server components
	s.setupMiddleware()
	s.setupRoutes()

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Handler:      s.gin,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}

	return s
}

// Start starts the HTTP server
func (s *Server) Start(ctx context.Context) error {
	// Start server in goroutine
	go func() {
		s.logger.Info("Starting API server on %s:%d", s.config.Host, s.config.Port)
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Failed to start HTTP server: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Graceful shutdown
	return s.Shutdown()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() error {
	s.logger.Info("Shutting down API server...")
	ui.ShowShutdown()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return s.httpServer.Shutdown(ctx)
}

// GetRouter returns the server's Gin engine (renamed for compatibility)
func (s *Server) GetRouter() *gin.Engine {
	return s.gin
}

// GetConfig returns the server's configuration
func (s *Server) GetConfig() config.ServerConfig {
	return s.config
}

// setupMiddleware sets up HTTP middleware
func (s *Server) setupMiddleware() {
	// Recovery middleware
	s.gin.Use(gin.Recovery())

	// CORS middleware
	s.gin.Use(s.corsMiddleware())

	// Custom logging middleware
	s.gin.Use(s.loggingMiddleware())
}
