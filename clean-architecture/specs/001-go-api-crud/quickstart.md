# クイックスタートガイド: Go API サーバー

**対象**: 開発者向け環境セットアップとAPIテスト手順
**前提条件**: Go 1.24+、Docker、PostgreSQL

## 📋 前提条件

### 必要ソフトウェア
- **Go**: バージョン 1.24 以上
- **Docker & Docker Compose**: PostgreSQL コンテナ実行用
- **curl または Postman**: API テスト用
- **make**: ビルドコマンド実行用

### 確認コマンド
```bash
go version    # Go 1.24+ であることを確認
docker --version
docker compose version
make --version
```

## 🚀 環境セットアップ

### 1. リポジトリ設定
```bash
# プロジェクトルートに移動
cd go-api-server-sample

# 依存関係の確認
go mod tidy
```

### 2. 環境ファイル作成
```bash
# .env ファイルを作成（サンプルをコピー）
cp .env.example .env

# 環境変数の内容確認・編集
cat .env
```

**推奨環境変数**:
```bash
# データベース設定
DB_HOST=localhost
DB_PORT=5432
DB_USER=api_user
DB_PASSWORD=api_password
DB_NAME=api_db

# サーバー設定
PORT=8080
GIN_MODE=debug

# ログ設定
LOG_LEVEL=debug
```

### 3. PostgreSQL セットアップ
```bash
# PostgreSQL コンテナ起動
make docker-up

# コンテナ起動確認
docker ps | grep postgres

# データベース接続テスト
make db-ping
```

### 4. データベース初期化
```bash
# マイグレーション実行
make migrate

# 初期データ投入（オプション）
make seed

# テーブル作成確認
make db-status
```

## 🏃‍♂️ アプリケーション起動

### 開発モード（ホットリロード）
```bash
# air を使った自動リロード起動
make dev

# または直接実行
air
```

### 通常実行
```bash
# ビルドして実行
make run

# または手動ビルド
make build
./bin/api-server
```

### 起動確認
```bash
# ヘルスチェック
curl http://localhost:8080/health

# 期待レスポンス
{
  "status": "healthy",
  "database": "connected",
  "timestamp": "2025-09-28T15:30:00Z"
}
```

## 🧪 API テストシナリオ

### 1. ヘルスチェック
```bash
curl -X GET http://localhost:8080/health
```

**期待結果**: 200 OK, システム健全性情報

### 2. コンテンツ作成
```bash
curl -X POST http://localhost:8080/api/v1/contents \
  -H "Content-Type: application/json" \
  -d '{
    "title": "テスト記事",
    "body": "これはテスト用の記事です。",
    "content_type": "article",
    "author": "テストユーザー"
  }'
```

**期待結果**: 201 Created, 作成されたコンテンツ情報（IDを含む）

### 3. コンテンツ一覧取得
```bash
curl -X GET http://localhost:8080/api/v1/contents
```

**期待結果**: 200 OK, コンテンツ一覧（作成したコンテンツを含む）

### 4. コンテンツ詳細取得
```bash
# ID=1 のコンテンツを取得
curl -X GET http://localhost:8080/api/v1/contents/1
```

**期待結果**: 200 OK, 指定IDのコンテンツ詳細

### 5. コンテンツ更新
```bash
curl -X PUT http://localhost:8080/api/v1/contents/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "更新されたテスト記事",
    "body": "これは更新されたテスト用の記事です。",
    "content_type": "blog",
    "author": "更新者"
  }'
```

**期待結果**: 200 OK, 更新されたコンテンツ情報

### 6. コンテンツ削除
```bash
curl -X DELETE http://localhost:8080/api/v1/contents/1
```

**期待結果**: 204 No Content

### 7. 削除確認
```bash
curl -X GET http://localhost:8080/api/v1/contents/1
```

**期待結果**: 404 Not Found

## 🔍 フィルタリング・ページネーション テスト

### コンテンツタイプフィルタ
```bash
curl -X GET "http://localhost:8080/api/v1/contents?content_type=article"
```

