package server

import (
	"runtime"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/rabbytesoftware/quiver/internal/logger"
	"github.com/rabbytesoftware/quiver/internal/metadata"
	"github.com/rabbytesoftware/quiver/internal/packages"
	"github.com/rabbytesoftware/quiver/internal/server/response"
)

// Handler handles server-related HTTP requests
type Handler struct {
	pkgManager *packages.ArrowsServer
	logger     *logger.Logger
	startTime  time.Time
}

// NewHandler creates a new server handler instance
func NewHandler(pkgManager *packages.ArrowsServer, logger *logger.Logger) *Handler {
	return &Handler{
		pkgManager: pkgManager,
		logger:     logger.WithService("server-handler"),
		startTime:  time.Now(),
	}
}

// GetServerInfo handles server info requests
func (h *Handler) GetServerInfo(c *gin.Context) {
	info := gin.H{
		"name":        "Quiver", // TODO get public IP and machine Hostname
		"version":     metadata.Project.Version,
		"description": metadata.Project.Description,
		"author":      metadata.Project.Author,
		"license":     metadata.Project.License,
		"url":         metadata.Project.URL,
		"copyright":   metadata.Project.Copyright,
		"build_info": gin.H{
			"go_version": runtime.Version(),
			"os":         runtime.GOOS,
			"arch":       runtime.GOARCH,
		},
	}

	response.Success(c, "Server information retrieved successfully", info)
}

// GetServerStatus handles server status requests
func (h *Handler) GetServerStatus(c *gin.Context) {
	installed := h.pkgManager.GetInstalledArrows()
	repositories := h.pkgManager.GetRepositories()
	
	uptime := time.Since(h.startTime)
	
	status := gin.H{
		"status":           "running",
		"uptime":           uptime.String(),
		"uptime_seconds":   int(uptime.Seconds()),
		"packages_count":   len(installed),
		"repositories_count": len(repositories),
		"memory": gin.H{
			"allocated": formatBytes(getMemStats().Alloc),
			"sys":       formatBytes(getMemStats().Sys),
		},
		"goroutines": runtime.NumGoroutine(),
	}

	response.Success(c, "Server status retrieved successfully", status)
}

// GetServerHealth handles server health check requests
func (h *Handler) GetServerHealth(c *gin.Context) {
	health := gin.H{
		"status": "healthy",
		"checks": gin.H{
			"packages":     "ok",
			"repositories": "ok",
			"memory":       checkMemoryHealth(),
		},
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	response.Success(c, "Server health check completed", health)
}

// GetServerMetrics handles server metrics requests
func (h *Handler) GetServerMetrics(c *gin.Context) {
	installed := h.pkgManager.GetInstalledArrows()
	repositories := h.pkgManager.GetRepositories()
	
	memStats := getMemStats()
	
	metrics := gin.H{
		"packages": gin.H{
			"total_installed": len(installed),
			"by_status":       getPackageStatusCounts(h.pkgManager),
		},
		"repositories": gin.H{
			"total": len(repositories),
		},
		"system": gin.H{
			"uptime_seconds":  int(time.Since(h.startTime).Seconds()),
			"goroutines":      runtime.NumGoroutine(),
			"memory": gin.H{
				"allocated_bytes": memStats.Alloc,
				"sys_bytes":       memStats.Sys,
				"heap_objects":    memStats.HeapObjects,
				"gc_cycles":       memStats.NumGC,
			},
		},
	}

	response.Success(c, "Server metrics retrieved successfully", metrics)
}

// GetServerVersion handles server version requests
func (h *Handler) GetServerVersion(c *gin.Context) {
	version := gin.H{
		"version":     metadata.Project.Version,
		"description": metadata.Project.Description,
		"author":      metadata.Project.Author,
	}

	response.Success(c, "Server version retrieved successfully", version)
}

// Helper functions

func getMemStats() *runtime.MemStats {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return &m
}

func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return ""
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return ""
}

func checkMemoryHealth() string {
	memStats := getMemStats()
	// Simple check: if allocated memory is less than 100MB, consider healthy
	if memStats.Alloc < 100*1024*1024 {
		return "ok"
	}
	return "warning"
}

func getPackageStatusCounts(pkgManager *packages.ArrowsServer) gin.H {
	installed := pkgManager.GetInstalledArrows()
	statusCounts := make(gin.H)
	
	for name := range installed {
		status, err := pkgManager.GetArrowStatus(name)
		if err != nil {
			if count, exists := statusCounts["error"]; exists {
				statusCounts["error"] = count.(int) + 1
			} else {
				statusCounts["error"] = 1
			}
		} else {
			statusKey := string(status)
			if count, exists := statusCounts[statusKey]; exists {
				statusCounts[statusKey] = count.(int) + 1
			} else {
				statusCounts[statusKey] = 1
			}
		}
	}
	
	return statusCounts
} 