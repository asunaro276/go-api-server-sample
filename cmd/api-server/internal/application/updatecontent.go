package application

import (
	"context"
	"errors"

	"go-api-server-sample/internal/domain/entities"
	"gorm.io/gorm"
)

type ContentUpdater interface {
	GetByID(ctx context.Context, id uint) (*entities.Content, error)
	Update(ctx context.Context, content *entities.Content) error
}

type UpdateContentUseCase struct {
	contentRepo ContentUpdater
}

type UpdateContentRequest struct {
	ID          uint   `json:"-"`
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Body        string `json:"body" binding:"required,min=1"`
	ContentType string `json:"content_type" binding:"required,oneof=article blog news page"`
	Author      string `json:"author" binding:"required,min=1,max=100"`
}

type UpdateContentResponse struct {
	*entities.Content
}

func NewUpdateContentUseCase(contentRepo ContentUpdater) *UpdateContentUseCase {
	return &UpdateContentUseCase{
		contentRepo: contentRepo,
	}
}

func (uc *UpdateContentUseCase) Execute(ctx context.Context, req *UpdateContentRequest) (*UpdateContentResponse, error) {
	content, err := uc.contentRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrContentNotFound
		}
		return nil, err
	}

	if err := content.Update(req.Title, req.Body, req.ContentType, req.Author); err != nil {
		return nil, err
	}

	if err := uc.contentRepo.Update(ctx, content); err != nil {
		return nil, err
	}

	return &UpdateContentResponse{Content: content}, nil
}
