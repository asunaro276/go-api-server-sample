package content

import (
	"net/http"

	"go-api-server-sample/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

// CreateContentRequest はコンテンツ作成リクエストの構造体
type CreateContentRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Body        string `json:"body" binding:"required,min=1"`
	ContentType string `json:"content_type" binding:"required,oneof=article blog news page"`
	Author      string `json:"author" binding:"required,min=1,max=100"`
}

// Create はコンテンツを作成するHTTPハンドラー
func (api *ContentAPI) Create(c *gin.Context) {
	var req CreateContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なリクエストです",
			"details": err.Error(),
		})
		return
	}

	// ドメインエンティティ作成
	content, err := entities.NewContent(req.Title, req.Body, req.ContentType, req.Author)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なリクエストです",
			"details": err.Error(),
		})
		return
	}

	// リポジトリでDB保存
	if err := api.repo.Create(c.Request.Context(), content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの作成に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, content)
}
