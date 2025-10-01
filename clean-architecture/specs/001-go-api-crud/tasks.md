# タスク: Go API CRUD

**入力**: `/specs/001-go-api-crud/` からの設計ドキュメント
**前提条件**: plan.md (必須), research.md, data-model.md, contracts/

## 実行フロー (main)
```
1. 機能ディレクトリからplan.mdを読み込み
   → 見つからない場合: ERROR "実装計画が見つかりません"
   → 抽出: 技術スタック、ライブラリ、構成
2. オプションの設計ドキュメントを読み込み:
   → data-model.md: エンティティを抽出 → モデルタスク
   → contracts/: 各ファイル → 契約テストタスク
   → research.md: 決定事項を抽出 → セットアップタスク
3. カテゴリ別にタスクを生成:
   → セットアップ: プロジェクト初期化、依存関係、リント
   → テスト: 契約テスト、統合テスト
   → コア: モデル、サービス、エンドポイント
   → 統合: DB、ミドルウェア、ログ
   → 仕上げ: 単体テスト、パフォーマンス、ドキュメント
4. タスクルールを適用:
   → 異なるファイル = 並列用に[P]をマーク
   → 同じファイル = 順次（[P]なし）
   → 実装前にテスト（TDD）
5. タスクを順次番号付け（T001, T002...）
6. 依存関係グラフを生成
7. 並列実行例を作成
8. タスク完全性を検証:
   → すべての契約にテストがあるか？
   → すべてのエンティティにモデルがあるか？
   → すべてのエンドポイントが実装されているか？
9. 返却: SUCCESS (タスクは実行準備完了)
```

## 形式: `[ID] [P?] 説明`
- **[P]**: 並列実行可能（異なるファイル、依存関係なし）
- 説明に正確なファイルパスを含める

## パス規約
- **Goプロジェクト**: 既存のClean Architecture構造を活用
- cmd/api-server/internal/ (UseCase層)
- internal/domain/ (Domain層)
- internal/infrastructure/ (Infrastructure層)

## フェーズ 3.1: セットアップ
- [ ] T001 Clean Architecture構造に従ってプロジェクトディレクトリを整理
- [ ] T002 Go 1.24プロジェクトをGin+GORM+testify依存関係で初期化
- [ ] T003 [P] golangci-lintとairツールを設定

## フェーズ 3.2: テストファースト（TDD）⚠️ 3.3前に必須完了
**重要: これらのテストは実装前に記述し、失敗しなければならない**
- [X] T004 [P] cmd/api-server/main_test.go で各APIの契約テスト（正常系のみ）
- [X] T010 [P] api/test/integration/content_crud_test.go でコンテンツCRUD統合テスト
- [X] T011 [P] api/test/integration/health_check_test.go でヘルスチェック統合テスト

## フェーズ 3.3: コア実装（テスト失敗後のみ）
- [X] T012 [P] internal/domain/entities/content.go でContentエンティティ
- [X] T013 [P] internal/domain/repositories/content.go でContentRepositoryインターフェース
- [X] T014 [P] cmd/api-server/internal/application/createcontent.go でCreateContentユースケース
- [X] T015 [P] cmd/api-server/internal/application/getcontent.go でGetContentユースケース
- [X] T016 [P] cmd/api-server/internal/application/listcontents.go でListContentsユースケース
- [X] T017 [P] cmd/api-server/internal/application/updatecontent.go でUpdateContentユースケース
- [X] T018 [P] cmd/api-server/internal/application/deletecontent.go でDeleteContentユースケース
- [X] T019 [P] cmd/api-server/internal/application/healthcheck.go でHealthCheckユースケース
- [X] T020 internal/infrastructure/repositories/content.go でContentRepository実装
- [X] T021 cmd/api-server/internal/controller/content.go でContentController
- [X] T022 cmd/api-server/internal/controller/health.go でHealthController
- [X] T023 cmd/api-server/internal/controller/dto/content.go でリクエスト/レスポンスDTO
- [X] T024 cmd/api-server/main.go でメインエントリーポイント

