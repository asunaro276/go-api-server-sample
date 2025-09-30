package application

import (
	"context"
	"errors"

	"go-api-server-sample/internal/domain/entities"
	"gorm.io/gorm"
)

type ContentDeleter interface {
	GetByID(ctx context.Context, id uint) (*entities.Content, error)
	Delete(ctx context.Context, id uint) error
}

type DeleteContentUseCase struct {
	contentRepo ContentDeleter
}

type DeleteContentRequest struct {
	ID uint `uri:"id" binding:"required,min=1"`
}

func NewDeleteContentUseCase(contentRepo ContentDeleter) *DeleteContentUseCase {
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
