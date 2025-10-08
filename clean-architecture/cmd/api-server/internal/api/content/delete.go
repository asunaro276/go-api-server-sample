package content

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Delete はコンテンツを削除するHTTPハンドラー
func (api *ContentAPI) Delete(c *gin.Context) {
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

	// 存在確認
	_, err = api.repo.GetByID(c.Request.Context(), uint(id))
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

	// 削除実行
	if err := api.repo.Delete(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    http.StatusInternalServerError,
			"message": "コンテンツの削除に失敗しました",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
