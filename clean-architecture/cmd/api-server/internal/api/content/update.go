package content

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UpdateContentRequest はコンテンツ更新リクエストの構造体
type UpdateContentRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Body        string `json:"body" binding:"required,min=1"`
	ContentType string `json:"content_type" binding:"required,oneof=article blog news page"`
	Author      string `json:"author" binding:"required,min=1,max=100"`
}

// Update はコンテンツを更新するHTTPハンドラー
func (api *ContentAPI) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なIDです",
			"details": err.Error(),
		})
		return
	}

	var req UpdateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なリクエストです",
			"details": err.Error(),
		})
		return
	}

	// 既存コンテンツ取得
	content, err := api.repo.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"code":    http.StatusNotFound,
				"message": "指定されたコンテンツが見つかりません",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの取得に失敗しました",
			"details": err.Error(),
		})
		return
	}

	// コンテンツ更新
	if err := content.Update(req.Title, req.Body, req.ContentType, req.Author); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なリクエストです",
			"details": err.Error(),
		})
		return
	}

	// DB保存
	if err := api.repo.Update(c.Request.Context(), content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの更新に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, content)
}
