# Tasks: Go API Server with User Management

**Input**: `/specs/001-go-api-api/`の設計ドキュメント
**Tech Stack**: Go 1.21+, Gin HTTPフレームワーク, GORM v2, PostgreSQL, Mockery v3
**Architecture**: ドメイン駆動設計原則による階層アーキテクチャ

## Execution Flow Summary
research.md、data-model.md、contracts/user-api.yaml、quickstart.mdに基づく実装フロー：
- Go標準プロジェクトレイアウトで階層アーキテクチャをセットアップ
- 全APIコントラクトのテストを最初に生成（TDD）
- Userエンティティとリポジトリパターンを実装
- 適切なエラーハンドリングでAPIエンドポイントを構築
- 統合とポリッシュフェーズ

## Phase 3.1: Setup
- [x] T001 Go標準プロジェクトレイアウトに従った階層アーキテクチャでプロジェクト構造を作成
- [x] T002 go.modファイルを初期化し依存関係をインストール（Gin、GORM、PostgreSQLドライバー、Mockery v3）
- [x] T003 [P] golangci-lintとgoフォーマットツールの設定
- [x] T004 [P] cmd/api-server/internal/config/config.goで環境変数によるPostgreSQL接続設定
- [x] T005 [P] ビルド、テスト、lint、モック生成コマンドを含むMakefileの作成

## Phase 3.2: Tests First (TDD) ⚠️ 3.3の前に必ず完了
**重要: これらのテストは実装前に必ず作成し、失敗する必要があります**
- [x] T006 [P] api/test/contract/users_post_test.goでPOST /api/v1/usersのコントラクトテスト
- [x] T007 [P] api/test/contract/users_get_list_test.goでGET /api/v1/usersのコントラクトテスト
- [x] T008 [P] api/test/contract/users_get_by_id_test.goでGET /api/v1/users/{id}のコントラクトテスト
- [x] T009 [P] api/test/contract/users_put_test.goでPUT /api/v1/users/{id}のコントラクトテスト
- [x] T010 [P] api/test/contract/users_delete_test.goでDELETE /api/v1/users/{id}のコントラクトテスト
- [x] T011 [P] api/test/integration/user_registration_test.goでユーザー登録フローの統合テスト
- [x] T012 [P] api/test/integration/user_crud_test.goでデータベースを使用したユーザーCRUD操作の統合テスト

## Phase 3.3: Core Implementation (テストが失敗した後のみ)
- [x] T013 [P] internal/domain/entities/user.goでUserエンティティ構造体
- [x] T014 [P] internal/domain/repositories/user_repository.goでUserリポジトリインターフェース
- [x] T015 [P] internal/domain/services/user_domain_service.goでUserドメインサービス
- [x] T016 [P] internal/infrastructure/repositories/user_repository_impl.goでGORMを使用したUserリポジトリ実装
- [x] T017 [P] cmd/api-server/internal/application/user_service.goでビジネスロジックを含むUserアプリケーションサービス
- [x] T018 [P] cmd/api-server/internal/controller/dtos/user_dto.goでUserのDTOとリクエスト/レスポンスモデル
- [x] T019 cmd/api-server/internal/controller/user_controller.goでPOST /api/v1/usersエンドポイントハンドラー
- [x] T020 GET /api/v1/usersエンドポイントハンドラー（user_controller.goを拡張）
- [x] T021 GET /api/v1/users/{id}エンドポイントハンドラー（user_controller.goを拡張）
- [x] T022 PUT /api/v1/users/{id}エンドポイントハンドラー（user_controller.goを拡張）
- [x] T023 DELETE /api/v1/users/{id}エンドポイントハンドラー（user_controller.goを拡張）
- [x] T024 cmd/api-server/internal/middleware/validation.goで入力検証ミドルウェア
- [x] T025 cmd/api-server/internal/middleware/error_handler.goでエラーハンドリングミドルウェア

## Phase 3.4: Integration
- [x] T026 internal/infrastructure/database/migrations/migration.goでGORMオートマイグレーションを使用したデータベースマイグレーション設定
- [x] T027 internal/infrastructure/database/postgres/connection.goでPostgreSQL接続設定
- [x] T028 cmd/api-server/internal/container/container.goで依存性注入コンテナ
- [x] T029 cmd/api-server/internal/router/router.goでGinを使用したHTTPルーター設定
- [x] T030 cmd/api-server/internal/middleware/logger.goで構造化ログミドルウェア
- [x] T031 cmd/api-server/internal/controller/health_controller.goでヘルスチェックエンドポイント
- [x] T032 cmd/api-server/main.goでグレースフルシャットダウン処理

## Phase 3.5: Polish
- [x] T033 [P] internal/domain/entities/user_test.goでUserエンティティバリデーションのユニットテスト
- [x] T034 [P] cmd/api-server/internal/application/user_service_test.goでUserサービスビジネスロジックのユニットテスト
- [x] T035 [P] cmd/api-server/internal/controller/user_controller_test.goでUserコントローラーHTTPハンドラーのユニットテスト
- [x] T036 [P] scripts/generate-mocks.shでMockery v3を使用したリポジトリインターフェースのモック生成
- [x] T037 api/test/performance/user_api_performance_test.goでエンドポイントが200ms以内に応答することを確認するパフォーマンステスト
- [x] T038 [P] internal/infrastructure/database/postgres/connection.goでデータベース接続プールの最適化
- [x] T039 [P] docs/api.mdでOpenAPI仕様からのAPIドキュメント生成
- [x] T040 scripts/coverage.shでコードカバレッジ分析とレポート設定
- [x] T041 quickstart.mdのシナリオに従った手動テスト

