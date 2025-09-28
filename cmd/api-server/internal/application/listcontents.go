package application

import (
	"context"

	"go-api-server-sample/internal/domain/entities"
	"go-api-server-sample/internal/domain/repositories"
)

type ListContentsUseCase struct {
	contentRepo repositories.ContentRepository
}

type ListContentsRequest struct {
	ContentType *string `form:"content_type" binding:"omitempty,oneof=article blog news page"`
	Author      *string `form:"author" binding:"omitempty,max=100"`
	Limit       int     `form:"limit" binding:"omitempty,min=1,max=100"`
	Offset      int     `form:"offset" binding:"omitempty,min=0"`
}

type ListContentsResponse struct {
	Contents []*entities.Content `json:"contents"`
	Total    int64               `json:"total"`
	Limit    int                 `json:"limit"`
	Offset   int                 `json:"offset"`
}

func NewListContentsUseCase(contentRepo repositories.ContentRepository) *ListContentsUseCase {
	return &ListContentsUseCase{
		contentRepo: contentRepo,
	}
}

func (uc *ListContentsUseCase) Execute(ctx context.Context, req *ListContentsRequest) (*ListContentsResponse, error) {
	filters := repositories.NewContentFilters()

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

	contents, total, err := uc.contentRepo.List(ctx, filters)
	if err != nil {
		return nil, err
	}

	return &ListContentsResponse{
		Contents: contents,
		Total:    total,
		Limit:    filters.Limit,
		Offset:   filters.Offset,
	}, nil
}
