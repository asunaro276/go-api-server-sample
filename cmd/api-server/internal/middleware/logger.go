package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerConfig defines the configuration for the logger middleware
type LoggerConfig struct {
	Logger        *log.Logger
	SkipPaths     []string
	TimeFormat    string
	UTC           bool
	SkipBodyPaths []string
}

// StructuredLogEntry represents a structured log entry
type StructuredLogEntry struct {
	Timestamp    string            `json:"timestamp"`
	Method       string            `json:"method"`
	Path         string            `json:"path"`
	StatusCode   int               `json:"status_code"`
	Latency      string            `json:"latency"`
	ClientIP     string            `json:"client_ip"`
	UserAgent    string            `json:"user_agent"`
	RequestID    string            `json:"request_id,omitempty"`
	RequestBody  interface{}       `json:"request_body,omitempty"`
	ResponseBody interface{}       `json:"response_body,omitempty"`
	Error        string            `json:"error,omitempty"`
	Headers      map[string]string `json:"headers,omitempty"`
}

// DefaultLoggerConfig returns a default logger configuration
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		TimeFormat: "2006-01-02T15:04:05.000Z07:00",
		UTC:        true,
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
		SkipBodyPaths: []string{
			"/health",
		},
	}
}

// StructuredLogger returns a structured logging middleware
func StructuredLogger(config LoggerConfig) gin.HandlerFunc {
	if config.Logger == nil {
		config.Logger = log.Default()
	}

	skipPathsMap := make(map[string]bool)
	for _, path := range config.SkipPaths {
		skipPathsMap[path] = true
	}

	skipBodyPathsMap := make(map[string]bool)
	for _, path := range config.SkipBodyPaths {
		skipBodyPathsMap[path] = true
	}

	return gin.HandlerFunc(func(c *gin.Context) {
		// Skip logging for specified paths
		if skipPathsMap[c.Request.URL.Path] {
			c.Next()
			return
		}

		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Prepare log entry
		entry := StructuredLogEntry{
			Timestamp:  formatTimestamp(start, config),
			Method:     c.Request.Method,
			Path:       path,
			StatusCode: c.Writer.Status(),
			Latency:    latency.String(),
			ClientIP:   c.ClientIP(),
			UserAgent:  c.Request.UserAgent(),
		}

		// Add query parameters to path if present
		if raw != "" {
			entry.Path = path + "?" + raw
		}

		// Add request ID if present
		if requestID := c.GetHeader("X-Request-ID"); requestID != "" {
			entry.RequestID = requestID
		}

		// Add error information if present
		if len(c.Errors) > 0 {
			entry.Error = c.Errors.String()
		}

		// Add request/response bodies for non-skip paths
		if !skipBodyPathsMap[c.Request.URL.Path] {
			// Note: Request body logging would require buffering the body
			// which is complex and can impact performance. Skipping for now.
			
			// Add important headers
			entry.Headers = map[string]string{
				"Content-Type":   c.GetHeader("Content-Type"),
				"Authorization":  maskHeader(c.GetHeader("Authorization")),
				"Accept":         c.GetHeader("Accept"),
			}
		}

		// Log the structured entry
		logStructuredEntry(config.Logger, entry)
	})
}

// formatTimestamp formats the timestamp according to the configuration
func formatTimestamp(t time.Time, config LoggerConfig) string {
	if config.UTC {
		t = t.UTC()
	}
	return t.Format(config.TimeFormat)
}

// logStructuredEntry logs the structured entry as JSON
func logStructuredEntry(logger *log.Logger, entry StructuredLogEntry) {
	jsonData, err := json.Marshal(entry)
	if err != nil {
		logger.Printf("Failed to marshal log entry: %v", err)
		// Fallback to simple logging
		logger.Printf("%s %s %d %s %s",
			entry.Method,
			entry.Path,
			entry.StatusCode,
			entry.Latency,
			entry.ClientIP,
		)
		return
	}

	logger.Println(string(jsonData))
}

// maskHeader masks sensitive header values
func maskHeader(value string) string {
	if value == "" {
		return ""
	}
	if len(value) <= 8 {
		return "****"
	}
	return value[:4] + "****" + value[len(value)-4:]
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = generateRequestID()
		}

		c.Header("X-Request-ID", requestID)
		c.Set("RequestID", requestID)
		c.Next()
	})
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return fmt.Sprintf("req_%d", time.Now().UnixNano())
}

// SimpleLogger provides a simple text-based logger middleware
func SimpleLogger(logger *log.Logger) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		if raw != "" {
			path = path + "?" + raw
		}

		logger.Printf("[%s] %s %s %d %v %s",
			start.Format("2006/01/02 15:04:05"),
			c.Request.Method,
			path,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
		)
	})
}