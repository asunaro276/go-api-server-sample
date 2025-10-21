---
name: design-expert
description: ソフトウェア設計に関する包括的なガイド。設計パターン、アーキテクチャパターン、設計原則（SOLID、DRYなど）、システム設計、API設計、データモデル設計などに関するアドバイスや実装を求める場合に使用すべきスキルです。クリーンアーキテクチャ、ヘキサゴナルアーキテクチャ、DDDなどのアーキテクチャスタイルにも対応しています。
---

# 設計エキスパート

## 概要

このスキルは、ソフトウェア設計に関する包括的な知識とベストプラクティスを提供します。設計パターン、アーキテクチャパターン、設計原則、システム設計、API設計など、ソフトウェア開発における設計面での意思決定をサポートします。

## クイックスタート

### 基本的な使用パターン

設計に関する質問や実装を依頼する際の基本的なパターン：

```
# 設計パターンの適用
「ユーザー認証システムにFactory Patternを適用してください」

# アーキテクチャの相談
「このマイクロサービスにどのアーキテクチャパターンが適していますか？」

# 設計原則の確認
「このコードはSOLID原則に従っていますか？改善点を教えてください」

# API設計のレビュー
「このREST APIの設計をレビューして、改善点を提案してください」
```

### いつこのスキルを使うべきか

- 新しい機能やモジュールの設計を開始する前
- 既存コードのリファクタリングを検討する時
- アーキテクチャの選択に迷った時
- 設計パターンの適用方法が不明な時
- コードレビューで設計面の改善を提案したい時
- スケーラビリティやメンテナンス性を向上させたい時

## コア設計領域

### 1. 設計パターン

#### 生成パターン（Creational Patterns）

**Singleton Pattern**
```go
// スレッドセーフなSingleton実装
type DatabaseConnection struct {
    conn *sql.DB
}

var (
    instance *DatabaseConnection
    once     sync.Once
)

func GetInstance() *DatabaseConnection {
    once.Do(func() {
        instance = &DatabaseConnection{
            conn: initConnection(),
        }
    })
    return instance
}
```

**Factory Pattern**
```go
// インターフェースベースのFactory
type UserRepository interface {
    Create(user *User) error
    FindByID(id string) (*User, error)
}

type RepositoryFactory struct{}

func (f *RepositoryFactory) CreateUserRepository(dbType string) UserRepository {
    switch dbType {
    case "postgres":
        return NewPostgresUserRepository()
    case "mongodb":
        return NewMongoUserRepository()
    default:
        return NewInMemoryUserRepository()
    }
}
```

**Builder Pattern**
```go
// 複雑なオブジェクト構築のためのBuilder
type RequestBuilder struct {
    request *http.Request
}

func NewRequestBuilder() *RequestBuilder {
    return &RequestBuilder{request: &http.Request{}}
}

func (b *RequestBuilder) WithMethod(method string) *RequestBuilder {
    b.request.Method = method
    return b
}

func (b *RequestBuilder) WithURL(url string) *RequestBuilder {
    b.request.URL, _ = neturl.Parse(url)
    return b
}

func (b *RequestBuilder) Build() *http.Request {
    return b.request
}
```

#### 構造パターン（Structural Patterns）

**Adapter Pattern**
```go
// 外部ライブラリを内部インターフェースに適合させる
type Logger interface {
    Info(msg string)
    Error(msg string)
}

type ZapLoggerAdapter struct {
    zapLogger *zap.Logger
}

func (a *ZapLoggerAdapter) Info(msg string) {
    a.zapLogger.Info(msg)
}

func (a *ZapLoggerAdapter) Error(msg string) {
    a.zapLogger.Error(msg)
}
```

**Decorator Pattern**
```go
// 機能を動的に追加
type Handler interface {
    Handle(ctx context.Context) error
}

type LoggingDecorator struct {
    handler Handler
    logger  Logger
}

func (d *LoggingDecorator) Handle(ctx context.Context) error {
    d.logger.Info("処理開始")
    err := d.handler.Handle(ctx)
    if err != nil {
        d.logger.Error("処理失敗")
    }
    return err
}
```

#### 振る舞いパターン（Behavioral Patterns）