## フェーズ 3.4: 統合
- [X] T025 internal/infrastructure/database/connection.go でPostgreSQL接続
- [X] T026 internal/infrastructure/database/migrate.go でマイグレーション
- [X] T027 [P] cmd/api-server/internal/middleware/error.go でエラーハンドリングミドルウェア
- [X] T028 [P] cmd/api-server/internal/middleware/cors.go でCORSミドルウェア
- [X] T029 [P] cmd/api-server/internal/middleware/logging.go でログミドルウェア
- [X] T030 cmd/api-server/internal/container/container.go で依存性注入コンテナ
- [X] T031 config/config.go で設定管理

## フェーズ 3.5: 仕上げ
- [X] T032 [P] internal/domain/entities/content_test.go でContentエンティティ単体テスト
- [X] T033 [P] internal/infrastructure/repositories/content_test.go でContentRepository単体テスト（testcontainers使用）
- [ ] T034 [P] cmd/api-server/internal/application/createcontent_test.go でCreateContentユースケース単体テスト
- [ ] T035 [P] cmd/api-server/internal/application/getcontent_test.go でGetContentユースケース単体テスト
- [ ] T036 [P] cmd/api-server/internal/application/listcontents_test.go でListContentsユースケース単体テスト
- [ ] T037 [P] cmd/api-server/internal/application/updatecontent_test.go でUpdateContentユースケース単体テスト
- [ ] T038 [P] cmd/api-server/internal/application/deletecontent_test.go でDeleteContentユースケース単体テスト
- [X] T039 [P] api/test/performance/api_performance_test.go でパフォーマンステスト（<200ms）
- [X] T040 重複削除とコードリファクタリング
- [X] T041 quickstart.mdのテストシナリオ実行検証

## 依存関係
- テスト（T004-T011）は実装前（T012-T024）
- T012はT013、T020をブロック
- T014-T018はT021をブロック
- T025はT026、T020をブロック
- T030はT024をブロック
- 実装は仕上げ前（T032-T041）

## 並列実行例
```
# T004-T011を一緒に起動:
Task: "cmd/api-server/main_test.go で各APIの契約テスト（正常系のみ）"
Task: "api/test/integration/content_crud_test.go でコンテンツCRUD統合テスト"
Task: "api/test/integration/health_check_test.go でヘルスチェック統合テスト"
```

```
# T012-T019を一緒に起動:
Task: "internal/domain/entities/content.go でContentエンティティ"
Task: "internal/domain/repositories/content.go でContentRepositoryインターフェース"
Task: "cmd/api-server/internal/application/createcontent.go でCreateContentユースケース"
Task: "cmd/api-server/internal/application/getcontent.go でGetContentユースケース"
Task: "cmd/api-server/internal/application/listcontents.go でListContentsユースケース"
Task: "cmd/api-server/internal/application/updatecontent.go でUpdateContentユースケース"
Task: "cmd/api-server/internal/application/deletecontent.go でDeleteContentユースケース"
Task: "cmd/api-server/internal/application/healthcheck.go でHealthCheckユースケース"
```

## 注意事項
- [P]タスク = 異なるファイル、依存関係なし
- 実装前にテストの失敗を確認
- 各タスク後にコミット
- 避ける: 曖昧なタスク、同一ファイル競合

## タスク生成ルール
*main()実行中に適用*

1. **契約から**:
   - 各契約エンドポイント → 契約テストタスク [P]
   - 各エンドポイント → 実装タスク

2. **データモデルから**:
   - Contentエンティティ → モデル作成タスク [P]
   - リポジトリ → Repository実装タスク

3. **ユーザーストーリーから**:
   - CRUD操作 → 統合テスト [P]
   - ヘルスチェック → 統合テスト [P]
   - quickstartシナリオ → 検証タスク

4. **順序**:
   - セットアップ → テスト → モデル → ユースケース → コントローラ → 統合 → 仕上げ
   - 依存関係は並列実行をブロック

## 検証チェックリスト
*ゲート: 返却前にmain()でチェック*

- [x] すべての契約に対応するテストがある（health + contents CRUD）
- [x] すべてのエンティティにモデルタスクがある（Content）
- [x] すべてのテストが実装前にある（T004-T011 → T012-T024）
- [x] 並列タスクが真に独立している（異なるファイル）
- [x] 各タスクが正確なファイルパスを指定している
- [x] 他の[P]タスクと同じファイルを変更するタスクがない
