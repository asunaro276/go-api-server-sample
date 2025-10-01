package repositories

import (
	"context"

	"go-api-server-sample/clean-architecture/internal/domain/entities"
	"go-api-server-sample/clean-architecture/internal/domain/repositories"
	"gorm.io/gorm"
)

type contentRepository struct {
	db *gorm.DB
}

func NewContentRepository(db *gorm.DB) repositories.ContentRepository {
	return &contentRepository{
		db: db,
	}
}

func (r *contentRepository) Create(ctx context.Context, content *entities.Content) error {
	return r.db.WithContext(ctx).Create(content).Error
}

func (r *contentRepository) GetByID(ctx context.Context, id uint) (*entities.Content, error) {
	var content entities.Content
	err := r.db.WithContext(ctx).First(&content, id).Error
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (r *contentRepository) List(ctx context.Context, filters repositories.ContentFilters) ([]*entities.Content, int64, error) {
	var contents []*entities.Content
	var total int64

	query := r.db.WithContext(ctx).Model(&entities.Content{})

	if filters.ContentType != nil {
		query = query.Where("content_type = ?", *filters.ContentType)
	}

	if filters.Author != nil {
		query = query.Where("author = ?", *filters.Author)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Limit(filters.Limit).Offset(filters.Offset).Order("created_at DESC").Find(&contents).Error
	if err != nil {
		return nil, 0, err
	}

	return contents, total, nil
}

func (r *contentRepository) Update(ctx context.Context, content *entities.Content) error {
	return r.db.WithContext(ctx).Save(content).Error
}

func (r *contentRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&entities.Content{}, id).Error
}