**Strategy Pattern**
```go
// アルゴリズムの切り替えを可能にする
type PaymentStrategy interface {
    Pay(amount float64) error
}

type CreditCardPayment struct{}
func (c *CreditCardPayment) Pay(amount float64) error { /* ... */ }

type PayPalPayment struct{}
func (p *PayPalPayment) Pay(amount float64) error { /* ... */ }

type PaymentProcessor struct {
    strategy PaymentStrategy
}

func (p *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

func (p *PaymentProcessor) ProcessPayment(amount float64) error {
    return p.strategy.Pay(amount)
}
```

**Observer Pattern**
```go
// イベント駆動アーキテクチャの基礎
type Observer interface {
    Update(event Event)
}

type Subject struct {
    observers []Observer
}

func (s *Subject) Attach(observer Observer) {
    s.observers = append(s.observers, observer)
}

func (s *Subject) Notify(event Event) {
    for _, observer := range s.observers {
        observer.Update(event)
    }
}
```

### 2. アーキテクチャパターン

#### クリーンアーキテクチャ

**レイヤー構造**
```
clean-architecture/
├── domain/              # エンティティとビジネスロジック
│   ├── entity/         # ドメインエンティティ
│   ├── repository/     # リポジトリインターフェース
│   └── service/        # ドメインサービス
├── usecase/            # ユースケース層（アプリケーションロジック）
│   ├── interactor/     # ビジネスロジックの実装
│   └── port/           # 入出力ポート
├── interface/          # インターフェース層
│   ├── handler/        # HTTPハンドラー
│   ├── presenter/      # プレゼンター
│   └── repository/     # リポジトリ実装
└── infrastructure/     # インフラストラクチャ層
    ├── database/       # データベース接続
    ├── external/       # 外部API
    └── config/         # 設定
```

**依存性の方向**
```
infrastructure → interface → usecase → domain
                                         ↑
                                    すべてがこの方向に依存
```

**実装例**
```go
// domain/entity/user.go
type User struct {
    ID        string
    Email     string
    CreatedAt time.Time
}

// domain/repository/user_repository.go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}

// usecase/interactor/user_interactor.go
type UserInteractor struct {
    userRepo repository.UserRepository
}

func (i *UserInteractor) CreateUser(ctx context.Context, email string) (*User, error) {
    user := &User{
        ID:        generateID(),
        Email:     email,
        CreatedAt: time.Now(),
    }
    if err := i.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }
    return user, nil
}

// interface/repository/postgres_user_repository.go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *User) error {
    // データベース実装
}
```

#### ヘキサゴナルアーキテクチャ（ポート&アダプター）

**原則**
- アプリケーションコアは外部依存から独立
- ポート：アプリケーションの境界を定義するインターフェース
- アダプター：外部システムとポートをつなぐ実装

**実装例**
```go
// ポート（インバウンド）
type UserService interface {
    RegisterUser(email string) error
}

// ポート（アウトバウンド）
type UserStore interface {
    Save(user User) error
}

// アプリケーションコア
type UserServiceImpl struct {
    store UserStore
}

func (s *UserServiceImpl) RegisterUser(email string) error {
    user := User{Email: email}
    return s.store.Save(user)
}

// アダプター（インバウンド - HTTP）
type HTTPAdapter struct {
    service UserService
}

func (h *HTTPAdapter) HandleRegister(w http.ResponseWriter, r *http.Request) {
    // HTTPリクエストをサービスコールに変換
}

// アダプター（アウトバウンド - PostgreSQL）
type PostgreSQLAdapter struct {
    db *sql.DB
}

func (p *PostgreSQLAdapter) Save(user User) error {
    // データベース保存処理
}
```

#### レイヤードアーキテクチャ

**構造**
```
プレゼンテーション層（UI/API）
        ↓
  ビジネスロジック層
        ↓
   データアクセス層
        ↓
    データベース
```

### 3. 設計原則

#### SOLID原則

**S - Single Responsibility Principle（単一責任の原則）**
```go
// 悪い例：複数の責任を持つ
type UserManager struct{}

func (m *UserManager) CreateUser(user User) error { /* ... */ }
func (m *UserManager) SendWelcomeEmail(user User) error { /* ... */ }
func (m *UserManager) LogActivity(activity string) error { /* ... */ }

// 良い例：責任を分離
type UserRepository struct{}
func (r *UserRepository) Create(user User) error { /* ... */ }

type EmailService struct{}
func (s *EmailService) SendWelcome(user User) error { /* ... */ }

type ActivityLogger struct{}
func (l *ActivityLogger) Log(activity string) error { /* ... */ }
```

