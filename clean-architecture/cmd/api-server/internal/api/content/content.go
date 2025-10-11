package content

import (
	"context"

	"go-api-server-sample/internal/domain/entities"
)

//go:generate mockery --name=ContentRepository --output=../../testing/mocks

// ContentRepository はコンテンツの永続化を担当するリポジトリインターフェース
type ContentRepository interface {
	Create(ctx context.Context, content *entities.Content) error
	GetByID(ctx context.Context, id uint) (*entities.Content, error)
	List(ctx context.Context, filters ContentFilters) ([]*entities.Content, int64, error)
	Update(ctx context.Context, content *entities.Content) error
	Delete(ctx context.Context, id uint) error
}

// ContentFilters はコンテンツ一覧取得時のフィルタ条件
type ContentFilters struct {
	ContentType *string
	Author      *string
	Limit       int
	Offset      int
}

// NewContentFilters はContentFiltersの新しいインスタンスを作成する
func NewContentFilters() ContentFilters {
	return ContentFilters{
		Limit:  20, // デフォルト取得件数
		Offset: 0,  // デフォルト開始位置
	}
}

// ContentAPI はContent関連のHTTPハンドラーを提供する構造体
type ContentAPI struct {
	repo ContentRepository
}

// NewContentAPI はContentAPIの新しいインスタンスを作成する
func NewContentAPI(repo ContentRepository) *ContentAPI {
	return &ContentAPI{
		repo: repo,
	}
}
