package application

import (
	"context"
	"errors"

	"go-api-server-sample/internal/domain/entities"
	"go-api-server-sample/internal/domain/repositories"
	"gorm.io/gorm"
)

type GetContentUseCase struct {
	contentRepo repositories.ContentRepository
}

type GetContentRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

type GetContentResponse struct {
	*entities.Content
}

var ErrContentNotFound = errors.New("指定されたコンテンツが見つかりません")

func NewGetContentUseCase(contentRepo repositories.ContentRepository) *GetContentUseCase {
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