### 作成者フィルタ
```bash
curl -X GET "http://localhost:8080/api/v1/contents?author=テストユーザー"
```

### ページネーション
```bash
# 最初の5件を取得
curl -X GET "http://localhost:8080/api/v1/contents?limit=5&offset=0"

# 次の5件を取得
curl -X GET "http://localhost:8080/api/v1/contents?limit=5&offset=5"
```

## 🛠 開発ツール

### テスト実行
```bash
# 全テスト実行
make test

# カバレッジ付きテスト
make test-coverage

# 統合テスト
make test-integration

# 特定パッケージのテスト
go test -v ./internal/domain/entities
```

### コード品質チェック
```bash
# リント実行
make lint

# コードフォーマット
make fmt

# 型チェック
make vet

# 全チェック
make check
```

### モック生成
```bash
# モック生成
make mock-gen

# 生成されたモックの確認
ls internal/testing/mocks/
```

## 📊 パフォーマンステスト

### 基本負荷テスト
```bash
# ヘルスチェックエンドポイント
ab -n 1000 -c 10 http://localhost:8080/health

# コンテンツ一覧取得
ab -n 1000 -c 10 http://localhost:8080/api/v1/contents
```

**期待値**: レスポンス時間 < 200ms p95

### ベンチマークテスト
```bash
make test-performance
```

## 🐛 トラブルシューティング

### よくある問題と解決方法

#### 1. PostgreSQL 接続エラー
```bash
# エラー例: connection refused
# 解決: PostgreSQL コンテナの状態確認
docker ps
make docker-up

# ポート競合の場合
docker port postgres_container
```

#### 2. マイグレーションエラー
```bash
# テーブル作成失敗
# 解決: データベースリセット
make migrate-reset
make migrate
```

#### 3. ポート競合
```bash
# ポート 8080 が使用中
# 解決: ポート確認と変更
lsof -i :8080
export PORT=8081
make run
```

#### 4. 依存関係エラー
```bash
# モジュール解決エラー
# 解決: 依存関係の再取得
go clean -modcache
go mod download
go mod tidy
```

#### 5. テスト失敗
```bash
# testcontainers エラー
# 解決: Docker デーモン確認
systemctl status docker
docker info
```

## 📝 ログ確認

### アプリケーションログ
```bash
# 開発モードでのログ監視
tail -f logs/app.log

# リアルタイムログ（JSON形式）
make dev | jq '.'
```

### データベースログ
```bash
# PostgreSQL コンテナログ
docker logs postgres_container -f
```

## 🔧 開発設定

### VSCode 設定推奨
```json
{
  "go.formatTool": "goimports",
  "go.lintTool": "golangci-lint",
  "go.testFlags": ["-v"],
  "go.coverOnSave": true
}
```

### デバッグ設定
```json
{
  "type": "go",
  "request": "launch",
  "mode": "debug",
  "program": "./cmd/api-server",
  "env": {
    "GIN_MODE": "debug"
  }
}
```

## ✅ 環境準備完了チェックリスト

- [ ] Go 1.24+ インストール完了
- [ ] Docker & Docker Compose 起動確認
- [ ] PostgreSQL コンテナ起動完了
- [ ] 環境変数設定完了
- [ ] マイグレーション実行完了
- [ ] アプリケーション起動成功
- [ ] ヘルスチェック API 応答確認
- [ ] CRUD API 基本動作確認
- [ ] テスト実行成功
- [ ] リント・フォーマット確認

## 📚 次のステップ

1. **実装フェーズ**: `/tasks` コマンドでタスク一覧生成
2. **TDD サイクル**: Red-Green-Refactor の実践
3. **CI/CD 設定**: GitHub Actions パイプライン構築
4. **監視設定**: メトリクス・ログ監視の導入

## 💡 開発のヒント

- **コミット前に必ず `make check` を実行**
- **Clean Architecture 原則を意識した実装**
- **テスト駆動開発の実践**
- **適切なエラーハンドリングの実装**
- **セキュリティベストプラクティスの遵守**