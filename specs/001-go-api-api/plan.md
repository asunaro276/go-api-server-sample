
# Implementation Plan: Go API Server with Layered Architecture

**Branch**: `001-go-api-api` | **Date**: 2025-09-21 | **Spec**: [spec.md](./spec.md)
**Input**: 仕様書（`/specs/001-go-api-api/spec.md`）

## Execution Flow (/plan command scope)
```
1. 入力パスから機能仕様を読み込み
   → 見つからない場合: ERROR "No feature spec at {path}"
2. Technical Contextを記入（NEEDS CLARIFICATIONをスキャン）
   → コンテキストからProject Typeを検出（web=frontend+backend, mobile=app+api）
   → プロジェクトタイプに基づいてStructure Decisionを設定
3. 憲法文書の内容に基づいてConstitution Checkセクションを記入
4. 以下のConstitution Checkセクションを評価
   → 違反が存在する場合: Complexity Trackingに文書化
   → 正当化が不可能な場合: ERROR "Simplify approach first"
   → Progress Trackingを更新: Initial Constitution Check
5. フェーズ0を実行 → research.md
   → NEEDS CLARIFICATIONが残っている場合: ERROR "Resolve unknowns"
6. フェーズ1を実行 → contracts, data-model.md, quickstart.md, エージェント固有テンプレートファイル（例：Claude Code用`CLAUDE.md`、GitHub Copilot用`.github/copilot-instructions.md`、Gemini CLI用`GEMINI.md`、Qwen Code用`QWEN.md`、opencode用`AGENTS.md`）
7. Constitution Checkセクションを再評価
   → 新しい違反がある場合: 設計をリファクタ、フェーズ1に戻る
   → Progress Trackingを更新: Post-Design Constitution Check
8. フェーズ2を計画 → タスク生成アプローチを説明（tasks.mdは作成しない）
9. 停止 - /tasksコマンドの準備完了
```

**重要**: /planコマンドはステップ7で停止します。フェーズ2-4は他のコマンドで実行されます：
- フェーズ2: /tasksコマンドでtasks.mdを作成
- フェーズ3-4: 実装実行（手動またはツール使用）

## Summary
ユーザー管理CRUD操作のためのレイヤードアーキテクチャ（ドメイン、アプリケーション、インフラストラクチャ）を持つGo REST APIサーバー。PostgreSQLデータベース、GORM ORM、Mockeryテストモックを使用。Go標準プロジェクトレイアウトに従い、domain/infrastructureはinternal/に、application/controllerはcmd/api-server/internal/に配置。

## Technical Context
**Language/Version**: Go 1.21+（最新安定版）
**Primary Dependencies**: GORM（PostgreSQL ORM）、Gin/Echo（HTTPフレームワーク）、Mockery（モック生成）
**Storage**: GORMマイグレーションを使用したPostgreSQLデータベース
**Testing**: Mockery生成モックを使用したGoテストパッケージ、ソースファイルに隣接するテスト
**Target Platform**: Linuxサーバー（コンテナ化デプロイメント）
**Project Type**: 単一バックエンドAPIサーバー
**Performance Goals**: 標準的なREST APIパフォーマンス（1秒未満のレスポンス時間）
**Constraints**: レイヤードアーキテクチャの分離、Go標準プロジェクトレイアウトの準拠
**Scale/Scope**: 単一のUserエンティティでのCRUD操作、より大きなシステムの基盤

## Constitution Check
*ゲート: フェーズ0リサーチ前に合格が必要。フェーズ1設計後に再チェック。*

**Library-First Principle**: ✅ PASS - コアドメインロジックをinternal/domainに再利用可能なライブラリとして配置
**Interface Separation**: ✅ PASS - リポジトリインターフェースはドメインに、実装はインフラストラクチャに配置
**Test-First Approach**: ✅ PASS - Mockeryモックを使用してソースファイルに隣接するテストでTDD実施
**Standard Layout Compliance**: ✅ PASS - golang-standards/project-layoutに準拠
**Simplicity**: ✅ PASS - 単一エンティティ（User）の標準CRUD、早期複雑化なし

**初期憲法チェック**: PASS

## Project Structure

### Documentation (this feature)
```
specs/[###-feature]/
├── plan.md              # このファイル（/planコマンド出力）
├── research.md          # フェーズ0出力（/planコマンド）
├── data-model.md        # フェーズ1出力（/planコマンド）
├── quickstart.md        # フェーズ1出力（/planコマンド）
├── contracts/           # フェーズ1出力（/planコマンド）
└── tasks.md             # フェーズ2出力（/tasksコマンド - /planでは作成されない）
```