**O - Open/Closed Principle（開放閉鎖の原則）**
```go
// インターフェースを使って拡張可能に
type PaymentProcessor interface {
    Process(amount float64) error
}

type PaymentService struct {
    processors []PaymentProcessor
}

func (s *PaymentService) AddProcessor(p PaymentProcessor) {
    s.processors = append(s.processors, p)
}

// 新しいプロセッサを追加しても既存コードは変更不要
type BitcoinProcessor struct{}
func (b *BitcoinProcessor) Process(amount float64) error { /* ... */ }
```

**L - Liskov Substitution Principle（リスコフの置換原則）**
```go
// 基底インターフェース
type Storage interface {
    Save(key string, value []byte) error
    Load(key string) ([]byte, error)
}

// どの実装も同じ契約を守る
type MemoryStorage struct{}
type FileStorage struct{}
type S3Storage struct{}

// すべてStorage型として置換可能
func ProcessData(storage Storage, key string, data []byte) error {
    return storage.Save(key, data)
}
```

**I - Interface Segregation Principle（インターフェース分離の原則）**
```go
// 悪い例：大きすぎるインターフェース
type Repository interface {
    Create(entity Entity) error
    Read(id string) (Entity, error)
    Update(entity Entity) error
    Delete(id string) error
    Search(query string) ([]Entity, error)
    Export(format string) ([]byte, error)
}

// 良い例：必要な機能だけのインターフェース
type Reader interface {
    Read(id string) (Entity, error)
}

type Writer interface {
    Create(entity Entity) error
    Update(entity Entity) error
}

type Searcher interface {
    Search(query string) ([]Entity, error)
}
```

**D - Dependency Inversion Principle（依存性逆転の原則）**
```go
// 悪い例：具体実装に依存
type UserService struct {
    db *PostgresDB  // 具体的な実装に依存
}

// 良い例：抽象に依存
type UserService struct {
    repo UserRepository  // インターフェースに依存
}

type UserRepository interface {
    Save(user User) error
}

// 実装は差し替え可能
type PostgresUserRepository struct{}
type MongoUserRepository struct{}
```

#### DRY原則（Don't Repeat Yourself）

```go
// 悪い例：重複コード
func (s *Service) CreateUser(user User) error {
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email")
    }
    // ...
}

func (s *Service) UpdateUser(user User) error {
    if user.Email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(user.Email, "@") {
        return errors.New("invalid email")
    }
    // ...
}

// 良い例：バリデーションを共通化
func validateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(email, "@") {
        return errors.New("invalid email")
    }
    return nil
}

func (s *Service) CreateUser(user User) error {
    if err := validateEmail(user.Email); err != nil {
        return err
    }
    // ...
}
```

#### KISS原則（Keep It Simple, Stupid）

```go
// 悪い例：過度に複雑
func (s *Service) ProcessData(data interface{}) (interface{}, error) {
    switch v := data.(type) {
    case map[string]interface{}:
        // 複雑な処理...
    case []interface{}:
        // さらに複雑な処理...
    default:
        // ...
    }
}

// 良い例：シンプルで明確
func (s *Service) ProcessUser(user User) (*ProcessedUser, error) {
    // 明確な入力と出力
}

func (s *Service) ProcessUsers(users []User) ([]*ProcessedUser, error) {
    // 明確な入力と出力
}
```

#### YAGNI原則（You Aren't Gonna Need It）

今必要でない機能は実装しない。将来必要になるかもしれない機能のための過度な抽象化を避ける。

```go
// 悪い例：現在不要な機能を先に実装
type UserRepository interface {
    Create(user User) error
    Read(id string) (User, error)
    Update(user User) error
    Delete(id string) error
    BulkCreate(users []User) error        // まだ不要
    Archive(id string) error               // まだ不要
    Restore(id string) error               // まだ不要
    ExportToCSV() ([]byte, error)         // まだ不要
}

// 良い例：現在必要な機能のみ
type UserRepository interface {
    Create(user User) error
    Read(id string) (User, error)
}

// 必要になったら追加
```

### 4. システム設計

#### スケーラビリティパターン

**水平スケーリング**
```go
// ステートレスなサービス設計
type APIHandler struct {
    db    Database      // 共有状態
    cache CacheService  // 共有キャッシュ
}

// インスタンス間で状態を共有しない
func (h *APIHandler) HandleRequest(w http.ResponseWriter, r *http.Request) {
    // リクエスト毎に独立した処理
}
```

