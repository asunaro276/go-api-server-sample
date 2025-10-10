# Go API Server - Clean Architecture

Clean Architectureの原則に基づいて構築されたシンプルなGo言語のRESTful APIサーバーです。

## プロジェクト概要

このプロジェクトは、保守性と拡張性を重視した3層アーキテクチャを採用しています：

- **API層** (`cmd/api-server/internal/api/`): HTTPリクエストの処理とビジネスロジック
- **Domain層** (`internal/domain/`): エンティティとリポジトリインターフェース定義
- **Infrastructure層** (`internal/infrastructure/`): データベースアクセスなど外部依存の実装

## 技術スタック

- **言語**: Go 1.23
- **Webフレームワーク**: Gin
- **ORM**: GORM
- **データベース**: PostgreSQL
- **テスト**: testify, testcontainers-go
- **リント**: golangci-lint
- **開発ツール**: air (ホットリロード)

## クイックスタート

### 前提条件

- Go 1.23以上
- Docker & Docker Compose
- Make

### セットアップ

```bash
# 環境ファイル作成、Docker起動、マイグレーション実行を一括で実行
make quickstart
```

上記コマンドで以下が自動実行されます：
1. `.env`ファイルの作成
2. PostgreSQLコンテナの起動
3. データベースマイグレーションの実行

### アプリケーションの起動

```bash
# 通常実行
make run

# 開発モード（ホットリロード）
make dev
```

サーバーは `http://localhost:8080` で起動します。

### ヘルスチェック

```bash
curl http://localhost:8080/health
```

## ディレクトリ構成

```
.
├── cmd/api-server/              # アプリケーションエントリーポイント
│   ├── main.go                 # メイン関数
│   ├── internal/               # API層とインフラ設定
│   │   ├── api/               # API層（HTTPハンドラー）
│   │   │   ├── content/       # コンテンツ関連API
│   │   │   └── health/        # ヘルスチェックAPI
│   │   ├── container/         # 依存性注入コンテナ
│   │   └── middleware/        # HTTPミドルウェア
│   └── test/                  # APIテスト
│       ├── integration/       # 統合テスト
│       └── performance/       # パフォーマンステスト
├── internal/                   # 内部パッケージ
│   ├── domain/                # ドメイン層
│   │   ├── entities/          # ドメインエンティティ
│   │   └── repositories/      # リポジトリインターフェース
│   └── infrastructure/        # インフラストラクチャ層
│       ├── database/          # データベース接続管理
│       └── repositories/      # リポジトリ実装
├── config/                     # 設定管理
├── scripts/                    # 開発・運用スクリプト
└── specs/                      # 仕様書・設計ドキュメント
```

## 開発コマンド

### アプリケーション実行

```bash
# 通常実行
make run

# 開発モード（自動リロード）
make dev

# ヘルプを表示
make help
```

### ビルド

```bash
# ローカルビルド
make build
```

### テスト

```bash
# 全テスト実行
make test

# カバレッジ付きテスト
make test-coverage

# 統合テスト
make test-integration

# パフォーマンステスト
make test-performance

# 全種類のテスト実行
make test-all
```

### コード品質

```bash
# リント実行
make lint

# コードフォーマット
make fmt

# go vet実行
make vet

# 全チェック（lint + vet + test）
make check

# CIパイプライン
make ci
```

### データベース

```bash
# マイグレーション実行
make migrate

# データベースリセット（開発環境のみ）
make migrate-reset
```

### Docker

```bash
# PostgreSQLコンテナ起動
make docker-up

# コンテナ停止
make docker-down

# ログ確認
make docker-logs
```

### 開発ツール

```bash
# 開発ツールインストール
make install-tools
```

## API エンドポイント

### ヘルスチェック

```bash
GET /health
```

### コンテンツ管理

```bash
# コンテンツ一覧取得
GET /contents

# コンテンツ取得
GET /contents/:id

# コンテンツ作成
POST /contents
Content-Type: application/json

{
  "title": "タイトル",
  "body": "本文"
}

# コンテンツ更新
PUT /contents/:id
Content-Type: application/json

{
  "title": "更新されたタイトル",
  "body": "更新された本文"
}

# コンテンツ削除
DELETE /contents/:id
```

## アーキテクチャ

### レイヤー間の依存関係

```
API層 → Domain層
   ↓
Infrastructure層
```

### 依存関係の原則

1. **内側への依存のみ**: 外側の層から内側の層への依存のみ許可
2. **抽象化による依存性逆転**: Infrastructure層への依存はインターフェース経由
3. **Domain層の独立性**: Domain層は他の層に依存しない

詳細は `.claude/memories/dependency.md` を参照してください。

## テスト方針

### テスト種類

- **単体テスト**: 各層のロジックを独立してテスト
- **統合テスト**: API層からリポジトリまでの統合動作をテスト
- **パフォーマンステスト**: APIのレスポンスタイムと負荷耐性をテスト

### テストの書き方

- `testify/suite` を使用したテストスイート
- リポジトリテストでは `testcontainers-go` を使用して実際のPostgreSQLでテスト
- モックは `mockery` で自動生成

詳細は `.claude/memories/testing.md` を参照してください。

## 環境変数

`.env`ファイルで以下の環境変数を設定します：

```bash
# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=go_api_db
DB_SSLMODE=disable

# サーバー設定
SERVER_PORT=8080
GIN_MODE=debug
```

## トラブルシューティング

### データベース接続エラー

```bash
# PostgreSQLコンテナが起動しているか確認
docker ps

# コンテナが起動していない場合
make docker-up
```

### ポートが既に使用されている

```bash
# ポート8080を使用しているプロセスを確認
lsof -i :8080

# または別のポートを使用
SERVER_PORT=8081 make run
```

## 開発ガイドライン

このプロジェクトでは、コーディング規約やアーキテクチャルールを `.claude/memories/` ディレクトリで管理しています：

- **dependency.md**: レイヤー間の依存関係ルール
- **naming.md**: ファイル・パッケージ・変数の命名規約
- **organization.md**: ディレクトリ構成とファイル配置ルール
- **testing.md**: テスト方針とベストプラクティス
- **git.md**: Git運用ガイドライン
- **security.md**: セキュリティガイドライン

## ライセンス

このプロジェクトはサンプルプロジェクトです。

## 貢献

プルリクエストを歓迎します。大きな変更の場合は、まずissueを開いて変更内容を議論してください。

### ブランチ戦略

```bash
# 機能追加
feature/###-[実装機能名]

# 例
feature/001-user-management
feature/002-authentication
```

詳細は `.claude/memories/git.md` を参照してください。