### Source Code (repository root)
```
# Go標準プロジェクトレイアウト（採用）
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
├── api/                   # OpenAPI/Swagger仕様、JSONスキーマファイル等
├── docs/                  # 設計とユーザードキュメント
├── scripts/               # ビルド、インストール、解析等を行うスクリプト
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

**Structure Decision**: Go標準プロジェクトレイアウトを採用し、レイヤードアーキテクチャを実装

## Phase 0: Outline & Research
1. **Technical Contextから不明点を抽出**:
   - 各NEEDS CLARIFICATIONに対してリサーチタスク作成
   - 各依存関係に対してベストプラクティスタスク作成
   - 各統合に対してパターンタスク作成

2. **リサーチエージェントの生成と派遣**:
   ```
   Technical Contextの各不明点に対して:
     タスク: "{機能コンテキスト}のための{不明点}のリサーチ"
   各技術選択に対して:
     タスク: "{ドメイン}における{技術}のベストプラクティス発見"
   ```

3. **research.mdでの調査結果統合**（以下の形式を使用）:
   - Decision: [選択されたもの]
   - Rationale: [選択理由]
   - Alternatives considered: [評価した他の選択肢]

**出力**: 全てのNEEDS CLARIFICATIONが解決されたresearch.md

## Phase 1: Design & Contracts
*前提条件: research.mdの完了*

1. **機能仕様からエンティティを抽出** → `data-model.md`:
   - エンティティ名、フィールド、関係性
   - 要求仕様からのバリデーションルール
   - 該当する場合の状態遷移

2. **機能要求からAPIコントラクトを生成**:
   - 各ユーザーアクション → エンドポイント
   - 標準的なREST/GraphQLパターンを使用
   - OpenAPI/GraphQLスキーマを`/contracts/`に出力

3. **コントラクトからコントラクトテストを生成**:
   - エンドポイントごとに1つのテストファイル
   - リクエスト/レスポンススキーマをアサート
   - テストは失敗する必要がある（まだ実装なし）

4. **ユーザーストーリーからテストシナリオを抽出**:
   - 各ストーリー → 統合テストシナリオ
   - クイックスタートテスト = ストーリー検証ステップ

5. **エージェントファイルの漸進的更新**（O(1)操作）:
   - `.specify/scripts/bash/update-agent-context.sh claude`を実行
     **重要**: 上記の通り正確に実行。引数の追加・削除は不可。
   - 存在する場合: 現在のプランからの新しい技術のみ追加
   - マーカー間の手動追加を保持
   - 最近の変更を更新（最新3つを保持）
   - トークン効率のため150行以下を維持
   - リポジトリルートに出力

**出力**: data-model.md、/contracts/*、失敗するテスト、quickstart.md、エージェント固有ファイル

## Phase 2: Task Planning Approach
*このセクションは/tasksコマンドが実行する内容を説明 - /plan実行中は実行しない*

**タスク生成戦略**:
- `.specify/templates/tasks-template.md`をベースとして読み込み
- フェーズ1設計ドキュメント（コントラクト、データモデル、クイックスタート）からタスクを生成
- 各コントラクト → コントラクトテストタスク [P]
- 各エンティティ → モデル作成タスク [P]
- 各ユーザーストーリー → 統合テストタスク
- テストを通すための実装タスク

**順序戦略**:
- TDD順序: 実装前にテスト
- 依存関係順序: モデル → サービス → UI
- [P]マークで並列実行（独立ファイル）をマーク

**推定出力**: tasks.mdに25-30の番号付き順序タスク

**重要**: このフェーズは/tasksコマンドで実行され、/planでは実行されない

## Phase 3+: Future Implementation
*これらのフェーズは/planコマンドの範囲外*

**フェーズ3**: タスク実行（/tasksコマンドでtasks.mdを作成）
**フェーズ4**: 実装（憲法原則に従ってtasks.mdを実行）
**フェーズ5**: 検証（テスト実行、quickstart.md実行、パフォーマンス検証）

## Complexity Tracking
*憲法チェックで正当化が必要な違反がある場合のみ記入*

| 違反 | 必要な理由 | より単純な代替案が却下された理由 |
|------|------------|--------------------------------|
| [例: 4番目のプロジェクト] | [現在の必要性] | [3つのプロジェクトでは不十分な理由] |
| [例: リポジトリパターン] | [具体的な問題] | [直接DB アクセスでは不十分な理由] |


## Progress Tracking
*このチェックリストは実行フロー中に更新*

**フェーズステータス**:
- [x] フェーズ0: リサーチ完了（/planコマンド）
- [x] フェーズ1: 設計完了（/planコマンド）
- [x] フェーズ2: タスク計画完了（/planコマンド - アプローチ説明のみ）
- [ ] フェーズ3: タスク生成（/tasksコマンド）
- [ ] フェーズ4: 実装完了
- [ ] フェーズ5: 検証合格

**ゲートステータス**:
- [x] 初期憲法チェック: PASS
- [x] 設計後憲法チェック: PASS
- [x] 全てのNEEDS CLARIFICATION解決
- [x] 複雑性逸脱の文書化（N/A - 逸脱なし）

---
*憲法 v2.1.1に基づく - `/memory/constitution.md`を参照*
