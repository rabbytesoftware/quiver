package server

import (
	"github.com/gorilla/mux"
)

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
	// Health check - direct on root router
	s.router.HandleFunc("/health", s.handlers.HealthHandler).Methods("GET")

	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Package routes
	s.setupPackageRoutes(api)

	// Server management routes
	s.setupServerRoutes(api)
}

// setupPackageRoutes sets up package-related routes
func (s *Server) setupPackageRoutes(api *mux.Router) {
	packageRouter := api.PathPrefix("/packages").Subrouter()

	// Package collection routes
	packageRouter.HandleFunc("", s.handlers.ListPackagesHandler).Methods("GET")

	// Individual package routes
	packageRouter.HandleFunc("/{id}", s.handlers.GetPackageHandler).Methods("GET")
	packageRouter.HandleFunc("/{id}/start", s.handlers.StartPackageHandler).Methods("POST")
	packageRouter.HandleFunc("/{id}/stop", s.handlers.StopPackageHandler).Methods("POST")
	packageRouter.HandleFunc("/{id}/status", s.handlers.PackageStatusHandler).Methods("GET")
}

// setupServerRoutes sets up server management routes
func (s *Server) setupServerRoutes(api *mux.Router) {
	serverRouter := api.PathPrefix("/server").Subrouter()

	serverRouter.HandleFunc("/info", s.handlers.ServerInfoHandler).Methods("GET")
	serverRouter.HandleFunc("/status", s.handlers.ServerStatusHandler).Methods("GET")
} 