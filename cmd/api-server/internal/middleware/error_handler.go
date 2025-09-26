package middleware

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"go-api-server-sample/cmd/api-server/internal/controller/dtos"
)

// ErrorHandlerMiddleware provides centralized error handling
type ErrorHandlerMiddleware struct {
	logger *log.Logger
}

// NewErrorHandlerMiddleware creates a new error handler middleware
func NewErrorHandlerMiddleware(logger *log.Logger) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{
		logger: logger,
	}
}

// HandleErrors is a Gin middleware that catches panics and converts them to proper HTTP responses
func (e *ErrorHandlerMiddleware) HandleErrors() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Log the error and stack trace
				e.logger.Printf("Panic recovered: %v\n%s", err, debug.Stack())

				// Return internal server error
				errorResp := dtos.NewErrorResponse("INTERNAL_ERROR", "内部サーバーエラーが発生しました")
				c.JSON(http.StatusInternalServerError, errorResp)
				c.Abort()
			}
		}()

		c.Next()

		// Handle any errors that were set during request processing
		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			e.handleError(c, err)
		}
	})
}

// handleError handles different types of errors
func (e *ErrorHandlerMiddleware) handleError(c *gin.Context, err *gin.Error) {
	// Log the error
	e.logger.Printf("Request error: %v", err.Err)

	// If response was already written, don't try to write again
	if c.Writer.Written() {
		return
	}

	// Convert error to appropriate HTTP response
	switch err.Type {
	case gin.ErrorTypeBind:
		errorResp := dtos.NewErrorResponse("VALIDATION_ERROR", "リクエストのバリデーションに失敗しました")
		c.JSON(http.StatusBadRequest, errorResp)
	case gin.ErrorTypePublic:
		errorResp := dtos.NewErrorResponse("REQUEST_ERROR", err.Error())
		c.JSON(http.StatusBadRequest, errorResp)
	default:
		errorResp := dtos.NewErrorResponse("INTERNAL_ERROR", "内部サーバーエラーが発生しました")
		c.JSON(http.StatusInternalServerError, errorResp)
	}
}

// RecoveryWithLogger returns a middleware that recovers from panics and logs them
func RecoveryWithLogger(logger *log.Logger) gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultWriter, func(c *gin.Context, recovered interface{}) {
		if logger != nil {
			logger.Printf("Panic recovered: %v\n%s", recovered, debug.Stack())
		}

		errorResp := dtos.NewErrorResponse("INTERNAL_ERROR", "内部サーバーエラーが発生しました")
		c.JSON(http.StatusInternalServerError, errorResp)
	})
}

// CORS middleware for cross-origin requests
func CORS() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Header("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})
}

// SecurityHeaders adds security-related headers
func SecurityHeaders() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		c.Header("Content-Security-Policy", "default-src 'self'")
		c.Next()
	})
}