package content

import (
	"net/http"

	"go-api-server-sample/internal/domain/entities"

	"github.com/gin-gonic/gin"
)

// ListContentsRequest は一覧取得リクエストの構造体
type ListContentsRequest struct {
	ContentType *string `form:"content_type" binding:"omitempty,oneof=article blog news page"`
	Author      *string `form:"author" binding:"omitempty,max=100"`
	Limit       int     `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset      int     `form:"offset" binding:"omitempty,min=0"`
}

// ListContentsResponse は一覧取得レスポンスの構造体
type ListContentsResponse struct {
	Contents []*entities.Content `json:"contents"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

// List はコンテンツ一覧を取得するHTTPハンドラー
func (api *ContentAPI) List(c *gin.Context) {
	var req ListContentsRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    http.StatusBadRequest,
			"message": "不正なクエリパラメータです",
			"details": err.Error(),
		})
		return
	}

	// フィルタ構築
	filters := NewContentFilters()

	if req.ContentType != nil {
		filters.ContentType = req.ContentType
	}

	if req.Author != nil {
		filters.Author = req.Author
	}

	if req.Limit > 0 {
		filters.Limit = req.Limit
	}

	if req.Offset >= 0 {
		filters.Offset = req.Offset
	}

	// リポジトリから取得
	contents, total, err := api.repo.List(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツ一覧の取得に失敗しました",
			"details": err.Error(),
		})
		return
	}

	response := &ListContentsResponse{
		Contents: contents,
		Total:    total,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	}

	c.JSON(http.StatusOK, response)
}
