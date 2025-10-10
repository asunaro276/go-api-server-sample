# レイヤー間の依存関係

このファイルは、Claude Codeがコードを作成・編集する際に遵守すべき依存関係注入の原則について定義します。

## レイヤー間の依存関係ルール

```
API層 → Domain層
   ↓
Infrastructure層
```

### アーキテクチャ概要

このプロジェクトは、シンプルなClean Architecture構成を採用しています：

- **API層** (`cmd/api-server/internal/api/`): HTTPリクエストの処理とレスポンス生成
- **Domain層** (`internal/domain/`): ビジネスロジックとリポジトリインターフェース定義
  - `entities/`: ドメインエンティティとバリデーション
  - `repositories/`: リポジトリインターフェース
- **Infrastructure層** (`internal/infrastructure/`): 外部依存の実装
  - `database/`: データベース接続管理
  - `repositories/`: リポジトリインターフェースの実装

### 絶対遵守事項

1. **内側への依存のみ許可**: 外側の層から内側の層への依存のみ許可
2. **抽象化による依存性逆転**: Infrastructure層への依存はインターフェース経由
3. **Domain層の独立性**: Domain層は他の層に依存しない

### 依存関係の詳細

#### API層の責務
- HTTPリクエストの受け取りとバリデーション
- ビジネスロジックの実行
- リポジトリインターフェースを通じたデータアクセス
- HTTPレスポンスの生成とエラーハンドリング

#### Domain層の責務
- エンティティの定義とバリデーション
- リポジトリインターフェースの定義
- ビジネスルールの表現

#### Infrastructure層の責務
- データベース接続の管理
- リポジトリインターフェースの実装
- 外部サービスとの統合

### 禁止事項（アーキテクチャ違反）

- **Domain層からの外部ライブラリ直接使用** (GORM等のインポート禁止)
- **API層からInfrastructure層の具象型直接使用** (リポジトリインターフェース経由必須)
- **Domain層から上位層への依存** (API層への依存禁止)

### 依存性注入

依存関係は `container` パッケージで管理します：

```go
// container/container.go
type Container struct {
    ContentRepo repositories.ContentRepository
}

func NewContainer(db *gorm.DB) *Container {
    return &Container{
        ContentRepo: infrastructure.NewContentRepository(db),
    }
}
```

## コードレビューポイント

### 必須チェック項目
- import文での依存関係チェック
- インターフェースの適切な定義
- 責務の分離が適切に行われているか
- API層が具象型に依存していないか
- Domain層が外部ライブラリに依存していないか
