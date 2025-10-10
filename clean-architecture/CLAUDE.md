# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 表示言語
ユーザーとのやり取りには日本語を利用してください
また、文書を書く際の日本語は文字化けを避けるためUTF-8エンコーディングを利用してください

## プロジェクト概要

このプロジェクトは、Clean Architectureの原則に基づいて構築されたシンプルなGo言語のRESTful APIサーバーです。

### アーキテクチャ概要

3層のシンプルな構成を採用しています：
- **API層**: HTTPリクエストの処理とビジネスロジックの実行
- **Domain層**: エンティティとリポジトリインターフェースの定義
- **Infrastructure層**: データベースアクセスなど外部依存の実装

### レイヤー間依存関係
@.claude/memories/dependency.md

### 命名規則
@.claude/memories/naming.md

### ディレクトリ構成
@.claude/memories/organization.md

### テスト方針
@.claude/memories/testing.md

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

# CI パイプライン
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

## 調査と探索
既存コードの調査を行う際にはSerena MCPを使って効率的に行ってください
