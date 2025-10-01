package repositories

import (
	"context"

	"go-api-server-sample/internal/domain/entities"
)

//go:generate mockery --name=ContentRepository --output=../../testing/mocks

type ContentRepository interface {
	Create(ctx context.Context, content *entities.Content) error
	GetByID(ctx context.Context, id uint) (*entities.Content, error)
	List(ctx context.Context, filters ContentFilters) ([]*entities.Content, int64, error)
	Update(ctx context.Context, content *entities.Content) error
	Delete(ctx context.Context, id uint) error
}

type ContentFilters struct {
	ContentType *string
	Author      *string
	Limit       int
	Offset      int
}

func NewContentFilters() ContentFilters {
	return ContentFilters{
		Limit:  20, // デフォルト取得件数
		Offset: 0,  // デフォルト開始位置
	}
}
