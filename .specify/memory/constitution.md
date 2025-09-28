<!--
Sync Impact Report - Constitution Update to v1.0.0
===========================================
Version Change: Initial → v1.0.0 (NEW)
Ratified: 2025-09-28 | Last Amended: 2025-09-28

Core Principles Established:
- I. Clean Architecture（NON-NEGOTIABLE）: 4層アーキテクチャの厳格な遵守
- II. テスト駆動開発（TDD）: Test Suites、Testcontainers、モック生成
- III. 命名規約・ファイル構成: 小文字繋ぎファイル名（スネーク・ケバブケース禁止）、1ユースケース1ファイル
- IV. Git運用ガイドライン: 機能単位コミット、禁止コマンド遵守
- V. 依存関係管理: 依存性注入、抽象化による分離

Added Sections:
- 技術スタック標準
- 開発ワークフロー
- ガバナンス

テンプレート更新状況:
✅ plan-template.mdのConstitution Checkルールと整合済み
✅ .claude/memories/への参照構造を維持
✅ tasks-template.mdがTDD原則と互換性確認済み
✅ spec-template.mdがClean Architecture要件と互換性確認済み

フォローアップTODO:
- なし - 既存の.claude/memories/コンテンツを使用してすべてのプレースホルダーを解決済み
-->

# go-api-server-sample Constitution

## Core Principles

### I. Clean Architecture（NON-NEGOTIABLE）
本プロジェクトは4層Clean Architectureを厳格に遵守する。Domain層は他の層に依存せず、依存関係は内側への一方向のみとする。Infrastructure層への依存は抽象化（インターフェース）経由でのみ行う。

詳細なアーキテクチャガイドライン：@.claude/memories/dependency.md を参照

### II. テスト駆動開発（TDD）
すべての実装はテスト駆動で行う。testify/suiteを使用したTest Suite構造、リポジトリテストでのTestcontainers利用、mockery v3によるモック生成を標準とする。実装前にテストが書かれ、失敗することを確認してから実装を開始する。

詳細なテストガイドライン：@.claude/memories/testing.md を参照

### III. 命名規約・ファイル構成
ファイル名は小文字繋ぎとし、スネークケース（アンダースコア）およびケバブケース（ハイフン区切り）は使用禁止。1ユースケース1ファイルの原則を遵守し、ディレクトリ名の重複を避ける。

詳細な命名規約：@.claude/memories/naming.md および @.claude/memories/organization.md を参照

### IV. Git運用ガイドライン
機能単位での意味のあるコミットを行い、`git add .`や強制プッシュは禁止。コミット前には必ずビルド、テスト、リントが通ることを確認し、適切なコミットメッセージフォーマットを使用する。

詳細なGit運用ガイドライン：@.claude/memories/git.md を参照

### V. 依存関係管理
依存性注入パターンを使用し、各層の責務を明確に分離する。UseCase層はDomain層のみに依存し、Controller層はUseCase層経由でビジネスロジックにアクセスする。外部ライブラリとの結合度を最小化する。

## 技術スタック標準

**必須技術スタック**:
- 言語: Go 1.24+
- Webフレームワーク: Gin
- ORM: GORM
- データベース: PostgreSQL
- テスト: testify、testcontainers
- モック生成: mockery v3
- リント: golangci-lint
- 開発ツール: air（ホットリロード）

**パフォーマンス要件**: APIレスポンス < 200ms p95、メモリ使用量の監視、適切なデータベースインデックス設計

## 開発ワークフロー

**標準実行手順**:
1. makeコマンドによる開発環境セットアップ（`make quickstart`）
2. 機能分岐作成（`feature/[機能名]`形式）
3. TDDサイクル（Red-Green-Refactor）
4. 段階的コミット（機能単位）
5. CI/CDパイプライン通過確認（`make ci`）

**品質ゲート**: 全テスト通過、リントエラーゼロ、カバレッジ要件達成、コードレビュー承認

## ガバナンス

**優先順位**: 本憲章がすべての他の慣行に優先する。修正には文書化、承認、移行計画が必要。

**コンプライアンス審査**: すべてのPR/レビューで遵守確認を実施。複雑性は正当化が必要。実行時の開発ガイダンスは CLAUDE.md を使用。

**修正手順**: MINOR変更（新原則追加）またはPATCH変更（明確化・誤字修正）のみ。MAJOR変更（後方互換性のない変更）は避ける。

**Version**: 1.0.0 | **Ratified**: 2025-09-28 | **Last Amended**: 2025-09-28