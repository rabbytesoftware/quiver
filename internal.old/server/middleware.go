package server

import (
	"time"

	"github.com/gin-gonic/gin"
)

// loggingMiddleware logs HTTP requests for Gin
func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status and size
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Build full path
		if raw != "" {
			path = path + "?" + raw
		}

		// Log the request
		s.logger.Info("[%s] %s %s %d %v %dB",
			c.Request.Method,
			path,
			c.ClientIP(),
			status,
			latency,
			size,
		)
	}
}

// corsMiddleware handles CORS headers for Gin
func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
