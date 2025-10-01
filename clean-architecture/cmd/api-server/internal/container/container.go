package container

import (
	"go-api-server-sample/cmd/api-server/internal/application"
	"go-api-server-sample/internal/domain/repositories"
	infraRepos "go-api-server-sample/internal/infrastructure/repositories"

	"gorm.io/gorm"
)

type Container struct {
	// Use Cases
	CreateContentUseCase *application.CreateContentUseCase
	GetContentUseCase    *application.GetContentUseCase
	ListContentsUseCase  *application.ListContentsUseCase
	UpdateContentUseCase *application.UpdateContentUseCase
	DeleteContentUseCase *application.DeleteContentUseCase
	HealthCheckUseCase   *application.HealthCheckUseCase

	// Repositories
	ContentRepository repositories.ContentRepository
}

func NewContainer(db *gorm.DB) *Container {
	container := &Container{}

	container.initRepositories(db)
	container.initUseCases(db)

	return container
}

func (c *Container) initRepositories(db *gorm.DB) {
	c.ContentRepository = infraRepos.NewContentRepository(db)
}

func (c *Container) initUseCases(db *gorm.DB) {
	c.CreateContentUseCase = application.NewCreateContentUseCase(c.ContentRepository)
	c.GetContentUseCase = application.NewGetContentUseCase(c.ContentRepository)
	c.ListContentsUseCase = application.NewListContentsUseCase(c.ContentRepository)
	c.UpdateContentUseCase = application.NewUpdateContentUseCase(c.ContentRepository)
	c.DeleteContentUseCase = application.NewDeleteContentUseCase(c.ContentRepository)
	c.HealthCheckUseCase = application.NewHealthCheckUseCase(db)
}
