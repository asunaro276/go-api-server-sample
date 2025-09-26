# Quickstart Guide

## Prerequisites
このプロジェクトを実行するために必要な環境：

- Go 1.21以上
- PostgreSQL 13以上
- Make（Makefileの実行用）
- Docker & Docker Compose（任意、ローカル開発用）

## Environment Setup

### 1. Clone Repository
```bash
git clone <repository-url>
cd go-api-server-sample
```

### 2. Environment Variables
`.env`ファイルを作成し、以下の環境変数を設定：

```bash
# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=go_api_server
DB_SSLMODE=disable

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Application Configuration
APP_ENV=development
LOG_LEVEL=debug
```

### 3. Database Setup

#### Option A: Docker Compose (推奨)
```bash
# PostgreSQLコンテナを起動
docker-compose up -d postgres

# データベースが起動するまで待機
sleep 5
```

#### Option B: Local PostgreSQL
```bash
# データベースを作成
createdb go_api_server

# 必要に応じてユーザーを作成
psql -c "CREATE USER your_username WITH ENCRYPTED PASSWORD 'your_password';"
psql -c "GRANT ALL PRIVILEGES ON DATABASE go_api_server TO your_username;"
```

## Build and Run

### 1. Dependencies Installation
```bash
# Go dependenciesを取得
go mod download

# Mockeryのインストール（テスト用）
go install github.com/vektra/mockery/v2@latest
```

### 2. Generate Mocks
```bash
# モックを生成
mockery --all --output=./internal/mocks --case=underscore
```

### 3. Database Migration
```bash
# マイグレーションを実行（初回のみ）
go run cmd/api-server/main.go migrate
```

### 4. Run Application
```bash
# Development mode
go run cmd/api-server/main.go

# または Makefileを使用
make run
```

サーバーは`http://localhost:8080`で起動します。

## API Testing

### Health Check
```bash
curl http://localhost:8080/health
```

期待されるレスポンス：
```json
{
  "status": "ok",
  "timestamp": "2023-12-01T12:00:00Z"
}
```

### User CRUD Operations

#### Create User
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "田中太郎",
    "email": "tanaka@example.com"
  }'
```

期待されるレスポンス：
```json
{
  "id": 1,
  "name": "田中太郎",
  "email": "tanaka@example.com",
  "created_at": "2023-12-01T12:00:00Z",
  "updated_at": "2023-12-01T12:00:00Z"
}
```

#### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
```

#### Get Users List
```bash
curl http://localhost:8080/api/v1/users?limit=10&offset=0
```

#### Update User
```bash
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "田中次郎",
    "email": "tanaka.jiro@example.com"
  }'
```

#### Delete User
```bash
curl -X DELETE http://localhost:8080/api/v1/users/1
```

## Testing

### Unit Tests
```bash
# 全てのテストを実行
go test ./...

# カバレッジ付きでテストを実行
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests
```bash
# 統合テストを実行（テスト用データベースが必要）
go test -tags=integration ./test/integration/...
```

## Development Tools

### Code Formatting
```bash
# コードフォーマット
go fmt ./...

# インポートを整理
goimports -w .
```

### Linting
```bash
# golangci-lintのインストール
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Lintを実行
golangci-lint run
```

### Database Operations

#### Reset Database
```bash
# データベースをリセット（開発環境のみ）
go run cmd/api-server/main.go migrate-reset
```

#### Create Migration
```bash
# 新しいマイグレーションファイルを作成
go run cmd/api-server/main.go create-migration <migration_name>
```

## Makefile Commands

プロジェクトルートの`Makefile`で利用可能なコマンド：

```bash
make build          # アプリケーションをビルド
make run             # アプリケーションを実行
make test            # テストを実行
make test-coverage   # カバレッジ付きテスト
make lint            # Lintを実行
make fmt             # コードフォーマット
make clean           # ビルド成果物を削除
make docker-up       # Docker環境を起動
make docker-down     # Docker環境を停止
make migrate         # データベースマイグレーション
make mock-gen        # モック生成
```

## Troubleshooting

### Common Issues

#### Database Connection Error
```
Error: failed to connect to database
```
解決方法：
1. PostgreSQLサービスが起動していることを確認
2. `.env`ファイルのデータベース設定を確認
3. データベースとユーザーが存在することを確認

#### Port Already in Use
```
Error: bind: address already in use
```
解決方法：
1. 別のプロセスがポート8080を使用していないか確認
2. `.env`ファイルで別のポートを指定

#### Mock Generation Failed
```
Error: mockery command not found
```
解決方法：
```bash
go install github.com/vektra/mockery/v2@latest
```

### Debug Mode
デバッグモードでアプリケーションを実行：

```bash
# 詳細ログを有効にして実行
LOG_LEVEL=debug go run cmd/api-server/main.go
```

## Project Structure Overview

```
├── cmd/api-server/           # アプリケーションエントリーポイント
│   └── internal/             # アプリケーション固有のコード
├── internal/                 # プライベートコード
│   ├── domain/              # ドメイン層
│   └── infrastructure/     # インフラストラクチャ層
├── test/                    # テストファイル
├── api/                     # API仕様
└── docs/                    # ドキュメント
```

各ディレクトリの詳細については、プロジェクト内の個別のREADMEファイルを参照してください。