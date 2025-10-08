package content

import (
	"go-api-server-sample/internal/domain/repositories"
)

// ContentAPI はContent関連のHTTPハンドラーを提供する構造体
type ContentAPI struct {
	repo repositories.ContentRepository
}

// NewContentAPI はContentAPIの新しいインスタンスを作成する
func NewContentAPI(repo repositories.ContentRepository) *ContentAPI {
	return &ContentAPI{
		repo: repo,
	}
}
