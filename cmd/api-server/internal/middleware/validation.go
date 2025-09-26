package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"go-api-server-sample/cmd/api-server/internal/controller/dtos"
)

// ValidationMiddleware provides request validation capabilities
type ValidationMiddleware struct {
	validator *validator.Validate
}

// NewValidationMiddleware creates a new validation middleware
func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{
		validator: validator.New(),
	}
}

// ValidateJSON is a Gin middleware that validates JSON request bodies
func (v *ValidationMiddleware) ValidateJSON() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Continue to the next handler - validation will be done in the controller
		// using Gin's ShouldBindJSON which internally uses the validator
		c.Next()
	})
}

// HandleValidationErrors converts validator errors to proper API responses
func (v *ValidationMiddleware) HandleValidationErrors(err error) *dtos.ErrorResponse {
	var details []dtos.ErrorDetail

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		for _, fieldError := range validationErrors {
			detail := dtos.ErrorDetail{
				Field:   fieldError.Field(),
				Message: getValidationErrorMessage(fieldError),
			}
			details = append(details, detail)
		}
	}

	if len(details) > 0 {
		return dtos.NewValidationErrorResponse("リクエストのバリデーションに失敗しました", details)
	}

	return dtos.NewErrorResponse("VALIDATION_ERROR", "リクエストのバリデーションに失敗しました")
}

// getValidationErrorMessage returns user-friendly error messages for validation errors
func getValidationErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fe.Field() + "は必須です"
	case "email":
		return "有効なメールアドレスを入力してください"
	case "max":
		return fe.Field() + "は" + fe.Param() + "文字以内で入力してください"
	case "min":
		return fe.Field() + "は" + fe.Param() + "文字以上で入力してください"
	default:
		return fe.Field() + "の値が無効です"
	}
}

// ContentTypeValidation ensures content type is application/json for POST/PUT requests
func ContentTypeValidation() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut {
			contentType := c.GetHeader("Content-Type")
			if contentType != "application/json" && contentType != "application/json; charset=utf-8" {
				errorResp := dtos.NewErrorResponse("INVALID_CONTENT_TYPE", "Content-Type must be application/json")
				c.JSON(http.StatusBadRequest, errorResp)
				c.Abort()
				return
			}
		}
		c.Next()
	})
}

// RequestSizeLimit limits the size of request bodies
func RequestSizeLimit(maxSize int64) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		if c.Request.ContentLength > maxSize {
			errorResp := dtos.NewErrorResponse("REQUEST_TOO_LARGE", "リクエストサイズが大きすぎます")
			c.JSON(http.StatusRequestEntityTooLarge, errorResp)
			c.Abort()
			return
		}

		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxSize)
		c.Next()
	})
}