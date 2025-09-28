# 技術調査結果: Go Clean Architecture + Gin + GORM

**調査日**: 2025-09-28
**対象**: Go言語でのClean Architecture実装パターン

## 1. Clean Architecture 4層構成とディレクトリ構造

### 決定: 既存ディレクトリ構造を活用
```
cmd/api-server/           # アプリケーションエントリーポイント
└── internal/             # アプリケーション層
    ├── application/      # UseCase層（アプリケーションサービス）
    ├── container/        # 依存性注入コンテナ
    ├── controller/       # HTTPコントローラ
    └── middleware/       # HTTPミドルウェア

internal/                 # 内部パッケージ
├── domain/              # ドメイン層
│   ├── entities/        # ドメインエンティティ
│   └── repositories/    # リポジトリインターフェース
└── infrastructure/      # インフラストラクチャ層
    ├── database/        # データベース関連
    └── repositories/    # リポジトリ実装
```

**根拠**: プロジェクト憲章で定義済みの構造が標準的なClean Architectureパターンに準拠している

### 検討した代替案
- ヘキサゴナルアーキテクチャ（ポート&アダプター）
- レイヤードアーキテクチャ

**採用理由**: 憲章でClean Architecture厳格遵守が明記されているため

## 2. Ginルーター初期化とミドルウェア設定

### 決定: 段階的ミドルウェア設定パターン
```go
// cmd/api-server/main.go での初期化例
func setupRouter() *gin.Engine {
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()

    // グローバルミドルウェア
    r.Use(gin.Logger())
    r.Use(gin.Recovery())
    r.Use(middleware.CORS())
    r.Use(middleware.RequestID())

    // API v1 グループ
    v1 := r.Group("/api/v1")
    v1.Use(middleware.ErrorHandler())

    return r
}
```

**根拠**:
- エラーハンドリングの一元化
- ミドルウェアの段階的適用による柔軟性
- ヘルスチェックとAPI エンドポイントの分離

### 検討した代替案
- 全ミドルウェアをグローバル適用
- ルート毎の個別ミドルウェア設定

## 3. GORMマイグレーション管理手法

### 決定: Auto Migrate + 初期データSeeder パターン
```go
// infrastructure/database/migrate.go
func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &entities.Content{},
        // 他のエンティティ
    )
}

// infrastructure/database/seeder.go
func SeedInitialData(db *gorm.DB) error {
    // 初期データ投入
}
```

**根拠**:
- 開発初期段階での迅速なスキーマ変更対応
- エンティティとDBスキーマの同期保証
- テスト環境での簡単なセットアップ

### 検討した代替案
- マイグレーションファイル管理（golang-migrate等）
- 手動SQL管理

**採用理由**: 認証なしの単純なCRUDアプリケーションでは複雑なマイグレーション管理は不要

## 4. testcontainers PostgreSQL統合テスト

### 決定: 共通TestSuite + testcontainersパターン
```go
// internal/infrastructure/repositories/testcontainers.go
type DatabaseTestSuite struct {
    suite.Suite
    container *postgres.PostgresContainer
    db        *gorm.DB
}

func (suite *DatabaseTestSuite) SetupSuite() {
    ctx := context.Background()

    // PostgreSQLコンテナ起動
    container, err := postgres.RunContainer(ctx,
        testcontainers.WithImage("postgres:15"),
        postgres.WithDatabase("testdb"),
        postgres.WithUsername("testuser"),
        postgres.WithPassword("testpass"),
    )

    // DB接続とマイグレーション実行
}
```

**根拠**:
- 実際のPostgreSQLを使った統合テスト
- テスト環境の独立性保証
- CI/CDでの再現可能性

### 検討した代替案
- SQLiteインメモリDB
- Docker Compose事前起動

## 5. mockery v3依存性注入とモック生成

### 決定: インターフェースベース + 自動モック生成
```go
//go:generate mockery --name=ContentRepository --output=../../../testing/mocks
type ContentRepository interface {
    Create(ctx context.Context, content *entities.Content) error
    GetByID(ctx context.Context, id uint) (*entities.Content, error)
    List(ctx context.Context) ([]*entities.Content, error)
    Update(ctx context.Context, content *entities.Content) error
    Delete(ctx context.Context, id uint) error
}
```

**根拠**:
- testify/suiteとの統合性
- モック生成の自動化
- 依存性注入の明確化

### 検討した代替案
- 手動モック実装
- GoMock使用

## 6. エラーハンドリングとログ出力

### 決定: カスタムエラー型 + 構造化ログパターン
```go
// pkg/errors/api_error.go
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}

// cmd/api-server/internal/middleware/error.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        if len(c.Errors) > 0 {
            err := c.Errors.Last().Err
            // エラーレスポンス生成
        }
    }
}
```

**根拠**:
- 一貫したエラーレスポンス形式
- ログの構造化による分析性向上
- セキュリティ原則（詳細情報の適切な隠蔽）

## 7. ヘルスチェックエンドポイント実装

### 決定: 多層ヘルスチェックパターン
```go
// cmd/api-server/internal/controller/health.go
type HealthController struct {
    db *gorm.DB
}

func (h *HealthController) Check(c *gin.Context) {
    response := HealthResponse{
        Status: "healthy",
        Database: h.checkDatabase(),
        Timestamp: time.Now(),
    }
    c.JSON(http.StatusOK, response)
}
```

**根拠**:
- データベース接続状態の確認
- システム全体の健全性監視
- 運用時の障害検知

## 8. リクエスト/レスポンス構造体設計

### 決定: DTO + バリデーションタグパターン
```go
// cmd/api-server/internal/controller/dto/content.go
type CreateContentRequest struct {
    Title       string `json:"title" binding:"required,min=1,max=200"`
    Body        string `json:"body" binding:"required,min=1"`
    ContentType string `json:"content_type" binding:"required,oneof=article blog news"`
    Author      string `json:"author" binding:"required,min=1,max=100"`
}

type ContentResponse struct {
    ID          uint      `json:"id"`
    Title       string    `json:"title"`
    Body        string    `json:"body"`
    ContentType string    `json:"content_type"`
    Author      string    `json:"author"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}
```

**根拠**:
- 入力値検証の自動化
- APIレスポンス形式の一貫性
- セキュリティ原則（入力値検証必須）

### 検討した代替案
- エンティティ直接使用
- 手動バリデーション

## 実装時の留意点

### セキュリティ要件
1. **入力値検証**: Ginのbindingタグを使用した自動検証
2. **SQLインジェクション対策**: GORM標準メソッド使用（Raw SQL禁止）
3. **ログ出力**: センシティブデータのマスキング実装

### パフォーマンス要件
1. **レスポンス時間**: < 200ms p95
2. **データベースインデックス**: 検索頻度の高いフィールドに適用
3. **接続プール**: GORM設定での最適化

### テスト戦略
1. **単体テスト**: mockery生成モック使用
2. **統合テスト**: testcontainers使用
3. **契約テスト**: OpenAPI仕様ベース

## 次フェーズでの活用

この調査結果を基に、フェーズ1では以下を作成：
- data-model.md: エンティティ設計とDBスキーマ
- contracts/: OpenAPI仕様とエンドポイント定義
- quickstart.md: 開発環境セットアップ手順
- 契約テスト: 各エンドポイントのテストファイル