**キャッシング戦略**
```go
// Cache-Aside Pattern
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // 1. キャッシュを確認
    if user, err := s.cache.Get(ctx, id); err == nil {
        return user, nil
    }

    // 2. キャッシュミス時はDBから取得
    user, err := s.db.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // 3. キャッシュに保存
    s.cache.Set(ctx, id, user, 10*time.Minute)

    return user, nil
}
```

**非同期処理**
```go
// メッセージキューを使った非同期処理
type EventPublisher interface {
    Publish(ctx context.Context, event Event) error
}

func (s *Service) CreateOrder(ctx context.Context, order Order) error {
    // 1. 注文を保存（同期）
    if err := s.orderRepo.Save(ctx, order); err != nil {
        return err
    }

    // 2. 通知などの重い処理は非同期で実行
    event := OrderCreatedEvent{OrderID: order.ID}
    if err := s.eventPublisher.Publish(ctx, event); err != nil {
        // ログを記録するが、エラーは返さない
        s.logger.Error("failed to publish event", err)
    }

    return nil
}
```

#### 障害耐性パターン

**Circuit Breaker Pattern**
```go
type CircuitBreaker struct {
    maxFailures int
    timeout     time.Duration
    failures    int
    lastFailure time.Time
    state       string // "closed", "open", "half-open"
    mu          sync.Mutex
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    if cb.state == "open" {
        if time.Since(cb.lastFailure) > cb.timeout {
            cb.state = "half-open"
        } else {
            return errors.New("circuit breaker is open")
        }
    }

    err := fn()
    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
        }
        return err
    }

    cb.failures = 0
    cb.state = "closed"
    return nil
}
```

**Retry Pattern**
```go
func RetryWithBackoff(ctx context.Context, maxRetries int, fn func() error) error {
    var err error
    for i := 0; i < maxRetries; i++ {
        if err = fn(); err == nil {
            return nil
        }

        backoff := time.Duration(math.Pow(2, float64(i))) * time.Second
        select {
        case <-time.After(backoff):
            continue
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return err
}
```

### 5. API設計

#### RESTful API設計原則

**リソース指向の設計**
```
GET    /api/v1/users          # ユーザー一覧取得
GET    /api/v1/users/:id      # 特定ユーザー取得
POST   /api/v1/users          # ユーザー作成
PUT    /api/v1/users/:id      # ユーザー更新
DELETE /api/v1/users/:id      # ユーザー削除

# ネストしたリソース
GET    /api/v1/users/:id/posts       # 特定ユーザーの投稿一覧
POST   /api/v1/users/:id/posts       # ユーザーの投稿作成
```

**適切なHTTPステータスコード**
```go
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest) // 400
        return
    }

    user, err := h.service.CreateUser(r.Context(), req)
    if err != nil {
        if errors.Is(err, ErrEmailAlreadyExists) {
            http.Error(w, "Email already exists", http.StatusConflict) // 409
            return
        }
        http.Error(w, "Internal error", http.StatusInternalServerError) // 500
        return
    }

    w.WriteHeader(http.StatusCreated) // 201
    json.NewEncoder(w).Encode(user)
}
```

**バージョニング**
```go
// URLパスでバージョン管理
router.HandleFunc("/api/v1/users", handlerV1.GetUsers)
router.HandleFunc("/api/v2/users", handlerV2.GetUsers)

// ヘッダーでバージョン管理
func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
    version := r.Header.Get("API-Version")
    switch version {
    case "2.0":
        h.getUsersV2(w, r)
    default:
        h.getUsersV1(w, r)
    }
}
```

**ページネーション**
```go
type PaginationRequest struct {
    Page     int `json:"page"`
    PageSize int `json:"page_size"`
}

type PaginationResponse struct {
    Items      []interface{} `json:"items"`
    TotalCount int           `json:"total_count"`
    Page       int           `json:"page"`
    PageSize   int           `json:"page_size"`
    TotalPages int           `json:"total_pages"`
}

// GET /api/v1/users?page=1&page_size=20
```

#### GraphQL設計原則

**スキーマファースト設計**
```graphql
type Query {
  user(id: ID!): User
  users(limit: Int, offset: Int): [User!]!
}

type Mutation {
  createUser(input: CreateUserInput!): User!
  updateUser(id: ID!, input: UpdateUserInput!): User!
}

type User {
  id: ID!
  email: String!
  posts: [Post!]!
}

input CreateUserInput {
  email: String!
  name: String!
}
```

### 6. データモデル設計

#### 正規化と非正規化

