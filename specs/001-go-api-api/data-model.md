# Data Model Specification

## User Entity

### Entity Definition
Userエンティティは、システムにおけるユーザーを表現するドメインエンティティです。基本的な識別情報と監査用のタイムスタンプを含みます。

### Fields

| Field | Type | Constraints | Description |
|-------|------|-------------|-------------|
| ID | uint | Primary Key, Auto Increment | ユーザーの一意識別子 |
| Name | string | Required, Max 100 chars | ユーザーの表示名 |
| Email | string | Required, Unique, Valid email format | ユーザーのメールアドレス |
| CreatedAt | time.Time | Auto-set on creation | レコード作成日時 |
| UpdatedAt | time.Time | Auto-set on update | レコード最終更新日時 |

### Go Struct Definition

```go
package entities

import (
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"size:100;not null" json:"name" validate:"required,max=100"`
    Email     string         `gorm:"size:255;uniqueIndex;not null" json:"email" validate:"required,email"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### Validation Rules

#### Name Field
- 必須フィールド
- 最大100文字
- 空文字列不可
- 先頭・末尾の空白は自動削除

#### Email Field
- 必須フィールド
- 有効なメールアドレス形式
- システム内で一意
- 大文字・小文字を区別しない
- 最大255文字

### Business Rules

#### User Creation
- 新規ユーザー作成時、EmailとNameは必須
- 同一メールアドレスでの重複登録は不可
- CreatedAtとUpdatedAtは自動設定

#### User Update
- IDは変更不可
- Emailの変更時は一意性制約を確認
- UpdatedAtは自動更新

#### User Deletion
- 論理削除（Soft Delete）を使用
- DeletedAtフィールドに削除日時を設定
- 削除されたユーザーは検索結果に含まれない

### State Transitions

```
[New] -> [Active] (Creation)
[Active] -> [Active] (Update)
[Active] -> [Deleted] (Soft Delete)
[Deleted] -> [Active] (Restore - future feature)
```

### Database Schema

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE UNIQUE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
```

### Repository Interface

```go
package repositories

import (
    "context"
    "github.com/your-project/internal/domain/entities"
)

type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    GetByID(ctx context.Context, id uint) (*entities.User, error)
    GetByEmail(ctx context.Context, email string) (*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uint) error
    List(ctx context.Context, limit, offset int) ([]*entities.User, error)
    Count(ctx context.Context) (int64, error)
}
```

### Error Conditions

#### Domain Errors
- `ErrUserNotFound`: 指定されたIDのユーザーが存在しない
- `ErrUserEmailExists`: 既に存在するメールアドレス
- `ErrInvalidEmail`: 無効なメールアドレス形式
- `ErrUserNameRequired`: ユーザー名が空または未設定

#### Infrastructure Errors
- `ErrDatabaseConnection`: データベース接続エラー
- `ErrDatabaseConstraint`: データベース制約違反

### Performance Considerations

#### Indexing Strategy
- Primary Key (ID): クラスター化インデックス
- Email: 一意インデックス（論理削除対応）
- DeletedAt: 論理削除フィルタリング用インデックス

#### Query Optimization
- List操作時は必要な分だけのレコード取得（LIMIT/OFFSET）
- 削除されたレコードは自動的にフィルタリング
- Email検索時は一意インデックスを活用