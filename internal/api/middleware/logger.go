package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rabbytesoftware/quiver/internal/core/watcher"
	"github.com/sirupsen/logrus"
)

func WatcherLogger(watcher *watcher.Watcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		bodySize := c.Writer.Size()
		
		if raw != "" {
			path = path + "?" + raw
		}

		message := fmt.Sprintf("%s %s %d %v %s %d",
			method,
			path,
			statusCode,
			latency,
			clientIP,
			bodySize,
		)

		switch {
		case statusCode >= 500:
			watcher.WithFields(logrus.Fields{
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"latency":    latency,
				"client_ip":  clientIP,
				"body_size":  bodySize,
				"type":       "http_request",
			}).Error(message)
		case statusCode >= 400:
			watcher.WithFields(logrus.Fields{
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"latency":    latency,
				"client_ip":  clientIP,
				"body_size":  bodySize,
				"type":       "http_request",
			}).Warn(message)
		default:
			watcher.WithFields(logrus.Fields{
				"method":     method,
				"path":       path,
				"status":     statusCode,
				"latency":    latency,
				"client_ip":  clientIP,
				"body_size":  bodySize,
				"type":       "http_request",
			}).Info(message)
		}
	}
}

func WatcherRecovery(watcher *watcher.Watcher) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the panic with Watcher
				message := fmt.Sprintf("Panic recovered: %v", err)
				watcher.WithFields(logrus.Fields{
					"method":     c.Request.Method,
					"path":       c.Request.URL.Path,
					"client_ip":  c.ClientIP(),
					"type":       "panic_recovery",
					"error":      err,
				}).Error(message)
				
				// Return 500 error
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