**正規化（RDBMS）**
```sql
-- 第3正規形
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE tags (
    id UUID PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);

CREATE TABLE post_tags (
    post_id UUID REFERENCES posts(id),
    tag_id UUID REFERENCES tags(id),
    PRIMARY KEY (post_id, tag_id)
);
```

**非正規化（NoSQL）**
```go
// MongoDB用のドキュメント設計
type Post struct {
    ID        primitive.ObjectID `bson:"_id"`
    UserID    primitive.ObjectID `bson:"user_id"`
    UserEmail string             `bson:"user_email"` // 非正規化
    Title     string             `bson:"title"`
    Content   string             `bson:"content"`
    Tags      []string           `bson:"tags"` // 埋め込み
    CreatedAt time.Time          `bson:"created_at"`
}
```

#### ドメイン駆動設計（DDD）

**エンティティと値オブジェクト**
```go
// エンティティ：IDを持ち、ライフサイクルがある
type User struct {
    id       UserID
    email    Email
    profile  Profile
    version  int
}

// 値オブジェクト：不変で、等価性は値で判断
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, errors.New("invalid email")
    }
    return Email{value: value}, nil
}

func (e Email) String() string {
    return e.value
}
```

**集約（Aggregate）**
```go
// 集約ルート
type Order struct {
    id         OrderID
    customerId CustomerID
    items      []OrderItem
    status     OrderStatus
    total      Money
}

// 集約内のエンティティ
type OrderItem struct {
    productId ProductID
    quantity  int
    price     Money
}

// 集約を通じてのみ変更
func (o *Order) AddItem(productId ProductID, quantity int, price Money) error {
    if o.status != OrderStatusDraft {
        return errors.New("cannot modify confirmed order")
    }

    item := OrderItem{
        productId: productId,
        quantity:  quantity,
        price:     price,
    }
    o.items = append(o.items, item)
    o.calculateTotal()
    return nil
}

func (o *Order) calculateTotal() {
    // 集約内で整合性を保つ
    total := Money{amount: 0}
    for _, item := range o.items {
        total = total.Add(item.price.Multiply(item.quantity))
    }
    o.total = total
}
```

## ベストプラクティス

### 1. 依存性注入（Dependency Injection）

**コンストラクタインジェクション**
```go
type UserService struct {
    repo   UserRepository
    logger Logger
    cache  CacheService
}

func NewUserService(
    repo UserRepository,
    logger Logger,
    cache CacheService,
) *UserService {
    return &UserService{
        repo:   repo,
        logger: logger,
        cache:  cache,
    }
}
```

**インターフェースの活用**
```go
// テスト可能な設計
type UserService struct {
    repo UserRepository // インターフェース
}

// モックを使ったテスト
type MockUserRepository struct {
    users map[string]*User
}

func (m *MockUserRepository) FindByID(id string) (*User, error) {
    user, ok := m.users[id]
    if !ok {
        return nil, ErrNotFound
    }
    return user, nil
}
```

### 2. エラーハンドリング

**カスタムエラー型**
```go
type DomainError struct {
    Code    string
    Message string
    Cause   error
}

func (e *DomainError) Error() string {
    return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
    return e.Cause
}

var (
    ErrNotFound      = &DomainError{Code: "NOT_FOUND", Message: "resource not found"}
    ErrUnauthorized  = &DomainError{Code: "UNAUTHORIZED", Message: "unauthorized access"}
    ErrInvalidInput  = &DomainError{Code: "INVALID_INPUT", Message: "invalid input"}
)
```

**エラーラッピング**
```go
func (s *Service) CreateUser(ctx context.Context, email string) (*User, error) {
    user := &User{Email: email}

    if err := s.repo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    return user, nil
}

// 呼び出し側でエラーチェック
if err := service.CreateUser(ctx, email); err != nil {
    if errors.Is(err, ErrEmailAlreadyExists) {
        // 特定のエラーに対する処理
    }
}
```

### 3. 設定管理

**環境変数とデフォルト値**
```go
type Config struct {
    Port         int
    DatabaseURL  string
    CacheTimeout time.Duration
}

func LoadConfig() (*Config, error) {
    return &Config{
        Port:         getEnvAsInt("PORT", 8080),
        DatabaseURL:  getEnv("DATABASE_URL", ""),
        CacheTimeout: getEnvAsDuration("CACHE_TIMEOUT", 10*time.Minute),
    }, nil
}

func getEnv(key, defaultValue string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return defaultValue
}
```

