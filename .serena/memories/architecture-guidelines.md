# アーキテクチャガイドライン

## 4層アーキテクチャ

本プロジェクトではClean Architectureに基づく4層構造を採用しています。

### 1. Domain層 (`internal/domain/`)

**責務**
- ドメインエンティティの定義
- ビジネスルールの実装
- リポジトリインターフェースの定義

**依存関係**
- 他の層に依存しない（最内層）
- 外部のフレームワークやライブラリに依存しない

**構成**
- `entities/` - ドメインエンティティ
- `valueobjects/` - バリューオブジェクト
- `repositories/` - リポジトリインターフェース

### 2. UseCase層 (`cmd/api-server/internal/usecases/`)

**責務**
- アプリケーションのビジネスロジック
- ユースケースの実装
- トランザクション境界の定義

**依存関係**
- Domain層のみに依存
- Infrastructure層への依存は抽象化（インターフェース）経由のみ


### 3. Infrastructure層 (`internal/infrastructure/`)

**責務**
- 外部システムとの連携（DB、API等）
- リポジトリの具象実装
- 設定管理
- 外部ライブラリの抽象化

**依存関係**
- Domain層に依存（リポジトリインターフェースの実装）
- UseCase層には依存しない

**構成**
- `database/` - データベース関連
- `repositories/` - リポジトリ実装

### 4. Controller層 (`cmd/api-server/internal/controller/`)

**責務**
- HTTPリクエスト/レスポンスの処理
- バリデーション
- UseCaseの呼び出し
- エラーハンドリング

**依存関係**
- UseCase層に依存
- Infrastructure層への直接依存は避ける
