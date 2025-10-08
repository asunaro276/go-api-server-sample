package container

import (
	"go-api-server-sample/cmd/api-server/internal/api/content"
	"go-api-server-sample/cmd/api-server/internal/api/health"
	"go-api-server-sample/internal/domain/repositories"
	infraRepos "go-api-server-sample/internal/infrastructure/repositories"

	"gorm.io/gorm"
)

type Container struct {
	// APIs
	ContentAPI *content.ContentAPI
	HealthAPI  *health.HealthAPI

	// Repositories
	ContentRepository repositories.ContentRepository
}

func NewContainer(db *gorm.DB) *Container {
	container := &Container{}

	container.initRepositories(db)
	container.initAPIs(db)

	return container
}

func (c *Container) initRepositories(db *gorm.DB) {
	c.ContentRepository = infraRepos.NewContentRepository(db)
}

func (c *Container) initAPIs(db *gorm.DB) {
	c.ContentAPI = content.NewContentAPI(c.ContentRepository)
	c.HealthAPI = health.NewHealthAPI(db)
}
