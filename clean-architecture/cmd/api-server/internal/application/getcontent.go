package application

import (
	"context"
	"errors"

	"go-api-server-sample/clean-architecture/internal/domain/entities"
	"gorm.io/gorm"
)

type ContentGetter interface {
	GetByID(ctx context.Context, id uint) (*entities.Content, error)
}

type GetContentUseCase struct {
	contentRepo ContentGetter
}

type GetContentRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type GetContentResponse struct {
	*entities.Content
}

var ErrContentNotFound = errors.New("指定されたコンテンツが見つかりません")

func NewGetContentUseCase(contentRepo ContentGetter) *GetContentUseCase {
	return &GetContentUseCase{
		contentRepo: contentRepo,
	}
}

func (uc *GetContentUseCase) Execute(ctx context.Context, req *GetContentRequest) (*GetContentResponse, error) {
	content, err := uc.contentRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContentNotFound
		}
		return nil, err
	}

	return &GetContentResponse{Content: content}, nil
}
