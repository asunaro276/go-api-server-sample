#!/bin/bash

# quickstart.md のテストシナリオ実行検証スクリプト

set -e

echo "🚀 Go API サーバーのクイックスタート検証を開始します"

# 1. ビルド確認
echo "📦 1. アプリケーションビルド確認"
if go build -o bin/api-server ./cmd/api-server; then
    echo "✅ ビルド成功"
else
    echo "❌ ビルド失敗"
    exit 1
fi

# 2. 基本テスト実行
echo "🧪 2. 基本テスト実行"
if go test ./internal/domain/entities -v > /dev/null 2>&1; then
    echo "✅ エンティティテスト成功"
else
    echo "❌ エンティティテスト失敗"
fi

# 3. コードフォーマット確認
echo "📝 3. コードフォーマット確認"
if go fmt ./...; then
    echo "✅ コードフォーマット完了"
else
    echo "❌ コードフォーマット失敗"
fi

# 4. go vet確認
echo "🔍 4. Go vet確認"
if go vet ./...; then
    echo "✅ Go vet成功"
else
    echo "❌ Go vet警告またはエラー"
fi

# 5. 依存関係確認
echo "📚 5. 依存関係確認"
if go mod tidy && go mod verify; then
    echo "✅ 依存関係確認成功"
else
    echo "❌ 依存関係エラー"
fi

# 6. ディレクトリ構造確認
echo "🏗️ 6. Clean Architectureディレクトリ構造確認"
required_dirs=(
    "cmd/api-server/internal/application"
    "cmd/api-server/internal/controller"
    "cmd/api-server/internal/middleware"
    "internal/domain/entities"
    "internal/domain/repositories"
    "internal/infrastructure/database"
    "internal/infrastructure/repositories"
    "config"
)

for dir in "${required_dirs[@]}"; do
    if [ -d "$dir" ]; then
        echo "✅ $dir 存在"
    else
        echo "❌ $dir が見つかりません"
    fi
done

# 7. 重要ファイル存在確認
echo "📄 7. 重要ファイル存在確認"
required_files=(
    "go.mod"
    "Makefile"
    ".air.toml"
    ".golangci.yml"
    ".env.example"
    "cmd/api-server/main.go"
    "internal/domain/entities/content.go"
    "internal/infrastructure/database/connection.go"
)

for file in "${required_files[@]}"; do
    if [ -f "$file" ]; then
        echo "✅ $file 存在"
    else
        echo "❌ $file が見つかりません"
    fi
done

echo "🎉 クイックスタート検証完了！"
echo ""
echo "次のステップ:"
echo "1. PostgreSQLコンテナを起動: make docker-up"
echo "2. マイグレーション実行: make migrate"
echo "3. 開発サーバー起動: make dev"
echo "4. ヘルスチェック確認: curl http://localhost:8080/health"