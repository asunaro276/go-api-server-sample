package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()

			log.Printf("エラーが発生しました: %v", err.Err)

			errorResponse := gin.H{
				"code":      http.StatusInternalServerError,
				"message":   "内部サーバーエラーが発生しました",
				"timestamp": time.Now().Format(time.RFC3339),
			}

			switch e := err.Err.(type) {
			case *gin.Error:
				if e.Type == gin.ErrorTypeBind {
					errorResponse["code"] = http.StatusBadRequest
					errorResponse["message"] = "リクエストデータが正しくありません"
					errorResponse["details"] = e.Error()
				}
			}

			c.AbortWithStatusJSON(errorResponse["code"].(int), errorResponse)
		}
	}
}
