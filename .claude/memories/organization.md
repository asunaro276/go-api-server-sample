# ファイル構成・命名規約

## ディレクトリ構成

```
├── cmd/api-server/           # アプリケーションエントリーポイント
│   └── internal/             # アプリケーション層
│       ├── application/      # UseCase層（アプリケーションサービス）
│       ├── container/        # 依存性注入コンテナ
│       ├── controller/       # HTTPコントローラ
│       └── middleware/       # HTTPミドルウェア
├── internal/                 # 内部パッケージ
│   ├── domain/              # ドメイン層
│   │   └── entities/        # ドメインエンティティ
│   └── infrastructure/      # インフラストラクチャ層
│       ├── database/        # データベース関連
│       └── repositories/    # リポジトリ実装
├── api/test/                # APIテスト
│   ├── integration/         # 統合テスト
│   └── performance/         # パフォーマンステスト
├── config/                  # 設定管理
└── pkg/                     # 共有パッケージ
```

## ファイル配置

プロジェクトでは以下の配置ルールに従っています：

- UseCase層: 1ユースケース1ファイルの原則
- テストファイル: 対象ファイルと同一ディレクトリに配置
- 設定・共通ファイル: 機能ごとに専用ディレクトリに配置
