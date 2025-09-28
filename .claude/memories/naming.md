# 命名規約（コーディング指針）

このファイルは、Claude Codeがコードを作成・編集する際に遵守すべき命名規約を定義します。

## ファイル命名規約

### 基本原則

1. **小文字繋ぎの名前を使用**
   - スネークケース（アンダースコア）とケバブケース（ハイフン区切り）は使用禁止

2. **ディレクトリ名重複回避**
   - ディレクトリ名で分かる内容はファイル名に含めない
   - 簡潔で意味のある名前を付ける

### レイヤー別命名規則

#### UseCase層（application/）
- 動作を表す動詞 + 対象オブジェクト
- 例：`getuser.go`, `createuser.go`, `updateproduct.go`

✅ **適切な例**
```
usecases/
├── getuser.go          # ユーザー取得ユースケース
├── createuser.go       # ユーザー作成ユースケース
├── updateuser.go       # ユーザー更新ユースケース
├── deleteuser.go       # ユーザー削除ユースケース
└── listusers.go        # ユーザー一覧取得ユースケース
```

❌ **不適切な例**
```
usecase/
├── get_user.go                    # スネークケース使用
├── user_usecase.go               # ディレクトリ名重複
├── getuserusecase.go             # ディレクトリ名重複
└── user-service.go               # 役割が曖昧
```

#### Controller層（controller/）
- オブジェクト名のみ
- 例：`user.go`, `product.go`

#### Domain層（domain/）
- **entities/**: オブジェクト名の単数形
  - 例：`user.go`, `product.go`
- **repositories/**: オブジェクト名のみ
  - 例：`user.go`, `product.go`

#### Infrastructure層（infrastructure/）
- **repositories/**: オブジェクト名のみ
  - 例：`user.go`

### テストファイル命名規則
- 対象ファイル名 + `_test.go`
- 例：`getuser_test.go`, `user_test.go`

## パッケージ命名規約

### 基本原則
1. **小文字のみ使用**
2. **短く、意味のある名前**
3. **複数形ではなく単数形を使用**（特別な場合を除く）

### 例
```go
package application  // ✅ 適切
package applications // ❌ 複数形
package app         // ❌ 略語（意味が不明確）
package Application // ❌ 大文字使用
```

## インポートパス

### 標準的なインポートパス
```go
import (
    "go-api-server-sample/internal/domain/entities"
    "go-api-server-sample/internal/domain/repositories"
    "go-api-server-sample/cmd/api-server/internal/application"
)
```

## Go言語標準命名規則

### 関数・メソッド名
- **public**: CamelCase（大文字始まり）
- **private**: camelCase（小文字始まり）

### 構造体名
- CamelCase（大文字始まり）

### 定数名
- CamelCase（大文字始まり）

### 変数名
- **public**: CamelCase（大文字始まり）
- **private**: camelCase（小文字始まり）

## コードレビューポイント

### 必須チェック項目
1. ファイル名にアンダースコアが含まれていないか
2. ディレクトリ名とファイル名が重複していないか
3. 1ユースケース1ファイルの原則が守られているか
4. 適切なレイヤーに配置されているか
5. パッケージ名が小文字単数形になっているか
6. Go言語標準の命名規則に従っているか
