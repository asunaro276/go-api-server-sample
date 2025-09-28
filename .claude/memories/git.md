# Git 運用ガイドライン

このファイルは、Claude Codeがプロジェクトでのgit操作を行う際に遵守すべきgit運用ガイドラインを定義します。

## 禁止コマンド

### 絶対に使用してはいけないコマンド

```bash
# 全ファイル一括追加（禁止）
git add .

# 強制プッシュ（禁止）
git push -f
git push --force
git push --force-with-lease  # 特別な理由がない限り禁止
```

### 理由

- **`git add .`**: 意図しないファイルや機密情報をコミットしてしまうリスク
- **強制プッシュ**: チーム開発でのコンフリクトやデータ損失の原因

### 代替手法

```bash
# ファイルを個別または関連ファイルごとに追加
git add internal/domain/entities/user.go
git add internal/domain/entities/user_test.go

# 特定のディレクトリ単位で追加
git add internal/domain/entities/

# 変更確認してから追加
git status
git diff
git add [具体的なファイル名]
```

## コミット単位の原則

### 基本方針

1. **命名しやすく意味のある単位でコミット**
2. **エラーや不整合が起こらない単位でコミット**
3. **機能的に独立した変更をまとめる**

### 適切なコミット単位の例

✅ **良いコミット例**
```
feat: ユーザーエンティティを追加
feat: ユーザー作成ユースケースを実装
feat: ユーザーリポジトリインターフェースを定義
test: ユーザーエンティティのバリデーションテストを追加
fix: ユーザー更新時のバリデーションエラーを修正
refactor: ユーザーコントローラのエラーハンドリングを改善
docs: READMEにユーザー管理APIの仕様を追加
```

❌ **悪いコミット例**
```
WIP: 途中まで実装  # 未完成のコミット
fix: バグ修正      # 具体性がない
update            # 何を更新したか不明
ユーザー機能を全部実装  # 範囲が広すぎる
```

### コミット単位の具体例

#### 1機能1コミットの例
```bash
# エンティティ追加
git add internal/domain/entities/user.go
git add internal/domain/entities/user_test.go
git commit -m "feat: ユーザーエンティティとバリデーションを追加"

# ユースケース実装
git add cmd/api-server/internal/application/createuser.go
git add cmd/api-server/internal/application/createuser_test.go
git commit -m "feat: ユーザー作成ユースケースを実装"

# リポジトリ実装
git add internal/infrastructure/repositories/user.go
git add internal/infrastructure/repositories/user_test.go
git commit -m "feat: ユーザーリポジトリ実装を追加"
```

#### エラー修正の例
```bash
# バグ修正（関連ファイルをまとめてコミット）
git add internal/domain/entities/user.go
git add internal/domain/entities/user_test.go
git commit -m "fix: ユーザーメールアドレスバリデーションの正規表現を修正"
```

### コミット前チェックリスト

1. **コンパイルエラーがないか確認**
   ```bash
   make build
   ```

2. **テストが通るか確認**
   ```bash
   make test
   ```

3. **リントエラーがないか確認**
   ```bash
   make lint
   ```

4. **関連ファイルがすべて含まれているか確認**
   ```bash
   git status
   git diff --cached
   ```

## ブランチ戦略

### ブランチ命名規則

```bash
# 機能追加
feature/[実装機能名]

# 例
feature/user-management
feature/user-authentication
feature/product-catalog
feature/order-processing
```

### ブランチ作成手順

```bash
# 1. mainブランチに移動
git checkout main

# 2. 最新化
git pull origin main

# 3. 新しいfeatureブランチを作成
git checkout -b feature/user-management

# 4. 作業開始
# ... 実装作業 ...

# 5. リモートにプッシュ
git push -u origin feature/user-management
```

### ブランチ運用ルール

1. **必ずmainブランチの最新から切る**
2. **featureブランチは短期間で完了させる**
3. **マージ前に必ずテストを実行**
4. **不要になったブランチは削除**

### プルリクエスト作成前チェック

```bash
# 1. 最新のmainブランチをマージ
git checkout main
git pull origin main
git checkout feature/your-feature
git merge main

# 2. コンフリクト解決後、テスト実行
make test-all

# 3. 全チェックが通ることを確認
make ci
```

## コミットメッセージ規約

### フォーマット

```
<type>: <description>

<body>
```

### Type一覧

- **feat**: 新機能追加
- **fix**: バグ修正
- **refactor**: リファクタリング
- **test**: テスト追加・修正
- **docs**: ドキュメント追加・修正
- **style**: コードフォーマット修正
- **chore**: その他（依存関係更新など）

### 例

```bash
git commit -m "feat: ユーザー認証機能を追加

- JWTトークンベースの認証を実装
- ログイン・ログアウトエンドポイントを追加
- 認証ミドルウェアを実装"
```

## 緊急時の対応

### 間違ったコミットをした場合

```bash
# 直前のコミットを取り消し（変更は保持）
git reset --soft HEAD~1

# コミットメッセージのみ修正
git commit --amend -m "正しいコミットメッセージ"
```

### 間違ったブランチで作業した場合

```bash
# 変更を一時保存
git stash

# 正しいブランチに移動
git checkout correct-branch

# 変更を復元
git stash pop
```

## コードレビューポイント

### 必須チェック項目

1. **禁止コマンドが使用されていないか**
2. **コミット単位が適切か**
3. **テストが通るか**
4. **リントエラーがないか**
5. **ブランチ命名規則に従っているか**
6. **コミットメッセージが適切か**