# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## プロジェクト概要

このプロジェクトは、Clean ArchitectureとDomain-Driven Designの原則に基づいて構築されたGo言語のRESTful APIサーバーです。

### アーキテクチャ
@.serena/memories/architecture-guidelines.md

### ディレクトリ構造

```
├── cmd/api-server/           # アプリケーションエントリーポイント
│   └── internal/             # アプリケーション層
│       ├── application/      # UseCase層（1ユースケース1ファイル）
│       ├── config/           # 設定管理
│       ├── container/        # 依存性注入コンテナ
│       ├── controller/       # HTTPコントローラ
│       ├── middleware/       # HTTPミドルウェア
│       └── router/           # ルーティング設定
├── internal/                 # 内部パッケージ
│   ├── domain/              # ドメイン層
│   │   ├── entities/        # ドメインエンティティ
│   │   ├── repositories/    # リポジトリインターフェース
│   │   └── services/        # ドメインサービス
│   └── infrastructure/      # インフラストラクチャ層
│       ├── database/        # データベース関連
│       └── repositories/    # リポジトリ実装
├── api/test/                # APIテスト
│   ├── contract/            # 契約テスト
│   ├── integration/         # 統合テスト
│   └── performance/         # パフォーマンステスト
└── pkg/                     # 共有パッケージ
```

### 技術スタック

- **言語**: Go 1.23
- **Webフレームワーク**: Gin
- **ORM**: GORM
- **データベース**: PostgreSQL
- **テスト**: testify, mockery
- **リント**: golangci-lint
- **開発ツール**: air (ホットリロード)

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

# Linux向けビルド
make build-linux
```

### テスト

```bash
# 全テスト実行
make test

# カバレッジ付きテスト
make test-coverage

# 統合テスト
make test-integration

# 契約テスト
make test-contract

# パフォーマンステスト
make test-performance

# 全種類のテスト実行
make test-all
```

#### テストガイドライン
@.claude/memories/testing-guidelines.md

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

# CI パイプライン
make ci
```

### データベース

```bash
# マイグレーション実行
make migrate

# データベースリセット（開発環境のみ）
make migrate-reset

# テストデータ投入
make seed
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

# モック生成
make mock-gen
```

## 単体テスト実行

```bash
# 特定のパッケージのテスト
go test -v ./internal/domain/entities

# 特定のテスト関数実行
go test -v -run TestUserValidation ./internal/domain/entities

# ベンチマークテスト
go test -bench=. ./internal/domain/entities
```

## データベース設定

### 開発環境セットアップ

```bash
# 環境ファイル作成とDocker起動、マイグレーション実行
make quickstart
```

## API エンドポイント

### ユーザー管理

- `POST /api/v1/users` - ユーザー作成
- `GET /api/v1/users` - ユーザー一覧取得
- `GET /api/v1/users/:id` - ユーザー詳細取得
- `PUT /api/v1/users/:id` - ユーザー更新
- `DELETE /api/v1/users/:id` - ユーザー削除

### ヘルスチェック

- `GET /health` - アプリケーション状態確認

## ファイル構成・命名規約
@.serena/memories/organization.md

## コーディング規約

### Go言語標準に従った命名規則

- パッケージ名：小文字、短く
- 関数名：CamelCase（public）、camelCase（private）
- 構造体：CamelCase
- 定数：CamelCase
- ファイル名：ケバブケース（ハイフン区切り）、1ユースケース1ファイル

### エラーハンドリング

- ドメインエラーは `internal/domain/entities` で定義
- インフラエラーはアプリケーション層で適切に変換
- HTTPエラーレスポンスは統一された形式を使用

### 依存性注入

- `container` パッケージで一元管理（`cmd/api-server/internal/container/container.go`）
- インターフェースを活用した疎結合設計
- 各ユースケースファイル内に必要な依存関係のインターフェースを定義

## リント設定

`.golangci.yml` で以下を設定済み：

- 行長制限：140文字
- 複雑度制限：15
- 関数長制限：100行、50ステートメント
- セキュリティチェック有効
- テストファイルは一部ルールを除外

## デバッグ

### ログレベル

開発環境では `LOG_LEVEL=debug` でデバッグ情報を出力

### ホットリロード

```bash
# airを使用した自動リロード
make dev
```

### データベースデバッグ

```bash
# PostgreSQLコンテナに接続
docker exec -it go-api-server-postgres psql -U postgres -d go_api_server
```

## 追加ガイドライン

### Claude Code操作ガイドライン (`.claude/memories/`)
- `architecture-rules.md` - アーキテクチャルール・依存関係
- `naming-conventions.md` - 命名規約・ファイル構成
- `testing-guidelines.md` - テスト駆動開発ガイドライン
- `git-workflow.md` - Git ワークフロー
- `release-process.md` - リリースプロセス
- `performance-guidelines.md` - パフォーマンスガイドライン
- `security-guidelines.md` - セキュリティガイドライン
- `coding.md` - コーディング規約

### プロジェクト構成情報 (`.serena/memories/`)
- `architecture-guidelines.md` - 4層アーキテクチャ構成
- `organization.md` - ディレクトリ構造・ファイル配置
- `api-design.md` - API設計ガイドライン
- `database-conventions.md` - データベース規約
- `environment-config.md` - 環境設定
- `logging-monitoring.md` - ログとモニタリング
- `test-coverage.md` - テストカバレッジ要件
