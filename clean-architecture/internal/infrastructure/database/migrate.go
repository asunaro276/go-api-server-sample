package database

import (
	"fmt"
	"log"

	"go-api-server-sample/clean-architecture/internal/domain/entities"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	log.Println("マイグレーションを開始します...")

	if err := db.AutoMigrate(&entities.Content{}); err != nil {
		return fmt.Errorf("Contentテーブルのマイグレーションに失敗しました: %w", err)
	}

	if err := createIndexes(db); err != nil {
		return fmt.Errorf("インデックス作成に失敗しました: %w", err)
	}

	log.Println("マイグレーションが完了しました")
	return nil
}

func Reset(db *gorm.DB) error {
	log.Println("データベースリセットを開始します...")

	if err := db.Migrator().DropTable(&entities.Content{}); err != nil {
		return fmt.Errorf("テーブル削除に失敗しました: %w", err)
	}

	log.Println("データベースリセットが完了しました")
	return nil
}

func createIndexes(db *gorm.DB) error {
	indexes := []struct {
		name  string
		query string
	}{
		{
			name:  "idx_contents_content_type",
			query: "CREATE INDEX IF NOT EXISTS idx_contents_content_type ON contents(content_type)",
		},
		{
			name:  "idx_contents_author",
			query: "CREATE INDEX IF NOT EXISTS idx_contents_author ON contents(author)",
		},
		{
			name:  "idx_contents_created_at",
			query: "CREATE INDEX IF NOT EXISTS idx_contents_created_at ON contents(created_at)",
		},
	}

	for _, idx := range indexes {
		log.Printf("インデックス %s を作成しています...", idx.name)
		if err := db.Exec(idx.query).Error; err != nil {
			return fmt.Errorf("インデックス %s の作成に失敗しました: %w", idx.name, err)
		}
	}

	return nil
}

func SeedSampleData(db *gorm.DB) error {
	log.Println("サンプルデータを投入しています...")

	sampleContents := []entities.Content{
		{
			Title:       "サンプル記事1",
			Body:        "これは最初のサンプル記事です。コンテンツ管理システムのテスト用データです。",
			ContentType: "article",
			Author:      "システム管理者",
		},
		{
			Title:       "サンプルブログ1",
			Body:        "ブログ形式のサンプル投稿です。日常的な情報を共有する際に使用します。",
			ContentType: "blog",
			Author:      "ブログ投稿者",
		},
		{
			Title:       "重要なお知らせ",
			Body:        "システムメンテナンスに関する重要なお知らせです。",
			ContentType: "news",
			Author:      "運営チーム",
		},
		{
			Title:       "利用規約",
			Body:        "本サービスの利用規約について説明しています。",
			ContentType: "page",
			Author:      "法務チーム",
		},
	}

	for _, content := range sampleContents {
		if err := db.Create(&content).Error; err != nil {
			return fmt.Errorf("サンプルデータ投入に失敗しました: %w", err)
		}
	}

	log.Printf("サンプルデータ %d 件を投入しました", len(sampleContents))
	return nil
}