### 4. ロギング

**構造化ロギング**
```go
type Logger interface {
    Info(msg string, fields ...Field)
    Error(msg string, err error, fields ...Field)
}

type Field struct {
    Key   string
    Value interface{}
}

// 使用例
logger.Info("user created",
    Field{Key: "user_id", Value: user.ID},
    Field{Key: "email", Value: user.Email},
)

logger.Error("failed to create user", err,
    Field{Key: "email", Value: email},
)
```

### 5. テスタビリティ

**テーブル駆動テスト**
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "user@example.com",
            wantErr: false,
        },
        {
            name:    "missing @",
            email:   "userexample.com",
            wantErr: true,
        },
        {
            name:    "empty email",
            email:   "",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("validateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

## よくある落とし穴とその回避方法

### 1. 過度な抽象化

**問題**
```go
// 過度に抽象化された不必要なレイヤー
type DataAccessObject interface {
    Execute(query Query) (Result, error)
}

type Query interface {
    GetSQL() string
    GetParams() []interface{}
}

type Result interface {
    GetData() interface{}
}
```

**解決策**
```go
// 必要な抽象化のみ
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id string) (*User, error)
}
```

### 2. 神オブジェクト（God Object）

**問題**
```go
// すべての責任を持つ巨大なサービス
type ApplicationService struct {
    // 数十のフィールド...
}

func (s *ApplicationService) DoEverything() {
    // 数百行のコード...
}
```

**解決策**
```go
// 責任を分離
type UserService struct {
    repo UserRepository
}

type AuthService struct {
    userRepo UserRepository
    tokenGen TokenGenerator
}

type NotificationService struct {
    emailSender EmailSender
}
```

### 3. 循環依存

**問題**
```go
// package A
import "projectb"

type ServiceA struct {
    serviceB *b.ServiceB
}

// package B
import "projecta"

type ServiceB struct {
    serviceA *a.ServiceA  // 循環依存
}
```

**解決策**
```go
// インターフェースを使って依存を切る
// package A
type ServiceBInterface interface {
    DoSomething() error
}

type ServiceA struct {
    serviceB ServiceBInterface
}

// package B
// import "projecta" は不要

type ServiceB struct {
    // ServiceAに直接依存しない
}
```

### 4. プリマチュアオプティマイゼーション

**問題**
最初から複雑なキャッシング、シャーディング、非同期処理を導入する。

**解決策**
まずはシンプルに実装し、パフォーマンステストで問題が確認されてから最適化する。

```go
// まずはシンプルに
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    return s.repo.FindByID(ctx, id)
}

// 必要になったらキャッシングを追加
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    if user, err := s.cache.Get(ctx, id); err == nil {
        return user, nil
    }

    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    s.cache.Set(ctx, id, user, 10*time.Minute)
    return user, nil
}
```

## リファレンス

### バンドルされた参考資料

- `references/design_patterns.md` - 主要な設計パターンの詳細とGoでの実装例
- `references/architecture_patterns.md` - アーキテクチャパターンの詳細ガイド
- `references/solid_principles.md` - SOLID原則の深い理解と実践例

### 推奨リソース

**書籍**
- 『Clean Architecture』Robert C. Martin
- 『Domain-Driven Design』Eric Evans
- 『Design Patterns』Gang of Four
- 『Refactoring』Martin Fowler

**オンラインリソース**
- Go言語の公式ドキュメント
- クリーンアーキテクチャのサンプル実装
- マイクロサービスパターンのカタログ

## 使用上の注意

### このスキルを効果的に使うために

1. **コンテキストを提供する**：現在のアーキテクチャや制約条件を説明してください
2. **具体的な要件を明記する**：スケーラビリティ、メンテナンス性など、優先する品質を伝えてください
3. **既存コードを共有する**：リファクタリングの場合は、現在のコード構造を見せてください
4. **トレードオフを確認する**：設計の選択には必ずトレードオフがあります。要件に応じて最適な選択を検討します

### 設計の進め方

1. **要件の理解**：機能要件と非機能要件を明確にする
2. **ドメインモデリング**：ビジネスドメインを理解し、適切な抽象化を見つける
3. **アーキテクチャ選択**：システムの特性に合ったアーキテクチャを選ぶ
4. **パターンの適用**：問題に合った設計パターンを適用する
5. **反復的改善**：フィードバックを得ながら設計を改善する

設計は一度で完璧にする必要はありません。継続的に改善していくプロセスです。
