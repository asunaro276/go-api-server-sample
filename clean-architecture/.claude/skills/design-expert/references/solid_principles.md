# SOLID原則詳細リファレンス

このドキュメントでは、オブジェクト指向設計の基本原則であるSOLID原則について、詳細な説明とGo言語での実践例を提供します。

## 目次

1. [Single Responsibility Principle（単一責任の原則）](#single-responsibility-principle)
2. [Open/Closed Principle（開放閉鎖の原則）](#openclosed-principle)
3. [Liskov Substitution Principle（リスコフの置換原則）](#liskov-substitution-principle)
4. [Interface Segregation Principle（インターフェース分離の原則）](#interface-segregation-principle)
5. [Dependency Inversion Principle（依存性逆転の原則）](#dependency-inversion-principle)

---

## Single Responsibility Principle

**定義**：クラス（またはモジュール）は、変更する理由を1つだけ持つべきである。

### 原則の理解

単一責任の原則は、各コンポーネントが1つの明確な責任を持ち、その責任に関連する変更のみを受けるべきであることを意味します。

### 悪い例：複数の責任を持つ構造体

```go
package bad

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "net/smtp"
)

// UserManager は複数の責任を持っている
type UserManager struct {
    db *sql.DB
}

// ユーザーの作成（ビジネスロジック）
func (m *UserManager) CreateUser(email, name string) error {
    // バリデーション
    if email == "" || name == "" {
        return fmt.Errorf("email and name are required")
    }

    // データベース保存（データアクセス）
    query := "INSERT INTO users (email, name) VALUES (?, ?)"
    _, err := m.db.Exec(query, email, name)
    if err != nil {
        return err
    }

    // メール送信（外部通信）
    m.sendWelcomeEmail(email)

    // ログ記録（ロギング）
    m.logUserCreation(email)

    return nil
}

func (m *UserManager) sendWelcomeEmail(email string) error {
    // SMTP設定（設定管理）
    smtpHost := "smtp.example.com"
    smtpPort := "587"

    // メール送信処理
    auth := smtp.PlainAuth("", "sender@example.com", "password", smtpHost)
    msg := []byte("Welcome to our service!")
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, "sender@example.com", []string{email}, msg)
}

func (m *UserManager) logUserCreation(email string) {
    fmt.Printf("User created: %s\n", email)
}

// ユーザーのエクスポート（データ変換）
func (m *UserManager) ExportUsersToJSON() ([]byte, error) {
    rows, err := m.db.Query("SELECT email, name FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []map[string]string
    for rows.Next() {
        var email, name string
        rows.Scan(&email, &name)
        users = append(users, map[string]string{"email": email, "name": name})
    }

    return json.Marshal(users)
}
```

**問題点**
- UserManagerは以下の複数の責任を持っている：
  1. ビジネスロジック（ユーザー作成）
  2. データアクセス（データベース操作）
  3. 外部通信（メール送信）
  4. ロギング
  5. 設定管理
  6. データ変換（JSON エクスポート）

- 変更の理由が多すぎる：
  - データベーススキーマ変更
  - メールプロバイダー変更
  - ログ形式変更
  - エクスポート形式変更

### 良い例：責任の分離

```go
package good

import (
    "context"
    "encoding/json"
)

// User はドメインエンティティ（1つの責任：ユーザーデータの表現）
type User struct {
    ID    string
    Email string
    Name  string
}

// UserRepository はデータアクセスの責任のみ
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByEmail(ctx context.Context, email string) (*User, error)
    FindAll(ctx context.Context) ([]*User, error)
}

// UserService はビジネスロジックの責任のみ
type UserService struct {
    repo         UserRepository
    emailService EmailService
    logger       Logger
}

func NewUserService(repo UserRepository, emailService EmailService, logger Logger) *UserService {
    return &UserService{
        repo:         repo,
        emailService: emailService,
        logger:       logger,
    }
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*User, error) {
    // バリデーション
    if email == "" || name == "" {
        return nil, fmt.Errorf("email and name are required")
    }

    // ビジネスルール：重複チェック
    existing, _ := s.repo.FindByEmail(ctx, email)
    if existing != nil {
        return nil, fmt.Errorf("user already exists")
    }

    // ユーザー作成
    user := &User{
        Email: email,
        Name:  name,
    }

    // 保存
    if err := s.repo.Create(ctx, user); err != nil {
        return nil, err
    }

    // ウェルカムメール送信（依存サービスに委譲）
    if err := s.emailService.SendWelcomeEmail(ctx, email); err != nil {
        s.logger.Error("Failed to send welcome email", "email", email, "error", err)
        // メール送信失敗はユーザー作成の失敗とはしない
    }

    // ログ記録（依存サービスに委譲）
    s.logger.Info("User created", "email", email)

    return user, nil
}

// EmailService はメール送信の責任のみ
type EmailService interface {
    SendWelcomeEmail(ctx context.Context, email string) error
}

// Logger はロギングの責任のみ
type Logger interface {
    Info(msg string, keysAndValues ...interface{})
    Error(msg string, keysAndValues ...interface{})
}

// UserExporter はデータエクスポートの責任のみ
type UserExporter struct {
    repo UserRepository
}

func NewUserExporter(repo UserRepository) *UserExporter {
    return &UserExporter{repo: repo}
}

func (e *UserExporter) ExportToJSON(ctx context.Context) ([]byte, error) {
    users, err := e.repo.FindAll(ctx)
    if err != nil {
        return nil, err
    }

    return json.Marshal(users)
}
```

### 具体的な実装例

```go
// PostgreSQLUserRepository はPostgreSQLを使ったリポジトリ実装
type PostgreSQLUserRepository struct {
    db *sql.DB
}

func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *User) error {
    query := "INSERT INTO users (id, email, name) VALUES ($1, $2, $3)"
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Email, user.Name)
    return err
}

// SMTPEmailService はSMTPを使ったメールサービス実装
type SMTPEmailService struct {
    host     string
    port     string
    username string
    password string
}

func (s *SMTPEmailService) SendWelcomeEmail(ctx context.Context, email string) error {
    // SMTP送信ロジック
    return nil
}

// StructuredLogger は構造化ロギングの実装
type StructuredLogger struct {
    // ロギングライブラリ（zap, logrusなど）
}

func (l *StructuredLogger) Info(msg string, keysAndValues ...interface{}) {
    // ログ出力
}

func (l *StructuredLogger) Error(msg string, keysAndValues ...interface{}) {
    // エラーログ出力
}
```

### メリット

1. **保守性の向上**：各コンポーネントが明確な責任を持つため、変更範囲が限定される
2. **テスト容易性**：責任が分離されているため、モックやスタブを使った単体テストが容易
3. **再利用性**：単一責任のコンポーネントは他の文脈でも再利用しやすい
4. **理解しやすさ**：各コンポーネントの目的が明確

---

## Open/Closed Principle

**定義**：ソフトウェアのエンティティ（クラス、モジュール、関数など）は、拡張に対して開いていて、修正に対して閉じているべきである。

### 原則の理解

新しい機能を追加する際、既存のコードを変更せずに（閉じている）、新しいコードを追加することで（開いている）実現できるようにすべきです。

### 悪い例：修正が必要な構造

```go
package bad

type PaymentProcessor struct{}

func (p *PaymentProcessor) ProcessPayment(method string, amount float64) error {
    switch method {
    case "credit_card":
        return p.processCreditCard(amount)
    case "paypal":
        return p.processPayPal(amount)
    case "bank_transfer":
        return p.processBankTransfer(amount)
    default:
        return fmt.Errorf("unsupported payment method: %s", method)
    }
}

func (p *PaymentProcessor) processCreditCard(amount float64) error {
    fmt.Printf("Processing credit card payment: %.2f\n", amount)
    return nil
}

func (p *PaymentProcessor) processPayPal(amount float64) error {
    fmt.Printf("Processing PayPal payment: %.2f\n", amount)
    return nil
}

func (p *PaymentProcessor) processBankTransfer(amount float64) error {
    fmt.Printf("Processing bank transfer: %.2f\n", amount)
    return nil
}
```

**問題点**
- 新しい支払い方法を追加するたびに、`ProcessPayment`メソッドを修正する必要がある
- switch文が増え続ける
- 既存コードの変更はバグのリスクを高める

### 良い例：拡張可能な構造

```go
package good

import (
    "context"
    "fmt"
)

// PaymentMethod はすべての支払い方法が実装すべきインターフェース
type PaymentMethod interface {
    Process(ctx context.Context, amount float64) error
    GetName() string
}

// CreditCardPayment はクレジットカード支払い
type CreditCardPayment struct {
    cardNumber string
    cvv        string
}

func NewCreditCardPayment(cardNumber, cvv string) *CreditCardPayment {
    return &CreditCardPayment{
        cardNumber: cardNumber,
        cvv:        cvv,
    }
}

func (c *CreditCardPayment) Process(ctx context.Context, amount float64) error {
    fmt.Printf("Processing credit card payment: %.2f (card ending in %s)\n",
        amount, c.cardNumber[len(c.cardNumber)-4:])
    return nil
}

func (c *CreditCardPayment) GetName() string {
    return "Credit Card"
}

// PayPalPayment はPayPal支払い
type PayPalPayment struct {
    email string
}

func NewPayPalPayment(email string) *PayPalPayment {
    return &PayPalPayment{email: email}
}

func (p *PayPalPayment) Process(ctx context.Context, amount float64) error {
    fmt.Printf("Processing PayPal payment: %.2f (account: %s)\n", amount, p.email)
    return nil
}

func (p *PayPalPayment) GetName() string {
    return "PayPal"
}

// BankTransferPayment は銀行振込
type BankTransferPayment struct {
    accountNumber string
    routingNumber string
}

func NewBankTransferPayment(accountNumber, routingNumber string) *BankTransferPayment {
    return &BankTransferPayment{
        accountNumber: accountNumber,
        routingNumber: routingNumber,
    }
}

func (b *BankTransferPayment) Process(ctx context.Context, amount float64) error {
    fmt.Printf("Processing bank transfer: %.2f (account: %s)\n", amount, b.accountNumber)
    return nil
}

func (b *BankTransferPayment) GetName() string {
    return "Bank Transfer"
}

// PaymentProcessor は修正不要
type PaymentProcessor struct {
    method PaymentMethod
}

func NewPaymentProcessor(method PaymentMethod) *PaymentProcessor {
    return &PaymentProcessor{method: method}
}

func (p *PaymentProcessor) SetMethod(method PaymentMethod) {
    p.method = method
}

func (p *PaymentProcessor) ProcessPayment(ctx context.Context, amount float64) error {
    if p.method == nil {
        return fmt.Errorf("payment method not set")
    }

    fmt.Printf("Using payment method: %s\n", p.method.GetName())
    return p.method.Process(ctx, amount)
}
```

### 新しい支払い方法の追加（既存コードを変更せずに拡張）

```go
// CryptocurrencyPayment は新しい支払い方法（既存コードは一切変更しない）
type CryptocurrencyPayment struct {
    walletAddress string
    cryptocurrency string
}

func NewCryptocurrencyPayment(walletAddress, cryptocurrency string) *CryptocurrencyPayment {
    return &CryptocurrencyPayment{
        walletAddress: walletAddress,
        cryptocurrency: cryptocurrency,
    }
}

func (c *CryptocurrencyPayment) Process(ctx context.Context, amount float64) error {
    fmt.Printf("Processing %s payment: %.8f (wallet: %s)\n",
        c.cryptocurrency, amount, c.walletAddress)
    return nil
}

func (c *CryptocurrencyPayment) GetName() string {
    return c.cryptocurrency
}

// 使用例
func main() {
    processor := NewPaymentProcessor(nil)

    // クレジットカードで支払い
    processor.SetMethod(NewCreditCardPayment("1234-5678-9012-3456", "123"))
    processor.ProcessPayment(context.Background(), 100.00)

    // 新しい支払い方法を追加しても既存コードは変更不要
    processor.SetMethod(NewCryptocurrencyPayment("1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa", "Bitcoin"))
    processor.ProcessPayment(context.Background(), 0.001)
}
```

### 別の例：レポート生成

```go
// ReportGenerator は拡張可能な設計
type ReportGenerator interface {
    Generate(data interface{}) ([]byte, error)
    Format() string
}

type PDFReportGenerator struct{}

func (p *PDFReportGenerator) Generate(data interface{}) ([]byte, error) {
    // PDF生成ロジック
    return []byte("PDF content"), nil
}

func (p *PDFReportGenerator) Format() string {
    return "PDF"
}

type ExcelReportGenerator struct{}

func (e *ExcelReportGenerator) Generate(data interface{}) ([]byte, error) {
    // Excel生成ロジック
    return []byte("Excel content"), nil
}

func (e *ExcelReportGenerator) Format() string {
    return "Excel"
}

// 新しい形式を追加（既存コードは変更しない）
type MarkdownReportGenerator struct{}

func (m *MarkdownReportGenerator) Generate(data interface{}) ([]byte, error) {
    // Markdown生成ロジック
    return []byte("Markdown content"), nil
}

func (m *MarkdownReportGenerator) Format() string {
    return "Markdown"
}

// ReportService は修正不要
type ReportService struct {
    generators map[string]ReportGenerator
}

func NewReportService() *ReportService {
    return &ReportService{
        generators: make(map[string]ReportGenerator),
    }
}

func (s *ReportService) RegisterGenerator(generator ReportGenerator) {
    s.generators[generator.Format()] = generator
}

func (s *ReportService) GenerateReport(format string, data interface{}) ([]byte, error) {
    generator, ok := s.generators[format]
    if !ok {
        return nil, fmt.Errorf("unsupported format: %s", format)
    }

    return generator.Generate(data)
}
```

---

## Liskov Substitution Principle

**定義**：派生型は、その基本型と置換可能でなければならない。

### 原則の理解

サブタイプは、プログラムの正確性を損なうことなく、基本型と置換可能でなければなりません。つまり、インターフェースの契約を守る必要があります。

### 悪い例：契約違反

```go
package bad

type Rectangle struct {
    width  float64
    height float64
}

func (r *Rectangle) SetWidth(width float64) {
    r.width = width
}

func (r *Rectangle) SetHeight(height float64) {
    r.height = height
}

func (r *Rectangle) GetArea() float64 {
    return r.width * r.height
}

// Square は Rectangle を継承（コンポジション）
type Square struct {
    Rectangle
}

// SetWidth は正方形の特性を保つため、両方の辺を変更
func (s *Square) SetWidth(width float64) {
    s.width = width
    s.height = width // 正方形なので高さも変更
}

// SetHeight も同様
func (s *Square) SetHeight(height float64) {
    s.width = height
    s.height = height
}

// 問題のあるコード
func CalculateArea(r *Rectangle) {
    r.SetWidth(5)
    r.SetHeight(10)
    area := r.GetArea()

    // Rectangle を期待すると area = 50
    // Square だと area = 100 （予期しない動作）
    fmt.Printf("Expected: 50, Got: %.0f\n", area)
}
```

**問題点**
- `Square`は`Rectangle`の契約を守っていない
- `SetWidth`を呼んだ後に`SetHeight`を呼ぶと、幅も変わってしまう
- 置換可能性が破られている

### 良い例：契約を守る設計

```go
package good

// Shape は図形の抽象インターフェース
type Shape interface {
    GetArea() float64
    GetPerimeter() float64
}

// Rectangle は長方形
type Rectangle struct {
    width  float64
    height float64
}

func NewRectangle(width, height float64) *Rectangle {
    return &Rectangle{width: width, height: height}
}

func (r *Rectangle) GetArea() float64 {
    return r.width * r.height
}

func (r *Rectangle) GetPerimeter() float64 {
    return 2 * (r.width + r.height)
}

func (r *Rectangle) SetWidth(width float64) {
    r.width = width
}

func (r *Rectangle) SetHeight(height float64) {
    r.height = height
}

// Square は独立した型として定義
type Square struct {
    side float64
}

func NewSquare(side float64) *Square {
    return &Square{side: side}
}

func (s *Square) GetArea() float64 {
    return s.side * s.side
}

func (s *Square) GetPerimeter() float64 {
    return 4 * s.side
}

func (s *Square) SetSide(side float64) {
    s.side = side
}

// CalculateArea は Shape インターフェースを使用
func CalculateArea(shape Shape) float64 {
    return shape.GetArea()
}

// 使用例
func main() {
    rectangle := NewRectangle(5, 10)
    square := NewSquare(5)

    fmt.Printf("Rectangle area: %.0f\n", CalculateArea(rectangle)) // 50
    fmt.Printf("Square area: %.0f\n", CalculateArea(square))       // 25
}
```

### 別の例：データアクセス層

```go
// 悪い例：契約違反
type BadRepository interface {
    FindByID(id string) (*User, error)
    Save(user *User) error
}

type ReadOnlyRepository struct {
    // 読み取り専用リポジトリ
}

func (r *ReadOnlyRepository) FindByID(id string) (*User, error) {
    // 正常に動作
    return &User{}, nil
}

func (r *ReadOnlyRepository) Save(user *User) error {
    // 契約違反：常にエラーを返す
    return errors.New("read-only repository")
}

// 良い例：契約を分離
type Reader interface {
    FindByID(id string) (*User, error)
}

type Writer interface {
    Save(user *User) error
}

type Repository interface {
    Reader
    Writer
}

type ReadOnlyRepository struct {
    // 読み取り専用リポジトリ
}

func (r *ReadOnlyRepository) FindByID(id string) (*User, error) {
    return &User{}, nil
}

// ReadOnlyRepository は Reader のみを実装
// Writer は実装しないので、契約違反は発生しない

type FullRepository struct {
    // 読み書き可能なリポジトリ
}

func (r *FullRepository) FindByID(id string) (*User, error) {
    return &User{}, nil
}

func (r *FullRepository) Save(user *User) error {
    return nil
}
```

---

## Interface Segregation Principle

**定義**：クライアントは、使用しないメソッドへの依存を強制されるべきではない。

### 原則の理解

大きすぎるインターフェースを小さく、焦点を絞ったインターフェースに分割すべきです。

### 悪い例：大きすぎるインターフェース

```go
package bad

// Worker は大きすぎるインターフェース
type Worker interface {
    Work() error
    Eat() error
    Sleep() error
    GetSalary() float64
    TakeLunch() error
    DrinkCoffee() error
}

// HumanWorker はすべてのメソッドを実装できる
type HumanWorker struct {
    name   string
    salary float64
}

func (h *HumanWorker) Work() error {
    fmt.Println("Human working")
    return nil
}

func (h *HumanWorker) Eat() error {
    fmt.Println("Human eating")
    return nil
}

func (h *HumanWorker) Sleep() error {
    fmt.Println("Human sleeping")
    return nil
}

func (h *HumanWorker) GetSalary() float64 {
    return h.salary
}

func (h *HumanWorker) TakeLunch() error {
    fmt.Println("Human taking lunch")
    return nil
}

func (h *HumanWorker) DrinkCoffee() error {
    fmt.Println("Human drinking coffee")
    return nil
}

// RobotWorker は不要なメソッドも実装を強制される
type RobotWorker struct {
    model string
}

func (r *RobotWorker) Work() error {
    fmt.Println("Robot working")
    return nil
}

func (r *RobotWorker) Eat() error {
    // ロボットは食べない - 無意味な実装
    return errors.New("robots don't eat")
}

func (r *RobotWorker) Sleep() error {
    // ロボットは寝ない - 無意味な実装
    return errors.New("robots don't sleep")
}

func (r *RobotWorker) GetSalary() float64 {
    // ロボットに給与はない
    return 0
}

func (r *RobotWorker) TakeLunch() error {
    return errors.New("robots don't take lunch")
}

func (r *RobotWorker) DrinkCoffee() error {
    return errors.New("robots don't drink coffee")
}
```

**問題点**
- `RobotWorker`は不要なメソッド（`Eat`, `Sleep`, `TakeLunch`, `DrinkCoffee`）の実装を強制される
- インターフェースが大きすぎて、すべての実装者に適合しない

### 良い例：小さく焦点を絞ったインターフェース

```go
package good

// Workable は作業能力を表す
type Workable interface {
    Work() error
}

// Eatable は食事能力を表す
type Eatable interface {
    Eat() error
}

// Sleepable は睡眠能力を表す
type Sleepable interface {
    Sleep() error
}

// Salaried は給与を受け取る能力を表す
type Salaried interface {
    GetSalary() float64
}

// Breakable は休憩能力を表す
type Breakable interface {
    TakeLunch() error
    DrinkCoffee() error
}

// HumanWorker は必要なインターフェースのみを実装
type HumanWorker struct {
    name   string
    salary float64
}

func (h *HumanWorker) Work() error {
    fmt.Println("Human working")
    return nil
}

func (h *HumanWorker) Eat() error {
    fmt.Println("Human eating")
    return nil
}

func (h *HumanWorker) Sleep() error {
    fmt.Println("Human sleeping")
    return nil
}

func (h *HumanWorker) GetSalary() float64 {
    return h.salary
}

func (h *HumanWorker) TakeLunch() error {
    fmt.Println("Human taking lunch")
    return nil
}

func (h *HumanWorker) DrinkCoffee() error {
    fmt.Println("Human drinking coffee")
    return nil
}

// RobotWorker は必要なメソッドのみを実装
type RobotWorker struct {
    model string
}

func (r *RobotWorker) Work() error {
    fmt.Println("Robot working")
    return nil
}

// ロボットは Eatable, Sleepable, Salaried, Breakable を実装しない

// 使用例：必要なインターフェースのみを要求
func AssignWork(worker Workable) {
    worker.Work()
}

func PaySalary(employee Salaried) {
    salary := employee.GetSalary()
    fmt.Printf("Paying salary: %.2f\n", salary)
}

func ScheduleLunch(employee Breakable) {
    employee.TakeLunch()
}

func main() {
    human := &HumanWorker{name: "John", salary: 50000}
    robot := &RobotWorker{model: "X-100"}

    // 両方に作業を割り当て
    AssignWork(human)
    AssignWork(robot)

    // 人間にのみ給与を支払い
    PaySalary(human)

    // 人間にのみ昼食をスケジュール
    ScheduleLunch(human)

    // robot は Salaried や Breakable を実装していないため、
    // PaySalary(robot) や ScheduleLunch(robot) はコンパイルエラー
}
```

### 別の例：データアクセス

```go
// 悪い例
type BadRepository interface {
    Create(item interface{}) error
    Read(id string) (interface{}, error)
    Update(item interface{}) error
    Delete(id string) error
    BulkInsert(items []interface{}) error
    Search(query string) ([]interface{}, error)
    Export(format string) ([]byte, error)
    Import(data []byte, format string) error
}

// 良い例：焦点を絞ったインターフェース
type Reader interface {
    Read(id string) (interface{}, error)
}

type Writer interface {
    Create(item interface{}) error
    Update(item interface{}) error
    Delete(id string) error
}

type BulkWriter interface {
    BulkInsert(items []interface{}) error
}

type Searcher interface {
    Search(query string) ([]interface{}, error)
}

type Exporter interface {
    Export(format string) ([]byte, error)
}

type Importer interface {
    Import(data []byte, format string) error
}

// 必要な機能のみを実装
type ReadOnlyRepository struct{}

func (r *ReadOnlyRepository) Read(id string) (interface{}, error) {
    return nil, nil
}

type BasicRepository struct{}

func (r *BasicRepository) Create(item interface{}) error { return nil }
func (r *BasicRepository) Read(id string) (interface{}, error) { return nil, nil }
func (r *BasicRepository) Update(item interface{}) error { return nil }
func (r *BasicRepository) Delete(id string) error { return nil }

// BasicRepository は Reader と Writer のみを実装
// BulkWriter, Searcher, Exporter, Importer は実装しない
```

---

## Dependency Inversion Principle

**定義**：
1. 上位モジュールは下位モジュールに依存してはならない。両方とも抽象に依存すべきである。
2. 抽象は詳細に依存してはならない。詳細が抽象に依存すべきである。

### 原則の理解

具体的な実装ではなく、インターフェース（抽象）に依存することで、柔軟性とテスト容易性を高めます。

### 悪い例：具体実装への依存

```go
package bad

import (
    "database/sql"
    "fmt"
)

// UserService は具体的な PostgreSQL 実装に直接依存
type UserService struct {
    db *sql.DB // 具体的な実装への依存
}

func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}

func (s *UserService) GetUser(id string) (*User, error) {
    // PostgreSQL 固有のクエリ
    query := "SELECT id, name, email FROM users WHERE id = $1"

    var user User
    err := s.db.QueryRow(query, id).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (s *UserService) CreateUser(name, email string) error {
    query := "INSERT INTO users (name, email) VALUES ($1, $2)"
    _, err := s.db.Exec(query, name, email)
    return err
}
```

**問題点**
- `UserService`はPostgreSQLに直接依存している
- データベースを変更（MySQL、MongoDB など）するには`UserService`を変更する必要がある
- モックを使ったテストが困難
- 依存の方向が逆（上位モジュールが下位モジュールに依存）

### 良い例：抽象への依存

```go
package good

import (
    "context"
)

// User はドメインエンティティ
type User struct {
    ID    string
    Name  string
    Email string
}

// UserRepository は抽象（インターフェース）
type UserRepository interface {
    FindByID(ctx context.Context, id string) (*User, error)
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id string) error
}

// UserService は抽象に依存
type UserService struct {
    repo UserRepository // インターフェースへの依存
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // リポジトリの実装詳細を知らない
    return s.repo.FindByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, name, email string) error {
    user := &User{
        Name:  name,
        Email: email,
    }
    return s.repo.Create(ctx, user)
}

// PostgreSQL実装（詳細が抽象に依存）
type PostgreSQLUserRepository struct {
    db *sql.DB
}

func NewPostgreSQLUserRepository(db *sql.DB) *PostgreSQLUserRepository {
    return &PostgreSQLUserRepository{db: db}
}

func (r *PostgreSQLUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE id = $1"

    var user User
    err := r.db.QueryRowContext(ctx, query, id).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, err
    }

    return &user, nil
}

func (r *PostgreSQLUserRepository) Create(ctx context.Context, user *User) error {
    query := "INSERT INTO users (id, name, email) VALUES ($1, $2, $3)"
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email)
    return err
}

func (r *PostgreSQLUserRepository) Update(ctx context.Context, user *User) error {
    query := "UPDATE users SET name = $2, email = $3 WHERE id = $1"
    _, err := r.db.ExecContext(ctx, query, user.ID, user.Name, user.Email)
    return err
}

func (r *PostgreSQLUserRepository) Delete(ctx context.Context, id string) error {
    query := "DELETE FROM users WHERE id = $1"
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}

// MongoDB実装（別の実装も簡単に追加可能）
type MongoDBUserRepository struct {
    client *mongo.Client
}

func NewMongoDBUserRepository(client *mongo.Client) *MongoDBUserRepository {
    return &MongoDBUserRepository{client: client}
}

func (r *MongoDBUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    // MongoDB固有の実装
    return nil, nil
}

func (r *MongoDBUserRepository) Create(ctx context.Context, user *User) error {
    // MongoDB固有の実装
    return nil
}

func (r *MongoDBUserRepository) Update(ctx context.Context, user *User) error {
    // MongoDB固有の実装
    return nil
}

func (r *MongoDBUserRepository) Delete(ctx context.Context, id string) error {
    // MongoDB固有の実装
    return nil
}

// インメモリ実装（テスト用）
type InMemoryUserRepository struct {
    users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
    return &InMemoryUserRepository{
        users: make(map[string]*User),
    }
}

func (r *InMemoryUserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    user, ok := r.users[id]
    if !ok {
        return nil, fmt.Errorf("user not found: %s", id)
    }
    return user, nil
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *User) error {
    r.users[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) Update(ctx context.Context, user *User) error {
    r.users[user.ID] = user
    return nil
}

func (r *InMemoryUserRepository) Delete(ctx context.Context, id string) error {
    delete(r.users, id)
    return nil
}
```

### テスト例

```go
// テストでは簡単にモックを使用できる
func TestUserService_GetUser(t *testing.T) {
    // Arrange
    mockRepo := NewInMemoryUserRepository()
    service := NewUserService(mockRepo)

    expectedUser := &User{ID: "1", Name: "John", Email: "john@example.com"}
    mockRepo.Create(context.Background(), expectedUser)

    // Act
    user, err := service.GetUser(context.Background(), "1")

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if user.Name != "John" {
        t.Errorf("Expected name John, got %s", user.Name)
    }
}
```

### 使用例

```go
func main() {
    // PostgreSQL を使用
    db, _ := sql.Open("postgres", "connection-string")
    postgresRepo := NewPostgreSQLUserRepository(db)
    service := NewUserService(postgresRepo)

    // または MongoDB を使用（UserService のコードは変更不要）
    mongoClient, _ := mongo.Connect(context.Background(), options.Client())
    mongoRepo := NewMongoDBUserRepository(mongoClient)
    service = NewUserService(mongoRepo)

    // サービスの使用方法は同じ
    service.GetUser(context.Background(), "user123")
}
```

### メリット

1. **柔軟性**：実装を簡単に切り替えられる
2. **テスト容易性**：モックやスタブを使った単体テストが容易
3. **保守性**：上位モジュールを変更せずに下位モジュールを変更できる
4. **再利用性**：抽象に依存するため、異なる文脈で再利用しやすい

---

## SOLID原則のまとめ

### すべての原則を組み合わせた例

```go
package example

import "context"

// S - Single Responsibility: 各型は1つの責任のみ
type User struct {
    ID    string
    Email string
    Name  string
}

// O - Open/Closed: インターフェースを使って拡張可能に
type Notifier interface {
    Notify(ctx context.Context, user *User, message string) error
}

type EmailNotifier struct{}
type SMSNotifier struct{}
type PushNotifier struct{}

// L - Liskov Substitution: すべての実装が同じ契約を守る
func (e *EmailNotifier) Notify(ctx context.Context, user *User, message string) error {
    // Email送信
    return nil
}

func (s *SMSNotifier) Notify(ctx context.Context, user *User, message string) error {
    // SMS送信
    return nil
}

// I - Interface Segregation: 小さく焦点を絞ったインターフェース
type UserReader interface {
    FindByID(ctx context.Context, id string) (*User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
}

// D - Dependency Inversion: 抽象に依存
type UserService struct {
    reader   UserReader
    writer   UserWriter
    notifier Notifier
}

func NewUserService(reader UserReader, writer UserWriter, notifier Notifier) *UserService {
    return &UserService{
        reader:   reader,
        writer:   writer,
        notifier: notifier,
    }
}

func (s *UserService) RegisterUser(ctx context.Context, email, name string) error {
    user := &User{Email: email, Name: name}

    if err := s.writer.Create(ctx, user); err != nil {
        return err
    }

    return s.notifier.Notify(ctx, user, "Welcome!")
}
```

SOLID原則を適用することで、保守性が高く、拡張可能で、テストしやすいコードを書くことができます。
