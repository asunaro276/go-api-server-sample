# データモデル設計: Go API サーバー

**作成日**: 2025-09-28
**基準仕様**: `/specs/001-go-api-crud/spec.md`

## エンティティ設計

### Content（コンテンツ）エンティティ

機能仕様で定義された「コンテンツアイテム」をドメインエンティティとして設計。

#### フィールド定義

| フィールド名 | データ型 | 制約 | 説明 |
|------------|----------|------|------|
| ID | uint | PRIMARY KEY, AUTO_INCREMENT | コンテンツの一意識別子 |
| Title | string | NOT NULL, LENGTH(1-200) | コンテンツのタイトル |
| Body | text | NOT NULL, LENGTH(1-) | コンテンツの本文 |
| ContentType | string | NOT NULL, ENUM('article', 'blog', 'news', 'page') | コンテンツの種類 |
| Author | string | NOT NULL, LENGTH(1-100) | コンテンツ作成者名 |
| CreatedAt | timestamp | NOT NULL, DEFAULT: NOW() | 作成日時 |
| UpdatedAt | timestamp | NOT NULL, DEFAULT: NOW() ON UPDATE | 更新日時 |

#### エンティティ構造（Go実装）

```go
// internal/domain/entities/content.go
package entities

import (
    "time"
    "gorm.io/gorm"
)

type Content struct {
    ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
    Title       string         `gorm:"type:varchar(200);not null" json:"title"`
    Body        string         `gorm:"type:text;not null" json:"body"`
    ContentType string         `gorm:"type:varchar(50);not null" json:"content_type"`
    Author      string         `gorm:"type:varchar(100);not null" json:"author"`
    CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Content) TableName() string {
    return "contents"
}
```

## 検証ルール

### 作成時検証
- **Title**: 必須、1-200文字
- **Body**: 必須、1文字以上
- **ContentType**: 必須、許可値('article', 'blog', 'news', 'page')のいずれか
- **Author**: 必須、1-100文字

### 更新時検証
- 作成時と同じ検証ルール
- ID は変更不可
- CreatedAt は変更不可

### ビジネスルール
- 削除はソフトデリート（DeletedAtにタイムスタンプ設定）
- 重複チェックなし（同じタイトルでも作成可能）
- 権限チェックなし（認証・認可機能なし）

## データベーススキーマ

### PostgreSQLテーブル定義

```sql
CREATE TABLE contents (
    id SERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    body TEXT NOT NULL,
    content_type VARCHAR(50) NOT NULL CHECK (content_type IN ('article', 'blog', 'news', 'page')),
    author VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- インデックス
CREATE INDEX idx_contents_content_type ON contents(content_type);
CREATE INDEX idx_contents_author ON contents(author);
CREATE INDEX idx_contents_created_at ON contents(created_at);
CREATE INDEX idx_contents_deleted_at ON contents(deleted_at);
```

### インデックス戦略

#### 検索パフォーマンス最適化
- **content_type**: コンテンツタイプでの絞り込み検索用
- **author**: 作成者による検索用
- **created_at**: 日付ソート・範囲検索用
- **deleted_at**: ソフトデリートクエリ最適化用

#### 期待されるクエリパターン
1. 全コンテンツ一覧（deleted_at IS NULL）
2. コンテンツタイプ別一覧（content_type = ?）
3. 作成者別一覧（author = ?）
4. 日付範囲検索（created_at BETWEEN ? AND ?）

## 状態遷移

### コンテンツライフサイクル

```
[作成] → [存在] → [更新可能] → [削除済み]
   ↓        ↓         ↓           ↓
 INSERT   SELECT    UPDATE    UPDATE(soft delete)
```

#### 状態説明
- **作成**: 新規コンテンツ作成（POST /api/v1/contents）
- **存在**: 通常の状態、参照・更新可能
- **更新可能**: 任意のフィールド更新可能（PUT /api/v1/contents/:id）
- **削除済み**: ソフトデリート状態、API経由での参照不可

### 状態遷移ルール
- 削除済みコンテンツは復元不可（API機能として提供しない）
- 削除済みコンテンツは一覧・詳細取得から除外
- 更新操作時に削除済みコンテンツを対象にした場合、404エラー

## データ整合性

### 外部キー制約
- なし（単一エンティティのため）

### ユニーク制約
- なし（重複タイトル許可）

### NOT NULL制約
- title, body, content_type, author: 必須フィールド
- created_at, updated_at: GORM自動管理

### チェック制約
- content_type: 許可値のみ受け入れ

## マイグレーション計画

### Auto Migrate使用
```go
// infrastructure/database/migrate.go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(&entities.Content{})
}
```

### 初期データ（オプション）
```go
// infrastructure/database/seeder.go
func SeedSampleData(db *gorm.DB) error {
    sampleContents := []entities.Content{
        {
            Title:       "サンプル記事",
            Body:        "これはサンプルの記事です。",
            ContentType: "article",
            Author:      "管理者",
        },
        // 追加のサンプルデータ
    }

    for _, content := range sampleContents {
        if err := db.Create(&content).Error; err != nil {
            return err
        }
    }
    return nil
}
```

## パフォーマンス考慮事項

### クエリ最適化
- 一覧取得時のページネーション実装推奨
- deleted_at IS NULL条件の効率的な処理
- ORDER BY created_at DESC でのソート最適化

### 予想データ規模
- 初期: 100-1,000件程度
- スケール: 10,000-100,000件を想定
- レスポンス目標: < 200ms p95

### メモリ使用量
- 1件あたり平均 1-5KB を想定
- 一覧取得時のメモリ効率を考慮

## セキュリティ考慮事項

### 入力値サニタイゼーション
- HTMLタグ除去（XSS対策）
- 特殊文字エスケープ
- SQLインジェクション対策（GORM使用）

### 出力制御
- センシティブデータなし（認証情報等含まない）
- DeletedAt フィールドはJSONレスポンスから除外

### アクセス制御
- 認証・認可なし（仕様要件）
- 全データが公開データとして扱われる