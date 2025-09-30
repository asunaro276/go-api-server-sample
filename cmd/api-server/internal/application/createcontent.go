package application

import (
	"context"

	"go-api-server-sample/internal/domain/entities"
)

type ContentCreator interface {
	Create(ctx context.Context, content *entities.Content) error
}

type CreateContentUseCase struct {
	contentRepo ContentCreator
}

type CreateContentRequest struct {
	Title       string `json:"title" binding:"required,min=1,max=200"`
	Body        string `json:"body" binding:"required,min=1"`
	ContentType string `json:"content_type" binding:"required,oneof=article blog news page"`
	Author      string `json:"author" binding:"required,min=1,max=100"`
}

type CreateContentResponse struct {
	*entities.Content
}

func NewCreateContentUseCase(contentRepo ContentCreator) *CreateContentUseCase {
	return &CreateContentUseCase{
		contentRepo: contentRepo,
	}
}

func (uc *CreateContentUseCase) Execute(ctx context.Context, req *CreateContentRequest) (*CreateContentResponse, error) {
	content, err := entities.NewContent(req.Title, req.Body, req.ContentType, req.Author)
	if err != nil {
		return nil, err
	}

	if err := uc.contentRepo.Create(ctx, content); err != nil {
		return nil, err
	}

	return &CreateContentResponse{Content: content}, nil
}
