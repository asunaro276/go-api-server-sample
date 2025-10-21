# アーキテクチャパターン詳細リファレンス

このドキュメントでは、ソフトウェアアーキテクチャの主要なパターンについて詳細に説明します。

## 目次

1. [クリーンアーキテクチャ](#クリーンアーキテクチャ)
2. [ヘキサゴナルアーキテクチャ](#ヘキサゴナルアーキテクチャ)
3. [レイヤードアーキテクチャ](#レイヤードアーキテクチャ)
4. [マイクロサービスアーキテクチャ](#マイクロサービスアーキテクチャ)
5. [イベント駆動アーキテクチャ](#イベント駆動アーキテクチャ)
6. [CQRS](#cqrs)

---

## クリーンアーキテクチャ

### 概要

クリーンアーキテクチャは、Robert C. Martin（Uncle Bob）によって提唱されたアーキテクチャパターンです。ビジネスロジックを外部の詳細（フレームワーク、データベース、UIなど）から独立させることを目的としています。

### 主要原則

1. **フレームワーク独立性**：フレームワークに依存しない
2. **テスト可能性**：ビジネスルールはUIやDBなしでテスト可能
3. **UI独立性**：UIを変更してもビジネスルールは影響を受けない
4. **データベース独立性**：データベースを簡単に置き換えられる
5. **外部エージェント独立性**：ビジネスルールは外部世界について何も知らない

### レイヤー構造

```
┌─────────────────────────────────────────────┐
│         Infrastructure Layer                │
│    (Frameworks, Drivers, External Systems)  │
├─────────────────────────────────────────────┤
│         Interface Adapters Layer            │
│    (Controllers, Presenters, Gateways)      │
├─────────────────────────────────────────────┤
│         Application Business Rules          │
│           (Use Cases)                       │
├─────────────────────────────────────────────┤
│       Enterprise Business Rules             │
│           (Entities)                        │
└─────────────────────────────────────────────┘

依存の方向：外側 → 内側
```

### Go言語での実装

**ディレクトリ構造**

```
project/
├── domain/                      # Enterprise Business Rules
│   ├── entity/                 # エンティティ
│   │   ├── user.go
│   │   └── post.go
│   ├── repository/             # リポジトリインターフェース
│   │   ├── user_repository.go
│   │   └── post_repository.go
│   └── service/                # ドメインサービス
│       └── user_service.go
│
├── usecase/                    # Application Business Rules
│   ├── user/
│   │   ├── create_user.go
│   │   ├── get_user.go
│   │   └── update_user.go
│   └── post/
│       ├── create_post.go
│       └── list_posts.go
│
├── interface/                  # Interface Adapters
│   ├── handler/               # HTTPハンドラー
│   │   ├── user_handler.go
│   │   └── post_handler.go
│   ├── presenter/             # プレゼンター
│   │   ├── user_presenter.go
│   │   └── post_presenter.go
│   └── repository/            # リポジトリ実装
│       ├── postgres/
│       │   ├── user_repository.go
│       │   └── post_repository.go
│       └── inmemory/
│           └── user_repository.go
│
└── infrastructure/             # Infrastructure
    ├── database/              # データベース接続
    │   └── postgres.go
    ├── server/                # HTTPサーバー
    │   └── server.go
    ├── config/                # 設定
    │   └── config.go
    └── external/              # 外部API
        └── email_client.go
```

**実装例**

```go
// domain/entity/user.go
package entity

import (
    "errors"
    "time"
)

type User struct {
    ID        string
    Email     string
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func NewUser(email, name string) (*User, error) {
    if email == "" {
        return nil, errors.New("email is required")
    }
    if name == "" {
        return nil, errors.New("name is required")
    }

    return &User{
        Email:     email,
        Name:      name,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }, nil
}

func (u *User) UpdateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    u.Email = email
    u.UpdatedAt = time.Now()
    return nil
}
```

```go
// domain/repository/user_repository.go
package repository

import (
    "context"
    "project/domain/entity"
)

type UserRepository interface {
    Create(ctx context.Context, user *entity.User) error
    FindByID(ctx context.Context, id string) (*entity.User, error)
    FindByEmail(ctx context.Context, email string) (*entity.User, error)
    Update(ctx context.Context, user *entity.User) error
    Delete(ctx context.Context, id string) error
}
```

```go
// usecase/user/create_user.go
package user

import (
    "context"
    "errors"
    "project/domain/entity"
    "project/domain/repository"
)

type CreateUserInput struct {
    Email string
    Name  string
}

type CreateUserOutput struct {
    User *entity.User
}

type CreateUserUseCase struct {
    userRepo repository.UserRepository
}

func NewCreateUserUseCase(userRepo repository.UserRepository) *CreateUserUseCase {
    return &CreateUserUseCase{userRepo: userRepo}
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // ビジネスルール：同じメールアドレスのユーザーは登録できない
    existingUser, err := uc.userRepo.FindByEmail(ctx, input.Email)
    if err == nil && existingUser != nil {
        return nil, errors.New("user with this email already exists")
    }

    // エンティティを作成
    user, err := entity.NewUser(input.Email, input.Name)
    if err != nil {
        return nil, err
    }

    // リポジトリを通じて永続化
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return &CreateUserOutput{User: user}, nil
}
```

```go
// interface/handler/user_handler.go
package handler

import (
    "encoding/json"
    "net/http"
    "project/usecase/user"
)

type UserHandler struct {
    createUserUC *user.CreateUserUseCase
}

func NewUserHandler(createUserUC *user.CreateUserUseCase) *UserHandler {
    return &UserHandler{createUserUC: createUserUC}
}

type CreateUserRequest struct {
    Email string `json:"email"`
    Name  string `json:"name"`
}

type CreateUserResponse struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    output, err := h.createUserUC.Execute(r.Context(), user.CreateUserInput{
        Email: req.Email,
        Name:  req.Name,
    })
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    resp := CreateUserResponse{
        ID:        output.User.ID,
        Email:     output.User.Email,
        Name:      output.User.Name,
        CreatedAt: output.User.CreatedAt.Format("2006-01-02T15:04:05Z"),
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(resp)
}
```

```go
// interface/repository/postgres/user_repository.go
package postgres

import (
    "context"
    "database/sql"
    "project/domain/entity"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    query := `
        INSERT INTO users (id, email, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `
    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Email, user.Name, user.CreatedAt, user.UpdatedAt)
    return err
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*entity.User, error) {
    query := `SELECT id, email, name, created_at, updated_at FROM users WHERE id = $1`

    var user entity.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    query := `SELECT id, email, name, created_at, updated_at FROM users WHERE email = $1`

    var user entity.User
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID, &user.Email, &user.Name, &user.CreatedAt, &user.UpdatedAt)

    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
    query := `
        UPDATE users
        SET email = $2, name = $3, updated_at = $4
        WHERE id = $1
    `
    _, err := r.db.ExecContext(ctx, query,
        user.ID, user.Email, user.Name, user.UpdatedAt)
    return err
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
    query := `DELETE FROM users WHERE id = $1`
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
```

### テスト戦略

```go
// usecase/user/create_user_test.go
package user_test

import (
    "context"
    "testing"
    "project/domain/entity"
    "project/usecase/user"
)

// モックリポジトリ
type MockUserRepository struct {
    users map[string]*entity.User
}

func NewMockUserRepository() *MockUserRepository {
    return &MockUserRepository{
        users: make(map[string]*entity.User),
    }
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) error {
    m.users[user.ID] = user
    return nil
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
    for _, user := range m.users {
        if user.Email == email {
            return user, nil
        }
    }
    return nil, nil
}

func TestCreateUserUseCase_Execute(t *testing.T) {
    // Arrange
    mockRepo := NewMockUserRepository()
    uc := user.NewCreateUserUseCase(mockRepo)

    // Act
    output, err := uc.Execute(context.Background(), user.CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    })

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if output.User.Email != "test@example.com" {
        t.Errorf("Expected email test@example.com, got %s", output.User.Email)
    }
}

func TestCreateUserUseCase_DuplicateEmail(t *testing.T) {
    // Arrange
    mockRepo := NewMockUserRepository()
    uc := user.NewCreateUserUseCase(mockRepo)

    // 最初のユーザーを作成
    uc.Execute(context.Background(), user.CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    })

    // Act: 同じメールアドレスで再度作成を試みる
    _, err := uc.Execute(context.Background(), user.CreateUserInput{
        Email: "test@example.com",
        Name:  "Another User",
    })

    // Assert
    if err == nil {
        t.Fatal("Expected error for duplicate email, got nil")
    }
}
```

---

## ヘキサゴナルアーキテクチャ

### 概要

ヘキサゴナルアーキテクチャ（ポート&アダプターアーキテクチャ）は、Alistair Cockburnによって提唱されました。アプリケーションを中心に置き、外部システムとの接続をポートとアダプターで抽象化します。

### 主要概念

- **アプリケーションコア**：ビジネスロジックを含む
- **ポート**：アプリケーションの境界を定義するインターフェース
  - インバウンドポート：外部からアプリケーションへの入力
  - アウトバウンドポート：アプリケーションから外部への出力
- **アダプター**：ポートの具体的な実装
  - プライマリアダプター：アプリケーションを駆動する（HTTPハンドラー、CLIなど）
  - セカンダリアダプター：アプリケーションによって駆動される（データベース、外部APIなど）

### 構造図

```
        ┌───────────────────────────────┐
        │   Primary Adapters            │
        │   (HTTP, CLI, gRPC)           │
        └───────────┬───────────────────┘
                    │
        ┌───────────▼───────────────────┐
        │   Inbound Ports               │
        │   (Use Case Interfaces)       │
        ├───────────────────────────────┤
        │                               │
        │   Application Core            │
        │   (Business Logic)            │
        │                               │
        ├───────────────────────────────┤
        │   Outbound Ports              │
        │   (Repository, External APIs) │
        └───────────┬───────────────────┘
                    │
        ┌───────────▼───────────────────┐
        │   Secondary Adapters          │
        │   (PostgreSQL, Redis, SMTP)   │
        └───────────────────────────────┘
```

### Go言語での実装

```go
// アプリケーションコア
package core

import (
    "context"
    "errors"
)

// User はドメインエンティティ
type User struct {
    ID       string
    Email    string
    Password string
}

// インバウンドポート：ユーザーサービスのインターフェース
type UserService interface {
    RegisterUser(ctx context.Context, email, password string) (*User, error)
    AuthenticateUser(ctx context.Context, email, password string) (*User, error)
}

// アウトバウンドポート：ユーザーリポジトリのインターフェース
type UserRepository interface {
    Save(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
}

// アウトバウンドポート：パスワードハッシャーのインターフェース
type PasswordHasher interface {
    Hash(password string) (string, error)
    Compare(hashedPassword, password string) error
}

// アウトバウンドポート：メール送信のインターフェース
type EmailSender interface {
    SendWelcomeEmail(ctx context.Context, email string) error
}

// UserServiceImpl はユーザーサービスの実装
type UserServiceImpl struct {
    userRepo       UserRepository
    passwordHasher PasswordHasher
    emailSender    EmailSender
}

func NewUserService(
    userRepo UserRepository,
    passwordHasher PasswordHasher,
    emailSender EmailSender,
) *UserServiceImpl {
    return &UserServiceImpl{
        userRepo:       userRepo,
        passwordHasher: passwordHasher,
        emailSender:    emailSender,
    }
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, email, password string) (*User, error) {
    // ビジネスルール：既存ユーザーチェック
    existing, _ := s.userRepo.FindByEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("user already exists")
    }

    // パスワードハッシュ化
    hashedPassword, err := s.passwordHasher.Hash(password)
    if err != nil {
        return nil, err
    }

    // ユーザー作成
    user := &User{
        Email:    email,
        Password: hashedPassword,
    }

    // 保存
    if err := s.userRepo.Save(ctx, user); err != nil {
        return nil, err
    }

    // ウェルカムメール送信（非同期推奨だが簡略化のため同期）
    s.emailSender.SendWelcomeEmail(ctx, email)

    return user, nil
}

func (s *UserServiceImpl) AuthenticateUser(ctx context.Context, email, password string) (*User, error) {
    user, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if err := s.passwordHasher.Compare(user.Password, password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    return user, nil
}
```

```go
// プライマリアダプター：HTTPハンドラー
package adapter

import (
    "encoding/json"
    "net/http"
    "project/core"
)

type HTTPAdapter struct {
    userService core.UserService
}

func NewHTTPAdapter(userService core.UserService) *HTTPAdapter {
    return &HTTPAdapter{userService: userService}
}

type RegisterRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}

func (a *HTTPAdapter) HandleRegister(w http.ResponseWriter, r *http.Request) {
    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    user, err := a.userService.RegisterUser(r.Context(), req.Email, req.Password)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "id":    user.ID,
        "email": user.Email,
    })
}
```

```go
// セカンダリアダプター：PostgreSQLリポジトリ
package adapter

import (
    "context"
    "database/sql"
    "project/core"
)

type PostgreSQLAdapter struct {
    db *sql.DB
}

func NewPostgreSQLAdapter(db *sql.DB) *PostgreSQLAdapter {
    return &PostgreSQLAdapter{db: db}
}

func (a *PostgreSQLAdapter) Save(ctx context.Context, user *core.User) error {
    query := `INSERT INTO users (id, email, password) VALUES ($1, $2, $3)`
    _, err := a.db.ExecContext(ctx, query, user.ID, user.Email, user.Password)
    return err
}

func (a *PostgreSQLAdapter) FindByEmail(ctx context.Context, email string) (*core.User, error) {
    query := `SELECT id, email, password FROM users WHERE email = $1`

    var user core.User
    err := a.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Password)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}
```

```go
// セカンダリアダプター：bcryptパスワードハッシャー
package adapter

import "golang.org/x/crypto/bcrypt"

type BcryptAdapter struct{}

func NewBcryptAdapter() *BcryptAdapter {
    return &BcryptAdapter{}
}

func (a *BcryptAdapter) Hash(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func (a *BcryptAdapter) Compare(hashedPassword, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
```

```go
// セカンダリアダプター：SMTPメール送信
package adapter

import (
    "context"
    "fmt"
)

type SMTPAdapter struct {
    host     string
    port     int
    username string
    password string
}

func NewSMTPAdapter(host string, port int, username, password string) *SMTPAdapter {
    return &SMTPAdapter{
        host:     host,
        port:     port,
        username: username,
        password: password,
    }
}

func (a *SMTPAdapter) SendWelcomeEmail(ctx context.Context, email string) error {
    // 実際のSMTP実装
    fmt.Printf("Sending welcome email to %s\n", email)
    return nil
}
```

---

## レイヤードアーキテクチャ

### 概要

レイヤードアーキテクチャは、最も一般的なアーキテクチャパターンの1つです。アプリケーションを水平方向のレイヤーに分割します。

### レイヤー構成

```
┌─────────────────────────────────┐
│  Presentation Layer             │  ← UI, API, Controllers
├─────────────────────────────────┤
│  Business Logic Layer           │  ← Services, Domain Logic
├─────────────────────────────────┤
│  Data Access Layer              │  ← Repositories, DAOs
├─────────────────────────────────┤
│  Database Layer                 │  ← Database
└─────────────────────────────────┘

依存の方向：上から下へ
```

### Go言語での実装

```go
// プレゼンテーション層
package presentation

import (
    "encoding/json"
    "net/http"
    "project/business"
)

type UserController struct {
    userService *business.UserService
}

func NewUserController(userService *business.UserService) *UserController {
    return &UserController{userService: userService}
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")

    user, err := c.userService.GetUserByID(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }

    json.NewEncoder(w).Encode(user)
}
```

```go
// ビジネスロジック層
package business

import (
    "context"
    "errors"
    "project/data"
)

type User struct {
    ID    string
    Email string
    Name  string
}

type UserService struct {
    userRepo *data.UserRepository
}

func NewUserService(userRepo *data.UserRepository) *UserService {
    return &UserService{userRepo: userRepo}
}

func (s *UserService) GetUserByID(ctx context.Context, id string) (*User, error) {
    // ビジネスロジック
    if id == "" {
        return nil, errors.New("user ID is required")
    }

    user, err := s.userRepo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    return user, nil
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*User, error) {
    // バリデーション
    if email == "" || name == "" {
        return nil, errors.New("email and name are required")
    }

    // 重複チェック
    existing, _ := s.userRepo.FindByEmail(ctx, email)
    if existing != nil {
        return nil, errors.New("user already exists")
    }

    user := &User{
        Email: email,
        Name:  name,
    }

    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    return user, nil
}
```

```go
// データアクセス層
package data

import (
    "context"
    "database/sql"
    "project/business"
)

type UserRepository struct {
    db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*business.User, error) {
    query := `SELECT id, email, name FROM users WHERE id = $1`

    var user business.User
    err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Email, &user.Name)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*business.User, error) {
    query := `SELECT id, email, name FROM users WHERE email = $1`

    var user business.User
    err := r.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.Name)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *UserRepository) Create(ctx context.Context, user *business.User) error {
    query := `INSERT INTO users (id, email, name) VALUES ($1, $2, $3)`
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name)
    return err
}
```

---

## マイクロサービスアーキテクチャ

### 概要

マイクロサービスアーキテクチャは、アプリケーションを小さく独立したサービスの集合として構築します。各サービスは独自のプロセスで実行され、軽量なメカニズム（通常はHTTP REST API）で通信します。

### 主要原則

1. **単一責任**：各サービスは1つのビジネス機能を担当
2. **自律性**：独立してデプロイ・スケール可能
3. **分散データ**：各サービスは独自のデータベースを持つ
4. **障害の分離**：1つのサービスの障害が他に波及しない
5. **技術の多様性**：サービスごとに異なる技術スタックを選択可能

### サービス分割の例

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│  User Service   │    │  Order Service  │    │ Payment Service │
│                 │    │                 │    │                 │
│  - User CRUD    │    │  - Create Order │    │  - Process Pay  │
│  - Auth         │    │  - List Orders  │    │  - Refunds      │
│                 │    │  - Update Status│    │                 │
│  [PostgreSQL]   │    │  [MongoDB]      │    │  [PostgreSQL]   │
└────────┬────────┘    └────────┬────────┘    └────────┬────────┘
         │                      │                       │
         └──────────────────────┴───────────────────────┘
                                │
                        ┌───────▼────────┐
                        │  API Gateway   │
                        └────────────────┘
```

### 実装パターン

**サービス間通信**

```go
// HTTP RESTベースの通信
package client

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "net/http"
)

type UserServiceClient struct {
    baseURL string
    client  *http.Client
}

func NewUserServiceClient(baseURL string) *UserServiceClient {
    return &UserServiceClient{
        baseURL: baseURL,
        client:  &http.Client{},
    }
}

func (c *UserServiceClient) GetUser(ctx context.Context, userID string) (*User, error) {
    url := fmt.Sprintf("%s/users/%s", c.baseURL, userID)

    req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
    if err != nil {
        return nil, err
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }

    var user User
    if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
        return nil, err
    }

    return &user, nil
}
```

**サービスディスカバリー**

```go
package discovery

import (
    "context"
    "errors"
)

type ServiceRegistry interface {
    Register(ctx context.Context, service Service) error
    Deregister(ctx context.Context, serviceID string) error
    Discover(ctx context.Context, serviceName string) ([]Service, error)
}

type Service struct {
    ID      string
    Name    string
    Address string
    Port    int
    Tags    []string
}

// Consulを使った実装例
type ConsulRegistry struct {
    // Consulクライアント
}

func (r *ConsulRegistry) Register(ctx context.Context, service Service) error {
    // Consulにサービスを登録
    return nil
}

func (r *ConsulRegistry) Discover(ctx context.Context, serviceName string) ([]Service, error) {
    // Consulからサービスを検索
    return nil, nil
}
```

---

## イベント駆動アーキテクチャ

### 概要

イベント駆動アーキテクチャでは、システムの状態変化をイベントとして表現し、コンポーネント間の疎結合な通信を実現します。

### パターン

**イベントソーシング + CQRS**

```go
// イベント
package event

import "time"

type Event interface {
    EventType() string
    AggregateID() string
    OccurredAt() time.Time
}

type UserRegistered struct {
    UserID    string
    Email     string
    Name      string
    Timestamp time.Time
}

func (e UserRegistered) EventType() string    { return "UserRegistered" }
func (e UserRegistered) AggregateID() string  { return e.UserID }
func (e UserRegistered) OccurredAt() time.Time { return e.Timestamp }

// イベントストア
type EventStore interface {
    Save(ctx context.Context, events []Event) error
    Load(ctx context.Context, aggregateID string) ([]Event, error)
}

// イベントバス
type EventBus interface {
    Publish(ctx context.Context, event Event) error
    Subscribe(eventType string, handler EventHandler) error
}

type EventHandler func(ctx context.Context, event Event) error
```

---

## CQRS

Command Query Responsibility Segregation（コマンドクエリ責務分離）は、読み取りと書き込みのモデルを分離するパターンです。

```go
// コマンド側（書き込み）
package command

type CreateUserCommand struct {
    Email string
    Name  string
}

type CommandHandler interface {
    Handle(ctx context.Context, cmd interface{}) error
}

// クエリ側（読み取り）
package query

type GetUserQuery struct {
    UserID string
}

type UserDTO struct {
    ID    string
    Email string
    Name  string
}

type QueryHandler interface {
    Handle(ctx context.Context, query interface{}) (interface{}, error)
}
```

このリファレンスでは主要なアーキテクチャパターンを説明しました。実際のプロジェクトでは、これらのパターンを組み合わせて使用することが一般的です。