## Dependencies
- Setup（T001-T005）が全ての前提
- Tests（T006-T012）がImplementation（T013-T025）の前提
- T013（Userエンティティ）がT014、T015、T016、T017をブロック
- T014（リポジトリインターフェース）がT016をブロック
- T016（リポジトリ実装）がT017をブロック
- T017（Userサービス）がT019-T023をブロック
- T018（DTOs）がT019-T023をブロック
- T019-T023（エンドポイント）がT024、T025をブロック
- Implementation（T013-T025）がIntegration（T026-T032）の前提
- IntegrationがPolish（T033-T041）の前提

## Parallel Execution Examples

### Phase 3.2 - Test Creation (全て並列)
```bash
# T006-T012を同時に実行（異なるファイル、独立）:
Task: "api/test/contract/users_post_test.goでPOST /api/v1/usersのコントラクトテスト"
Task: "api/test/contract/users_get_list_test.goでGET /api/v1/usersのコントラクトテスト"
Task: "api/test/contract/users_get_by_id_test.goでGET /api/v1/users/{id}のコントラクトテスト"
Task: "api/test/contract/users_put_test.goでPUT /api/v1/users/{id}のコントラクトテスト"
Task: "api/test/contract/users_delete_test.goでDELETE /api/v1/users/{id}のコントラクトテスト"
Task: "api/test/integration/user_registration_test.goでユーザー登録フローの統合テスト"
Task: "api/test/integration/user_crud_test.goでデータベースを使用したユーザーCRUD操作の統合テスト"
```

### Phase 3.3 - Core Implementation (可能な限り並列)
```bash
# T013-T018を同時に実行（異なるファイル、依存関係なし）:
Task: "internal/domain/entities/user.goでUserエンティティ構造体"
Task: "internal/domain/repositories/user_repository.goでUserリポジトリインターフェース"
Task: "internal/domain/services/user_domain_service.goでUserドメインサービス"
Task: "internal/infrastructure/repositories/user_repository_impl.goでGORMを使用したUserリポジトリ実装"
Task: "cmd/api-server/internal/application/user_service.goでビジネスロジックを含むUserアプリケーションサービス"
Task: "cmd/api-server/internal/controller/dtos/user_dto.goでUserのDTOとリクエスト/レスポンスモデル"
```

### Phase 3.5 - Polish (テストファイル並列)
```bash
# T033-T035を同時に実行（異なるテストファイル）:
Task: "internal/domain/entities/user_test.goでUserエンティティバリデーションのユニットテスト"
Task: "cmd/api-server/internal/application/user_service_test.goでUserサービスビジネスロジックのユニットテスト"
Task: "cmd/api-server/internal/controller/user_controller_test.goでUserコントローラーHTTPハンドラーのユニットテスト"
```

## Project Structure (Plan.mdに基づく)
```
├── cmd/
│   └── api-server/           # APIサーバーエントリーポイント
│       ├── main.go          # アプリケーションエントリーポイント
│       └── internal/        # このアプリケーション専用のコード
│           ├── controller/  # HTTPハンドラー（コントローラー層）
│           ├── application/ # ユースケース（アプリケーション層）
│           ├── middleware/  # ミドルウェア
│           └── config/      # 設定管理
├── internal/                # プライベートなアプリケーションとライブラリコード
│   ├── domain/             # ドメイン層
│   │   ├── entities/       # エンティティ定義（User）
│   │   ├── repositories/   # リポジトリインターフェース
│   │   └── services/       # ドメインサービス
│   └── infrastructure/     # インフラストラクチャ層
│       ├── database/       # データベース実装
│       │   ├── postgres/   # PostgreSQL固有の実装
│       │   └── migrations/ # マイグレーションファイル
│       └── repositories/   # リポジトリ実装
├── api/                   # OpenAPI/Swagger仕様、JSONスキーマファイル、テスト
│   └── test/              # APIテストファイル
│       ├── contract/      # APIコントラクトテスト
│       ├── integration/   # 統合テスト
│       └── performance/   # パフォーマンステスト
├── docs/                  # 設計とユーザードキュメント
├── scripts/               # ビルド、インストール、解析等を行うスクリプト
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## Validation Checklist
- [x] 全APIコントラクト（5エンドポイント）に対応するテスト（T006-T010）
- [x] Userエンティティにモデル作成タスク（T013）
- [x] 全テストが実装前（T006-T012がT013-T025の前）
- [x] 並列タスクが真に独立（異なるファイル、共有依存関係なし）
- [x] 各タスクが正確なファイルパスを指定
- [x] 同じファイルを変更する[P]タスクなし
- [x] TDDアプローチ：失敗するテストから、その後実装
- [x] Go標準プロジェクトレイアウトがタスク全体で適切に実装
- [x] quickstart.mdの全シナリオが統合テストでカバー

## Notes
- [P]タスクは並列実行可能（異なるファイル、依存関係なし）
- 実装前にテストが失敗することを確認（TDD要件）
- 各タスク完了後にコミット
- Go標準プロジェクトレイアウトと命名規則に従う
- インターフェースモック生成にMockery v3を使用
- 全エンドポイントはOpenAPI仕様に従って適切なHTTPステータスコードを処理
- data-model.mdで指定されたUserエンティティのソフトデリートを実装