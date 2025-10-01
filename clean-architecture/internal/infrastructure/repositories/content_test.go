package repositories

import (
	"context"
	"fmt"
	"testing"
	"time"

	"go-api-server-sample/clean-architecture/internal/domain/entities"
	"go-api-server-sample/clean-architecture/internal/domain/repositories"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type ContentRepositoryTestSuite struct {
	suite.Suite
	container *postgres.PostgresContainer
	db        *gorm.DB
	repo      repositories.ContentRepository
}

func (suite *ContentRepositoryTestSuite) SetupSuite() {
	ctx := context.Background()

	// PostgreSQLコンテナ起動（カスタムwait strategyを使用）
	container, err := postgres.Run(ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	suite.Require().NoError(err)
	suite.container = container

	// DB接続とマイグレーション実行
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	suite.Require().NoError(err)

	// この時点でPostgreSQLは確実に接続可能
	suite.db, err = gorm.Open(postgresDriver.Open(connStr), &gorm.Config{})
	suite.Require().NoError(err)

	err = suite.db.AutoMigrate(&entities.Content{})
	suite.Require().NoError(err)

	suite.repo = NewContentRepository(suite.db)
}

func (suite *ContentRepositoryTestSuite) TearDownSuite() {
	ctx := context.Background()
	if suite.container != nil {
		suite.container.Terminate(ctx)
	}
}

func (suite *ContentRepositoryTestSuite) SetupSubTest() {
	// テストデータクリーンアップ
	suite.db.Exec("DELETE FROM contents")
}

func (suite *ContentRepositoryTestSuite) TestCreate() {
	suite.Run("正常にコンテンツを作成できる", func() {
		ctx := context.Background()
		content, _ := entities.NewContent("テストタイトル", "テスト本文", "article", "テスト作成者")

		err := suite.repo.Create(ctx, content)

		assert.NoError(suite.T(), err)
		assert.Greater(suite.T(), content.ID, uint(0))
		assert.NotZero(suite.T(), content.CreatedAt)
		assert.NotZero(suite.T(), content.UpdatedAt)
	})
}

func (suite *ContentRepositoryTestSuite) TestGetByID() {
	suite.Run("存在するIDでコンテンツを取得できる", func() {
		ctx := context.Background()

		// テストデータ作成
		original, _ := entities.NewContent("テストタイトル", "テスト本文", "article", "テスト作成者")
		err := suite.repo.Create(ctx, original)
		suite.Require().NoError(err)

		// 取得テスト
		retrieved, err := suite.repo.GetByID(ctx, original.ID)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), original.ID, retrieved.ID)
		assert.Equal(suite.T(), original.Title, retrieved.Title)
		assert.Equal(suite.T(), original.Body, retrieved.Body)
		assert.Equal(suite.T(), original.ContentType, retrieved.ContentType)
		assert.Equal(suite.T(), original.Author, retrieved.Author)
	})

	suite.Run("存在しないIDではエラーになる", func() {
		ctx := context.Background()

		_, err := suite.repo.GetByID(ctx, 99999)

		assert.Error(suite.T(), err)
	})
}

