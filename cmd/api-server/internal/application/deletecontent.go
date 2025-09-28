package application

import (
	"context"
	"errors"

	"go-api-server-sample/internal/domain/repositories"
	"gorm.io/gorm"
)

type DeleteContentUseCase struct {
	contentRepo repositories.ContentRepository
}

type DeleteContentRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func NewDeleteContentUseCase(contentRepo repositories.ContentRepository) *DeleteContentUseCase {
	return &DeleteContentUseCase{
		contentRepo: contentRepo,
	}
}

func (uc *DeleteContentUseCase) Execute(ctx context.Context, req *DeleteContentRequest) error {
	_, err := uc.contentRepo.GetByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrContentNotFound
		}
		return err
	}

	return uc.contentRepo.Delete(ctx, req.ID)
}
