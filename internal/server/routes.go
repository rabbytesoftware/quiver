package server

import (
	"github.com/gorilla/mux"
)

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
	// Health check - direct on root router
	s.router.HandleFunc("/health", s.handlers.Health.HealthHandler).Methods("GET")

	// API v1 routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// Package routes (legacy)
	s.setupPackageRoutes(api)

	// Arrow package management routes
	s.setupArrowRoutes(api)

	// Repository management routes
	s.setupRepositoryRoutes(api)

	// Server management routes
	s.setupServerRoutes(api)

	// Netbridge management routes
	s.setupNetbridgeRoutes(api)
}

// setupPackageRoutes sets up package-related routes
func (s *Server) setupPackageRoutes(api *mux.Router) {
	packageRouter := api.PathPrefix("/packages").Subrouter()

	// Package collection routes
	packageRouter.HandleFunc("", s.handlers.Packages.ListPackagesHandler).Methods("GET")

	// Individual package routes
	packageRouter.HandleFunc("/{id}", s.handlers.Packages.GetPackageHandler).Methods("GET")
	packageRouter.HandleFunc("/{id}/start", s.handlers.Packages.StartPackageHandler).Methods("POST")
	packageRouter.HandleFunc("/{id}/stop", s.handlers.Packages.StopPackageHandler).Methods("POST")
	packageRouter.HandleFunc("/{id}/status", s.handlers.Packages.PackageStatusHandler).Methods("GET")
}

// setupArrowRoutes sets up arrow package management routes
func (s *Server) setupArrowRoutes(api *mux.Router) {
	arrowRouter := api.PathPrefix("/arrows").Subrouter()

	// Search arrows
	arrowRouter.HandleFunc("/search", s.handlers.Arrows.SearchArrowsHandler).Methods("GET")

	// Installed arrows
	arrowRouter.HandleFunc("/installed", s.handlers.Arrows.GetInstalledArrowsHandler).Methods("GET")

	// Individual arrow operations
	arrowRouter.HandleFunc("/{name}/install", s.handlers.Arrows.InstallArrowHandler).Methods("POST")
	arrowRouter.HandleFunc("/{name}/execute", s.handlers.Arrows.ExecuteArrowHandler).Methods("POST")
	arrowRouter.HandleFunc("/{name}/uninstall", s.handlers.Arrows.UninstallArrowHandler).Methods("DELETE")
	arrowRouter.HandleFunc("/{name}/update", s.handlers.Arrows.UpdateArrowHandler).Methods("PUT")
	arrowRouter.HandleFunc("/{name}/validate", s.handlers.Arrows.ValidateArrowHandler).Methods("POST")
	arrowRouter.HandleFunc("/{name}/status", s.handlers.Arrows.GetArrowStatusHandler).Methods("GET")
}

// setupRepositoryRoutes sets up repository management routes
func (s *Server) setupRepositoryRoutes(api *mux.Router) {
	repoRouter := api.PathPrefix("/repositories").Subrouter()

	// Repository collection routes
	repoRouter.HandleFunc("", s.handlers.Repositories.GetRepositoriesHandler).Methods("GET")
	repoRouter.HandleFunc("", s.handlers.Repositories.AddRepositoryHandler).Methods("POST")
	repoRouter.HandleFunc("", s.handlers.Repositories.RemoveRepositoryHandler).Methods("DELETE")
}

// setupServerRoutes sets up server management routes
func (s *Server) setupServerRoutes(api *mux.Router) {
	serverRouter := api.PathPrefix("/server").Subrouter()

	serverRouter.HandleFunc("/info", s.handlers.Server.ServerInfoHandler).Methods("GET")
	serverRouter.HandleFunc("/status", s.handlers.Server.ServerStatusHandler).Methods("GET")
}

// setupNetbridgeRoutes sets up netbridge management routes
func (s *Server) setupNetbridgeRoutes(api *mux.Router) {
	netbridgeRouter := api.PathPrefix("/netbridge").Subrouter()

	// Port management routes
	netbridgeRouter.HandleFunc("/open", s.handlers.Netbridge.OpenPortHandler).Methods("POST")
	netbridgeRouter.HandleFunc("/close", s.handlers.Netbridge.ClosePortHandler).Methods("POST")
	netbridgeRouter.HandleFunc("/open/{port}/{protocol}", s.handlers.Netbridge.OpenPortByURLHandler).Methods("POST")
	netbridgeRouter.HandleFunc("/close/{port}/{protocol}", s.handlers.Netbridge.ClosePortByURLHandler).Methods("POST")
	
	// Automatic port discovery routes
	netbridgeRouter.HandleFunc("/open-auto", s.handlers.Netbridge.OpenPortAutoHandler).Methods("POST")
	netbridgeRouter.HandleFunc("/open-auto/{protocol}", s.handlers.Netbridge.OpenPortAutoByURLHandler).Methods("POST")
	
	// Status and information routes
	netbridgeRouter.HandleFunc("/ports", s.handlers.Netbridge.ListOpenPortsHandler).Methods("GET")
	netbridgeRouter.HandleFunc("/status", s.handlers.Netbridge.GetStatusHandler).Methods("GET")
	netbridgeRouter.HandleFunc("/refresh-ip", s.handlers.Netbridge.RefreshPublicIPHandler).Methods("POST")
} 