func (suite *ContentRepositoryTestSuite) TestList() {
	suite.Run("フィルタなしで全件取得できる", func() {
		ctx := context.Background()

		// テストデータ作成
		contents := []*entities.Content{
			func() *entities.Content {
				c, _ := entities.NewContent("記事1", "本文1", "article", "作成者A")
				return c
			}(),
			func() *entities.Content {
				c, _ := entities.NewContent("ブログ1", "本文2", "blog", "作成者B")
				return c
			}(),
			func() *entities.Content {
				c, _ := entities.NewContent("記事2", "本文3", "article", "作成者A")
				return c
			}(),
		}

		for _, content := range contents {
			err := suite.repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		// 取得テスト
		filters := repositories.NewContentFilters()
		result, total, err := suite.repo.List(ctx, filters)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(3), total)
		assert.Len(suite.T(), result, 3)
	})

	suite.Run("コンテンツタイプでフィルタリングできる", func() {
		ctx := context.Background()

		// テストデータ作成（上記と同じ）
		contents := []*entities.Content{
			func() *entities.Content {
				c, _ := entities.NewContent("記事1", "本文1", "article", "作成者A")
				return c
			}(),
			func() *entities.Content {
				c, _ := entities.NewContent("ブログ1", "本文2", "blog", "作成者B")
				return c
			}(),
			func() *entities.Content {
				c, _ := entities.NewContent("記事2", "本文3", "article", "作成者A")
				return c
			}(),
		}

		for _, content := range contents {
			err := suite.repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		// articleタイプでフィルタリング
		filters := repositories.NewContentFilters()
		articleType := "article"
		filters.ContentType = &articleType

		result, total, err := suite.repo.List(ctx, filters)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(2), total)
		assert.Len(suite.T(), result, 2)
		for _, content := range result {
			assert.Equal(suite.T(), "article", content.ContentType)
		}
	})

	suite.Run("作成者でフィルタリングできる", func() {
		ctx := context.Background()

		// テストデータ作成（上記と同じ）
		contents := []*entities.Content{
			func() *entities.Content {
				c, _ := entities.NewContent("記事1", "本文1", "article", "作成者A")
				return c
			}(),
			func() *entities.Content {
				c, _ := entities.NewContent("ブログ1", "本文2", "blog", "作成者B")
				return c
			}(),
			func() *entities.Content {
				c, _ := entities.NewContent("記事2", "本文3", "article", "作成者A")
				return c
			}(),
		}

		for _, content := range contents {
			err := suite.repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		// 作成者Aでフィルタリング
		filters := repositories.NewContentFilters()
		author := "作成者A"
		filters.Author = &author

		result, total, err := suite.repo.List(ctx, filters)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(2), total)
		assert.Len(suite.T(), result, 2)
		for _, content := range result {
			assert.Equal(suite.T(), "作成者A", content.Author)
		}
	})

	suite.Run("ページネーションが機能する", func() {
		ctx := context.Background()

		// 5件のテストデータ作成
		for i := 1; i <= 5; i++ {
			content, _ := entities.NewContent(
				fmt.Sprintf("記事%d", i),
				fmt.Sprintf("本文%d", i),
				"article",
				"作成者",
			)
			err := suite.repo.Create(ctx, content)
			suite.Require().NoError(err)
		}

		// 最初の2件取得
		filters := repositories.NewContentFilters()
		filters.Limit = 2
		filters.Offset = 0

		result, total, err := suite.repo.List(ctx, filters)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(5), total)
		assert.Len(suite.T(), result, 2)

		// 次の2件取得
		filters.Offset = 2
		result2, total2, err := suite.repo.List(ctx, filters)

		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), int64(5), total2)
		assert.Len(suite.T(), result2, 2)

		// 異なるコンテンツが取得されることを確認
		assert.NotEqual(suite.T(), result[0].ID, result2[0].ID)
	})
}

func (suite *ContentRepositoryTestSuite) TestUpdate() {
	suite.Run("正常にコンテンツを更新できる", func() {
		ctx := context.Background()

		// テストデータ作成
		content, _ := entities.NewContent("元のタイトル", "元の本文", "article", "元の作成者")
		err := suite.repo.Create(ctx, content)
		suite.Require().NoError(err)

		// 更新
		err = content.Update("新しいタイトル", "新しい本文", "blog", "新しい作成者")
		suite.Require().NoError(err)

		err = suite.repo.Update(ctx, content)
		assert.NoError(suite.T(), err)

		// 更新確認
		updated, err := suite.repo.GetByID(ctx, content.ID)
		suite.Require().NoError(err)

		assert.Equal(suite.T(), "新しいタイトル", updated.Title)
		assert.Equal(suite.T(), "新しい本文", updated.Body)
		assert.Equal(suite.T(), "blog", updated.ContentType)
		assert.Equal(suite.T(), "新しい作成者", updated.Author)
		assert.True(suite.T(), updated.UpdatedAt.After(updated.CreatedAt))
	})
}

func (suite *ContentRepositoryTestSuite) TestDelete() {
	suite.Run("正常にコンテンツを削除できる", func() {
		ctx := context.Background()

		// テストデータ作成
		content, _ := entities.NewContent("削除対象", "削除対象本文", "article", "作成者")
		err := suite.repo.Create(ctx, content)
		suite.Require().NoError(err)

		// 削除
		err = suite.repo.Delete(ctx, content.ID)
		assert.NoError(suite.T(), err)

		// 削除確認（ソフトデリートなのでエラーになる）
		_, err = suite.repo.GetByID(ctx, content.ID)
		assert.Error(suite.T(), err)
	})
}

func TestContentRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ContentRepositoryTestSuite))
}
