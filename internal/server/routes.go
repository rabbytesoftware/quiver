package server

// setupRoutes sets up the HTTP routes
func (s *Server) setupRoutes() {
	// Health routes (root level)
	s.handlers.Health.SetupRoutes(s.gin.Group(""))

	// API v1 routes
	api := s.gin.Group("/api/v1")
	
	// Setup routes for each handler module
	s.handlers.Arrows.SetupRoutes(api)
	s.handlers.Packages.SetupRoutes(api)
	s.handlers.Repositories.SetupRoutes(api)
	s.handlers.Server.SetupRoutes(api)
	s.handlers.Netbridge.SetupRoutes(api)
} 