package main

import (
	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/api/health"
	"go-api-server-sample/cmd/api-server/internal/infrastructure/repositories"

	"gorm.io/gorm"
)

// Container は依存性注入コンテナ
type Container struct {
	// APIs
	ContentAPI *content.ContentAPI
	HealthAPI  *health.HealthAPI

	// Repositories
	ContentRepository content.ContentRepository
}

// NewContainer は新しいContainerインスタンスを作成する
func NewContainer(db *gorm.DB) *Container {
	container := &Container{}

	container.initRepositories(db)
	container.initAPIs(db)

	return container
}

func (c *Container) initRepositories(db *gorm.DB) {
	c.ContentRepository = repositories.NewContentRepository(db)
}

func (c *Container) initAPIs(db *gorm.DB) {
	c.ContentAPI = content.NewContentAPI(c.ContentRepository)
	c.HealthAPI = health.NewHealthAPI(db)
}
