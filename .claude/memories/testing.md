# テストガイドライン

このプロジェクトにおけるテストの書き方と運用に関するガイドラインです。

## テストファイルの配置

### 基本原則

- **テストファイルはソースコードと同一ディレクトリに配置**
- **パッケージ名もソースコードと同一にする**
- **ファイル名は `*_test.go` の形式**

```
cmd/api-server/internal/usecase/
├── user.go
├── user_test.go         # 同一ディレクトリ、同一パッケージ
└── ...

internal/domain/entities/
├── user.go
├── user_test.go                 # 同一ディレクトリ、同一パッケージ
└── ...

internal/infrastructure/repositories/
├── user_repository.go
├── user_repository_test.go      # 同一ディレクトリ、同一パッケージ
└── ...
```

### パッケージ宣言例

```go
// user.go
package application

// user_test.go
package application  // 同一パッケージ
```

## モックの利用

### Mockery v3を使用したモック生成

インターフェースに基づく依存性のモックは `mockery v3` を使用して自動生成します。

#### モック生成コマンド

```bash
# 全モック生成
make mock

# 特定のインターフェースのモック生成
mockery --name=UserRepository --dir=internal/domain/repositories --output=internal/testing/mocks
```

#### モック使用例

```go
type UserUsecaseTestSuite struct {
    suite.Suite
    mockRepo *mocks.UserRepository
    usecase  UserUsecaseInterface
}

func (suite *UserUsecaseTestSuite) SetupTest() {
    // モックの初期化
}

func (suite *UserUsecaseTestSuite) TestCreateUser() {
    // Given-When-Thenパターンでテスト実装
}

func TestUserUsecaseTestSuite(t *testing.T) {
    suite.Run(t, new(UserUsecaseTestSuite))
}
```

### リポジトリテストでのTestcontainers利用

リポジトリのテストでは実際のデータベースを使用するため、testcontainersを利用します。

#### 共通テストスイートの作成

`internal/infrastructure/repositories/testcontainers.go` に共通のテストスイートを作成し、他のテストで再利用します。

```go
type DatabaseTestSuite struct {
    suite.Suite
    container *postgres.PostgresContainer
    db        *gorm.DB
}

func (suite *DatabaseTestSuite) SetupSuite() {
    // PostgreSQLコンテナ起動、DB接続、マイグレーション実行
}

func (suite *DatabaseTestSuite) TearDownSuite() {
    // コンテナ停止
}

func (suite *DatabaseTestSuite) SetupSubTest() {
    // テストデータクリーンアップ
}

func (suite *DatabaseTestSuite) GetDB() *gorm.DB {
    return suite.db
}
```

#### リポジトリテストでの利用例

```go
type UserRepositoryTestSuite struct {
    DatabaseTestSuite  // 共通テストスイートを埋め込み
    repo UserRepository
}

func (suite *UserRepositoryTestSuite) SetupTest() {
    suite.repo = NewUserRepository(suite.GetDB())
}

func (suite *UserRepositoryTestSuite) TestCreateUser() {
    suite.Run("正常にユーザーが作成される", func() {
        // Given-When-Thenパターンでテスト実装
    })

    suite.Run("重複メールアドレスでエラーになる", func() {
        // Given-When-Thenパターンでテスト実装
    })
}

func TestUserRepositoryTestSuite(t *testing.T) {
    suite.Run(t, new(UserRepositoryTestSuite))
}
```

## テストの書き方

### Test Suiteの使用

基本的に全てのテストで `testify/suite` を使用します。

#### 基本構造

```go
type ExampleTestSuite struct {
    suite.Suite
    // テスト用のフィールド
}

func (suite *ExampleTestSuite) SetupSuite() {
    // スイート開始時に1回実行
}

func (suite *ExampleTestSuite) TearDownSuite() {
    // スイート終了時に1回実行
}

func (suite *ExampleTestSuite) SetupTest() {
    // 各テスト開始時に実行
}

func (suite *ExampleTestSuite) SetupSubTest() {
    // 各サブテスト開始時に実行
}

func (suite *ExampleTestSuite) TestExample() {
    suite.Run("ケース1", func() {
        // テスト実装
    })
}

func TestExampleTestSuite(t *testing.T) {
    suite.Run(t, new(ExampleTestSuite))
}
```

### テーブルドリブンテストの書き方

```go
func (suite *UserUsecaseTestSuite) TestValidateUser() {
    tests := []struct {
        name    string
        user    *entities.User
        wantErr bool
        errType error
    }{
        // テストケース定義
    }

    for _, tt := range tests {
        suite.Run(tt.name, func() {
            // Given-When-Thenパターンでテスト実装
        })
    }
}
```

## テストの種類と実行方法

### 単体テスト

```bash
# 全単体テスト実行
make test

# 特定パッケージ
go test -v ./internal/domain/entities

# カバレッジ付き
make test-coverage
```

### 統合テスト

```bash
# 統合テスト実行
make test-integration
```

### 契約テスト

```bash
# 契約テスト実行
make test-contract
```

### パフォーマンステスト

```bash
# パフォーマンステスト実行
make test-performance
```

### 全テスト実行

```bash
# 全種類のテスト実行
make test-all
```

## ベストプラクティス

### 1. テストの独立性

- 各テスト/サブテストは独立して実行可能であること
- SetupSubTest/TearDownSubTestでデータをクリーンアップ
- テストの実行順序に依存しない設計

### 2. 冪等性の確保

- 同じテストを複数回実行しても同じ結果になること
- データベースの状態をテスト毎にリセット

### 3. 可読性

- テスト名は日本語で具体的に記述
- Given-When-Thenパターンを使用
- サブテストで複数のケースを整理

### 4. モックの適切な使用

- 外部依存はモックで代替
- アサーションでモックの呼び出しを検証
- 過度なモックは避け、必要最小限に留める

### 5. テストデータ管理

- テストデータは各テスト内で定義
- 共通データは適切にセットアップ/クリーンアップ
- テストcontainersでのデータ分離
