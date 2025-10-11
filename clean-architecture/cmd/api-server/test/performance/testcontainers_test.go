package performance

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"go-api-server-sample/internal/domain/entities"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	// 共通のPostgreSQLコンテナとDB接続
	testContainer *postgres.PostgresContainer
	testDB        *gorm.DB
)

// TestMain はパッケージ全体で1回だけ実行され、共通のPostgreSQLコンテナを起動
func TestMain(m *testing.M) {
	ctx := context.Background()

	// PostgreSQLコンテナ起動
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
	if err != nil {
		log.Fatalf("failed to start postgres container: %v", err)
	}
	testContainer = container

	// DB接続確立
	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		log.Fatalf("failed to get connection string: %v", err)
	}

	db, err := gorm.Open(postgresDriver.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	testDB = db

	// マイグレーション実行
	if err := db.AutoMigrate(&entities.Content{}); err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}

	// テスト実行
	code := m.Run()

	// クリーンアップ
	if err := container.Terminate(ctx); err != nil {
		log.Printf("failed to terminate container: %v", err)
	}

	os.Exit(code)
}

// getDB は共通のDB接続を返す
func getDB() *gorm.DB {
	return testDB
}

// cleanupDB はテストデータをクリーンアップする
func cleanupDB(b *testing.B) {
	if err := testDB.Exec("TRUNCATE TABLE contents RESTART IDENTITY CASCADE").Error; err != nil {
		b.Fatalf("failed to cleanup database: %v", err)
	}
}

// createTestContent はテスト用コンテンツを作成する
func createTestContent(b *testing.B, title, body string) *entities.Content {
	content := &entities.Content{
		Title:       title,
		Body:        body,
		ContentType: "article",
		Author:      "ベンチマーク作成者",
	}
	if err := testDB.Create(content).Error; err != nil {
		b.Fatalf("failed to create test content: %v", err)
	}
	return content
}

// createTestContents は複数のテスト用コンテンツを作成する
func createTestContents(b *testing.B, count int) []*entities.Content {
	contents := make([]*entities.Content, 0, count)
	for i := 0; i < count; i++ {
		content := createTestContent(b,
			fmt.Sprintf("ベンチマーク用コンテンツ %d", i+1),
			fmt.Sprintf("ベンチマーク用本文 %d", i+1),
		)
		contents = append(contents, content)
	}
	return contents
}
