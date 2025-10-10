# ファイル構成・命名規約

## ディレクトリ構成

```
├── cmd/api-server/           # アプリケーションエントリーポイント
│   ├── main.go              # エントリーポイント
│   ├── internal/            # API層とインフラストラクチャ設定
│   │   ├── api/            # API層（HTTPハンドラー）
│   │   │   ├── content/    # コンテンツ関連API
│   │   │   └── health/     # ヘルスチェックAPI
│   │   ├── container/      # 依存性注入コンテナ
│   │   └── middleware/     # HTTPミドルウェア
│   └── test/               # APIテスト
│       ├── integration/    # 統合テスト
│       └── performance/    # パフォーマンステスト
├── internal/               # 内部パッケージ（Domain層とInfrastructure層）
│   ├── domain/            # ドメイン層
│   │   ├── entities/      # ドメインエンティティ
│   │   └── repositories/  # リポジトリインターフェース
│   └── infrastructure/    # インフラストラクチャ層
│       ├── database/      # データベース接続管理
│       └── repositories/  # リポジトリ実装
├── config/                # 設定管理
├── scripts/               # 開発・運用スクリプト
└── specs/                 # 仕様書・設計ドキュメント
```

## レイヤー別ディレクトリ詳細

### API層 (`cmd/api-server/internal/api/`)

各リソースごとにサブディレクトリを作成し、CRUDなど関連する操作をまとめます。

```
api/
├── content/
│   ├── content.go    # API構造体定義
│   ├── create.go     # 作成処理
│   ├── get.go        # 取得処理
│   ├── list.go       # 一覧取得処理
│   ├── update.go     # 更新処理
│   └── delete.go     # 削除処理
└── health/
    └── check.go      # ヘルスチェック
```

### Domain層 (`internal/domain/`)

ドメインロジックとインターフェース定義を配置します。

```
domain/
├── entities/
│   ├── content.go        # エンティティ定義
│   └── content_test.go   # エンティティテスト
└── repositories/
    └── content.go        # リポジトリインターフェース
```

### Infrastructure層 (`internal/infrastructure/`)

外部依存の実装を配置します。

```
infrastructure/
├── database/
│   ├── connection.go     # DB接続管理
│   └── migrate.go        # マイグレーション
└── repositories/
    ├── content.go        # リポジトリ実装
    └── content_test.go   # リポジトリテスト
```

## ファイル配置ルール

### 基本原則

1. **API層**: リソースごとに1ディレクトリ、操作ごとに1ファイル
2. **Domain層**: エンティティとリポジトリインターフェースを分離
3. **Infrastructure層**: リポジトリ実装とDB接続を分離
4. **テストファイル**: 対象ファイルと同一ディレクトリに配置

### パッケージ構成の考え方

- **cmd/api-server/internal**: アプリケーション固有のコード
- **internal**: プロジェクト内部で共有されるコード（外部公開不可）
- **config**: 環境設定やアプリケーション設定
- **scripts**: 開発・運用で使用するスクリプト
- **specs**: プロジェクトの仕様書や設計ドキュメント
