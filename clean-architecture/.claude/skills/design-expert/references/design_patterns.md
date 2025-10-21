# 設計パターン詳細リファレンス

このドキュメントでは、主要な設計パターンについて詳細な説明と実践的なGo言語での実装例を提供します。

## 目次

1. [生成パターン](#生成パターン)
2. [構造パターン](#構造パターン)
3. [振る舞いパターン](#振る舞いパターン)

---

## 生成パターン

生成パターンは、オブジェクトの生成メカニズムを扱い、状況に応じて適切な方法でオブジェクトを作成します。

### Singleton Pattern

**目的**：クラスのインスタンスが1つだけ存在することを保証する。

**使用場面**
- データベース接続
- 設定管理
- ログシステム
- キャッシュマネージャー

**Go言語での実装**

```go
package singleton

import (
    "database/sql"
    "sync"
)

// DatabaseConnection はシングルトンインスタンス
type DatabaseConnection struct {
    conn *sql.DB
}

var (
    instance *DatabaseConnection
    once     sync.Once
)

// GetInstance はシングルトンインスタンスを返す
// スレッドセーフで、最初の呼び出し時のみインスタンスを生成
func GetInstance() *DatabaseConnection {
    once.Do(func() {
        db, err := sql.Open("postgres", "connection-string")
        if err != nil {
            panic(err)
        }
        instance = &DatabaseConnection{conn: db}
    })
    return instance
}

// Query はデータベースクエリを実行
func (d *DatabaseConnection) Query(query string) (*sql.Rows, error) {
    return d.conn.Query(query)
}
```

**使用例**

```go
func main() {
    db1 := singleton.GetInstance()
    db2 := singleton.GetInstance()

    // db1とdb2は同じインスタンス
    fmt.Printf("Same instance: %v\n", db1 == db2) // true
}
```

**注意点**
- グローバル状態を作るため、テストが困難になる可能性がある
- 依存性注入を使う方が望ましい場合が多い

---

### Factory Pattern

**目的**：オブジェクト生成のロジックをカプセル化し、インターフェースを通じてオブジェクトを作成する。

**使用場面**
- データベースの種類に応じたリポジトリの生成
- ファイル形式に応じたパーサーの生成
- プロトコルに応じたクライアントの生成

**実装**

```go
package factory

import "fmt"

// Storage はストレージのインターフェース
type Storage interface {
    Save(key string, value []byte) error
    Load(key string) ([]byte, error)
    Delete(key string) error
}

// MemoryStorage はメモリベースのストレージ
type MemoryStorage struct {
    data map[string][]byte
}

func (m *MemoryStorage) Save(key string, value []byte) error {
    if m.data == nil {
        m.data = make(map[string][]byte)
    }
    m.data[key] = value
    return nil
}

func (m *MemoryStorage) Load(key string) ([]byte, error) {
    value, ok := m.data[key]
    if !ok {
        return nil, fmt.Errorf("key not found: %s", key)
    }
    return value, nil
}

func (m *MemoryStorage) Delete(key string) error {
    delete(m.data, key)
    return nil
}

// FileStorage はファイルベースのストレージ
type FileStorage struct {
    basePath string
}

func (f *FileStorage) Save(key string, value []byte) error {
    // ファイルシステムへの保存処理
    return nil
}

func (f *FileStorage) Load(key string) ([]byte, error) {
    // ファイルシステムからの読み込み処理
    return nil, nil
}

func (f *FileStorage) Delete(key string) error {
    // ファイル削除処理
    return nil
}

// StorageFactory はストレージのファクトリー
type StorageFactory struct{}

// CreateStorage はストレージタイプに応じたインスタンスを生成
func (f *StorageFactory) CreateStorage(storageType string) (Storage, error) {
    switch storageType {
    case "memory":
        return &MemoryStorage{}, nil
    case "file":
        return &FileStorage{basePath: "/tmp/storage"}, nil
    default:
        return nil, fmt.Errorf("unknown storage type: %s", storageType)
    }
}
```

**使用例**

```go
func main() {
    factory := &factory.StorageFactory{}

    // メモリストレージを使用
    storage, err := factory.CreateStorage("memory")
    if err != nil {
        log.Fatal(err)
    }

    storage.Save("key1", []byte("value1"))
    value, _ := storage.Load("key1")
    fmt.Println(string(value))

    // ファイルストレージに切り替え
    fileStorage, _ := factory.CreateStorage("file")
    fileStorage.Save("key2", []byte("value2"))
}
```

---

### Abstract Factory Pattern

**目的**：関連するオブジェクトのファミリーを、具体的なクラスを指定せずに生成する。

**実装**

```go
package abstractfactory

// UIFactory は UI コンポーネントファミリーのファクトリー
type UIFactory interface {
    CreateButton() Button
    CreateTextBox() TextBox
}

// Button インターフェース
type Button interface {
    Render() string
}

// TextBox インターフェース
type TextBox interface {
    Render() string
}

// Windows用のファクトリー
type WindowsFactory struct{}

func (w *WindowsFactory) CreateButton() Button {
    return &WindowsButton{}
}

func (w *WindowsFactory) CreateTextBox() TextBox {
    return &WindowsTextBox{}
}

type WindowsButton struct{}

func (b *WindowsButton) Render() string {
    return "[Windows Button]"
}

type WindowsTextBox struct{}

func (t *WindowsTextBox) Render() string {
    return "[Windows TextBox]"
}

// Mac用のファクトリー
type MacFactory struct{}

func (m *MacFactory) CreateButton() Button {
    return &MacButton{}
}

func (m *MacFactory) CreateTextBox() TextBox {
    return &MacTextBox{}
}

type MacButton struct{}

func (b *MacButton) Render() string {
    return "[Mac Button]"
}

type MacTextBox struct{}

func (t *MacTextBox) Render() string {
    return "[Mac TextBox]"
}

// アプリケーション
func RenderUI(factory UIFactory) {
    button := factory.CreateButton()
    textBox := factory.CreateTextBox()

    fmt.Println(button.Render())
    fmt.Println(textBox.Render())
}
```

---

### Builder Pattern

**目的**：複雑なオブジェクトの構築過程を段階的に行い、同じ構築過程で異なる表現を作成できるようにする。

**実装**

```go
package builder

import "fmt"

// HTTPRequest は構築したいオブジェクト
type HTTPRequest struct {
    Method  string
    URL     string
    Headers map[string]string
    Body    []byte
    Timeout int
}

// HTTPRequestBuilder はビルダー
type HTTPRequestBuilder struct {
    request *HTTPRequest
}

// NewHTTPRequestBuilder はビルダーの新しいインスタンスを作成
func NewHTTPRequestBuilder() *HTTPRequestBuilder {
    return &HTTPRequestBuilder{
        request: &HTTPRequest{
            Headers: make(map[string]string),
            Method:  "GET", // デフォルト値
            Timeout: 30,    // デフォルト値
        },
    }
}

// WithMethod はHTTPメソッドを設定
func (b *HTTPRequestBuilder) WithMethod(method string) *HTTPRequestBuilder {
    b.request.Method = method
    return b
}

// WithURL はURLを設定
func (b *HTTPRequestBuilder) WithURL(url string) *HTTPRequestBuilder {
    b.request.URL = url
    return b
}

// WithHeader はヘッダーを追加
func (b *HTTPRequestBuilder) WithHeader(key, value string) *HTTPRequestBuilder {
    b.request.Headers[key] = value
    return b
}

// WithBody はリクエストボディを設定
func (b *HTTPRequestBuilder) WithBody(body []byte) *HTTPRequestBuilder {
    b.request.Body = body
    return b
}

// WithTimeout はタイムアウトを設定
func (b *HTTPRequestBuilder) WithTimeout(timeout int) *HTTPRequestBuilder {
    b.request.Timeout = timeout
    return b
}

// Build は構築したHTTPRequestを返す
func (b *HTTPRequestBuilder) Build() (*HTTPRequest, error) {
    if b.request.URL == "" {
        return nil, fmt.Errorf("URL is required")
    }
    return b.request, nil
}
```

**使用例**

```go
func main() {
    request, err := builder.NewHTTPRequestBuilder().
        WithMethod("POST").
        WithURL("https://api.example.com/users").
        WithHeader("Content-Type", "application/json").
        WithHeader("Authorization", "Bearer token123").
        WithBody([]byte(`{"name":"John"}`)).
        WithTimeout(60).
        Build()

    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", request)
}
```

---

### Prototype Pattern

**目的**：既存のオブジェクトをクローンして新しいオブジェクトを作成する。

**実装**

```go
package prototype

import "time"

// Cloneable はクローン可能なオブジェクトのインターフェース
type Cloneable interface {
    Clone() Cloneable
}

// Document はドキュメント
type Document struct {
    Title     string
    Content   string
    Author    string
    CreatedAt time.Time
    Tags      []string
}

// Clone はドキュメントのクローンを作成
func (d *Document) Clone() *Document {
    // ディープコピー
    tags := make([]string, len(d.Tags))
    copy(tags, d.Tags)

    return &Document{
        Title:     d.Title,
        Content:   d.Content,
        Author:    d.Author,
        CreatedAt: time.Now(), // 新しい作成日時
        Tags:      tags,
    }
}

// DocumentRegistry はプロトタイプのレジストリ
type DocumentRegistry struct {
    prototypes map[string]*Document
}

func NewDocumentRegistry() *DocumentRegistry {
    return &DocumentRegistry{
        prototypes: make(map[string]*Document),
    }
}

// Register はプロトタイプを登録
func (r *DocumentRegistry) Register(key string, prototype *Document) {
    r.prototypes[key] = prototype
}

// CreateDocument は登録されたプロトタイプからドキュメントを作成
func (r *DocumentRegistry) CreateDocument(key string) *Document {
    if prototype, ok := r.prototypes[key]; ok {
        return prototype.Clone()
    }
    return nil
}
```

**使用例**

```go
func main() {
    registry := prototype.NewDocumentRegistry()

    // プロトタイプを登録
    blogPost := &prototype.Document{
        Title:   "Default Blog Post",
        Author:  "System",
        Tags:    []string{"blog", "default"},
    }
    registry.Register("blog-post", blogPost)

    // プロトタイプから新しいドキュメントを作成
    newPost := registry.CreateDocument("blog-post")
    newPost.Title = "My First Post"
    newPost.Content = "Hello, World!"

    fmt.Printf("%+v\n", newPost)
}
```

---

## 構造パターン

構造パターンは、クラスやオブジェクトを組み合わせて、より大きな構造を作る方法を扱います。

### Adapter Pattern

**目的**：互換性のないインターフェースを持つクラスを一緒に動作させる。

**実装**

```go
package adapter

import "fmt"

// ModernLogger は新しいロギングインターフェース
type ModernLogger interface {
    LogInfo(message string)
    LogError(message string)
    LogDebug(message string)
}

// LegacyLogger は古いロギングシステム
type LegacyLogger struct{}

func (l *LegacyLogger) Log(level string, message string) {
    fmt.Printf("[%s] %s\n", level, message)
}

// LegacyLoggerAdapter は古いロガーを新しいインターフェースに適合させる
type LegacyLoggerAdapter struct {
    legacyLogger *LegacyLogger
}

func NewLegacyLoggerAdapter(legacy *LegacyLogger) *LegacyLoggerAdapter {
    return &LegacyLoggerAdapter{legacyLogger: legacy}
}

func (a *LegacyLoggerAdapter) LogInfo(message string) {
    a.legacyLogger.Log("INFO", message)
}

func (a *LegacyLoggerAdapter) LogError(message string) {
    a.legacyLogger.Log("ERROR", message)
}

func (a *LegacyLoggerAdapter) LogDebug(message string) {
    a.legacyLogger.Log("DEBUG", message)
}

// Application は新しいインターフェースを使用
func UseLogger(logger ModernLogger) {
    logger.LogInfo("Application started")
    logger.LogError("An error occurred")
    logger.LogDebug("Debug information")
}
```

---

### Decorator Pattern

**目的**：オブジェクトに動的に新しい責任を追加する。

**実装**

```go
package decorator

import (
    "fmt"
    "time"
)

// DataSource はデータソースのインターフェース
type DataSource interface {
    Write(data []byte) error
    Read() ([]byte, error)
}

// FileDataSource はファイルベースのデータソース
type FileDataSource struct {
    filename string
}

func (f *FileDataSource) Write(data []byte) error {
    fmt.Printf("Writing to file: %s\n", f.filename)
    return nil
}

func (f *FileDataSource) Read() ([]byte, error) {
    fmt.Printf("Reading from file: %s\n", f.filename)
    return []byte("file content"), nil
}

// EncryptionDecorator は暗号化機能を追加
type EncryptionDecorator struct {
    wrapped DataSource
}

func (e *EncryptionDecorator) Write(data []byte) error {
    encrypted := e.encrypt(data)
    return e.wrapped.Write(encrypted)
}

func (e *EncryptionDecorator) Read() ([]byte, error) {
    data, err := e.wrapped.Read()
    if err != nil {
        return nil, err
    }
    return e.decrypt(data), nil
}

func (e *EncryptionDecorator) encrypt(data []byte) []byte {
    fmt.Println("Encrypting data")
    return data
}

func (e *EncryptionDecorator) decrypt(data []byte) []byte {
    fmt.Println("Decrypting data")
    return data
}

// CompressionDecorator は圧縮機能を追加
type CompressionDecorator struct {
    wrapped DataSource
}

func (c *CompressionDecorator) Write(data []byte) error {
    compressed := c.compress(data)
    return c.wrapped.Write(compressed)
}

func (c *CompressionDecorator) Read() ([]byte, error) {
    data, err := c.wrapped.Read()
    if err != nil {
        return nil, err
    }
    return c.decompress(data), nil
}

func (c *CompressionDecorator) compress(data []byte) []byte {
    fmt.Println("Compressing data")
    return data
}

func (c *CompressionDecorator) decompress(data []byte) []byte {
    fmt.Println("Decompressing data")
    return data
}

// LoggingDecorator はロギング機能を追加
type LoggingDecorator struct {
    wrapped DataSource
}

func (l *LoggingDecorator) Write(data []byte) error {
    fmt.Printf("[%s] Writing %d bytes\n", time.Now().Format(time.RFC3339), len(data))
    return l.wrapped.Write(data)
}

func (l *LoggingDecorator) Read() ([]byte, error) {
    fmt.Printf("[%s] Reading data\n", time.Now().Format(time.RFC3339))
    return l.wrapped.Read()
}
```

**使用例**

```go
func main() {
    // 基本のファイルデータソース
    source := &decorator.FileDataSource{filename: "data.txt"}

    // デコレーターを重ねる
    source = &decorator.EncryptionDecorator{wrapped: source}
    source = &decorator.CompressionDecorator{wrapped: source}
    source = &decorator.LoggingDecorator{wrapped: source}

    // 使用時は透過的
    source.Write([]byte("Hello, World!"))
    // 出力：
    // [2024-01-01T12:00:00Z] Writing 13 bytes
    // Compressing data
    // Encrypting data
    // Writing to file: data.txt

    source.Read()
}
```

---

### Facade Pattern

**目的**：複雑なサブシステムに対して、シンプルなインターフェースを提供する。

**実装**

```go
package facade

import "fmt"

// 複雑なサブシステム
type VideoFile struct {
    filename string
}

type Codec interface {
    Encode(data []byte) []byte
    Decode(data []byte) []byte
}

type MPEG4Codec struct{}

func (m *MPEG4Codec) Encode(data []byte) []byte {
    fmt.Println("Encoding with MPEG4")
    return data
}

func (m *MPEG4Codec) Decode(data []byte) []byte {
    fmt.Println("Decoding with MPEG4")
    return data
}

type AudioMixer struct{}

func (a *AudioMixer) Mix(audio []byte) []byte {
    fmt.Println("Mixing audio")
    return audio
}

type BitrateReader struct{}

func (b *BitrateReader) Read(filename string) []byte {
    fmt.Printf("Reading file: %s\n", filename)
    return []byte("video data")
}

// Facade：シンプルなインターフェース
type VideoConverter struct {
    codec  Codec
    mixer  *AudioMixer
    reader *BitrateReader
}

func NewVideoConverter() *VideoConverter {
    return &VideoConverter{
        codec:  &MPEG4Codec{},
        mixer:  &AudioMixer{},
        reader: &BitrateReader{},
    }
}

// Convert は複雑な変換処理を隠蔽
func (v *VideoConverter) Convert(filename string, format string) []byte {
    fmt.Println("Starting video conversion...")

    // 複雑な処理を内部で実行
    data := v.reader.Read(filename)
    data = v.codec.Decode(data)
    data = v.mixer.Mix(data)
    data = v.codec.Encode(data)

    fmt.Println("Conversion complete!")
    return data
}
```

---

### Proxy Pattern

**目的**：他のオブジェクトへのアクセスを制御するための代理オブジェクトを提供する。

**実装**

```go
package proxy

import (
    "fmt"
    "time"
)

// Image はイメージのインターフェース
type Image interface {
    Display() error
}

// RealImage は実際のイメージ（重い処理）
type RealImage struct {
    filename string
    data     []byte
}

func NewRealImage(filename string) *RealImage {
    img := &RealImage{filename: filename}
    img.loadFromDisk()
    return img
}

func (r *RealImage) loadFromDisk() {
    fmt.Printf("Loading image from disk: %s\n", r.filename)
    time.Sleep(2 * time.Second) // 重い処理をシミュレート
    r.data = []byte("image data")
}

func (r *RealImage) Display() error {
    fmt.Printf("Displaying image: %s\n", r.filename)
    return nil
}

// ImageProxy はイメージのプロキシ（遅延ロード）
type ImageProxy struct {
    filename string
    realImage *RealImage
}

func NewImageProxy(filename string) *ImageProxy {
    return &ImageProxy{filename: filename}
}

func (p *ImageProxy) Display() error {
    // 実際に必要になるまでロードしない（遅延ロード）
    if p.realImage == nil {
        p.realImage = NewRealImage(p.filename)
    }
    return p.realImage.Display()
}

// CachingProxy はキャッシング機能を持つプロキシ
type CachingProxy struct {
    target Service
    cache  map[string]interface{}
}

type Service interface {
    GetData(key string) (interface{}, error)
}

func NewCachingProxy(target Service) *CachingProxy {
    return &CachingProxy{
        target: target,
        cache:  make(map[string]interface{}),
    }
}

func (c *CachingProxy) GetData(key string) (interface{}, error) {
    // キャッシュをチェック
    if data, ok := c.cache[key]; ok {
        fmt.Println("Cache hit!")
        return data, nil
    }

    // キャッシュミス時は実際のサービスから取得
    fmt.Println("Cache miss, fetching from service")
    data, err := c.target.GetData(key)
    if err != nil {
        return nil, err
    }

    // キャッシュに保存
    c.cache[key] = data
    return data, nil
}
```

---

## 振る舞いパターン

振る舞いパターンは、オブジェクト間の責任の分配とコミュニケーションを扱います。

### Strategy Pattern

**目的**：アルゴリズムのファミリーを定義し、それぞれをカプセル化して交換可能にする。

**実装**

```go
package strategy

import "fmt"

// PaymentStrategy は支払い戦略のインターフェース
type PaymentStrategy interface {
    Pay(amount float64) error
    GetName() string
}

// CreditCardStrategy はクレジットカード支払い
type CreditCardStrategy struct {
    cardNumber string
    cvv        string
}

func (c *CreditCardStrategy) Pay(amount float64) error {
    fmt.Printf("Paying %.2f using Credit Card (ending in %s)\n",
        amount, c.cardNumber[len(c.cardNumber)-4:])
    return nil
}

func (c *CreditCardStrategy) GetName() string {
    return "Credit Card"
}

// PayPalStrategy はPayPal支払い
type PayPalStrategy struct {
    email string
}

func (p *PayPalStrategy) Pay(amount float64) error {
    fmt.Printf("Paying %.2f using PayPal (%s)\n", amount, p.email)
    return nil
}

func (p *PayPalStrategy) GetName() string {
    return "PayPal"
}

// CryptocurrencyStrategy は暗号通貨支払い
type CryptocurrencyStrategy struct {
    walletAddress string
}

func (c *CryptocurrencyStrategy) Pay(amount float64) error {
    fmt.Printf("Paying %.2f using Cryptocurrency (wallet: %s)\n",
        amount, c.walletAddress)
    return nil
}

func (c *CryptocurrencyStrategy) GetName() string {
    return "Cryptocurrency"
}

// PaymentProcessor は支払い処理システム
type PaymentProcessor struct {
    strategy PaymentStrategy
}

// SetStrategy は支払い戦略を設定
func (p *PaymentProcessor) SetStrategy(strategy PaymentStrategy) {
    p.strategy = strategy
}

// ProcessPayment は支払いを処理
func (p *PaymentProcessor) ProcessPayment(amount float64) error {
    if p.strategy == nil {
        return fmt.Errorf("payment strategy not set")
    }

    fmt.Printf("Processing payment with %s strategy\n", p.strategy.GetName())
    return p.strategy.Pay(amount)
}
```

**使用例**

```go
func main() {
    processor := &strategy.PaymentProcessor{}

    // クレジットカードで支払い
    processor.SetStrategy(&strategy.CreditCardStrategy{
        cardNumber: "1234-5678-9012-3456",
        cvv:        "123",
    })
    processor.ProcessPayment(100.00)

    // PayPalに切り替え
    processor.SetStrategy(&strategy.PayPalStrategy{
        email: "user@example.com",
    })
    processor.ProcessPayment(50.00)

    // 暗号通貨に切り替え
    processor.SetStrategy(&strategy.CryptocurrencyStrategy{
        walletAddress: "1A1zP1eP5QGefi2DMPTfTL5SLmv7DivfNa",
    })
    processor.ProcessPayment(0.001)
}
```

---

### Observer Pattern

**目的**：オブジェクト間の1対多の依存関係を定義し、1つのオブジェクトの状態が変化したときに、依存するすべてのオブジェクトに自動的に通知する。

**実装**

```go
package observer

import (
    "fmt"
    "sync"
)

// Event はイベントデータ
type Event struct {
    Type string
    Data interface{}
}

// Observer はオブザーバーのインターフェース
type Observer interface {
    Update(event Event)
    GetID() string
}

// Subject はサブジェクト（監視対象）
type Subject struct {
    observers map[string]Observer
    mu        sync.RWMutex
}

func NewSubject() *Subject {
    return &Subject{
        observers: make(map[string]Observer),
    }
}

// Attach はオブザーバーを登録
func (s *Subject) Attach(observer Observer) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.observers[observer.GetID()] = observer
}

// Detach はオブザーバーを削除
func (s *Subject) Detach(observerID string) {
    s.mu.Lock()
    defer s.mu.Unlock()
    delete(s.observers, observerID)
}

// Notify はすべてのオブザーバーに通知
func (s *Subject) Notify(event Event) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, observer := range s.observers {
        go observer.Update(event) // 非同期で通知
    }
}

// EmailNotifier はメール通知オブザーバー
type EmailNotifier struct {
    id    string
    email string
}

func NewEmailNotifier(id, email string) *EmailNotifier {
    return &EmailNotifier{id: id, email: email}
}

func (e *EmailNotifier) Update(event Event) {
    fmt.Printf("[Email %s] Sending notification to %s: %v\n",
        e.id, e.email, event.Data)
}

func (e *EmailNotifier) GetID() string {
    return e.id
}

// SMSNotifier はSMS通知オブザーバー
type SMSNotifier struct {
    id          string
    phoneNumber string
}

func NewSMSNotifier(id, phoneNumber string) *SMSNotifier {
    return &SMSNotifier{id: id, phoneNumber: phoneNumber}
}

func (s *SMSNotifier) Update(event Event) {
    fmt.Printf("[SMS %s] Sending notification to %s: %v\n",
        s.id, s.phoneNumber, event.Data)
}

func (s *SMSNotifier) GetID() string {
    return s.id
}

// LogObserver はログ記録オブザーバー
type LogObserver struct {
    id string
}

func NewLogObserver(id string) *LogObserver {
    return &LogObserver{id: id}
}

func (l *LogObserver) Update(event Event) {
    fmt.Printf("[Log %s] Recording event: %s - %v\n",
        l.id, event.Type, event.Data)
}

func (l *LogObserver) GetID() string {
    return l.id
}
```

**使用例**

```go
func main() {
    // サブジェクトを作成
    subject := observer.NewSubject()

    // オブザーバーを登録
    emailNotifier := observer.NewEmailNotifier("email1", "user@example.com")
    smsNotifier := observer.NewSMSNotifier("sms1", "+1234567890")
    logObserver := observer.NewLogObserver("log1")

    subject.Attach(emailNotifier)
    subject.Attach(smsNotifier)
    subject.Attach(logObserver)

    // イベントを通知
    subject.Notify(observer.Event{
        Type: "UserRegistered",
        Data: map[string]string{"username": "john_doe"},
    })

    // オブザーバーを削除
    subject.Detach("sms1")

    // 再度通知（SMSは通知されない）
    subject.Notify(observer.Event{
        Type: "OrderPlaced",
        Data: map[string]interface{}{"order_id": 12345, "amount": 99.99},
    })
}
```

---

### Command Pattern

**目的**：リクエストをオブジェクトとしてカプセル化し、リクエストのパラメータ化、キュー化、ログ記録、アンドゥ操作を可能にする。

**実装**

```go
package command

import "fmt"

// Command はコマンドのインターフェース
type Command interface {
    Execute() error
    Undo() error
}

// Receiver：実際の処理を行うオブジェクト
type TextEditor struct {
    text string
}

func (t *TextEditor) Insert(text string) {
    t.text += text
}

func (t *TextEditor) Delete(length int) {
    if length > len(t.text) {
        length = len(t.text)
    }
    t.text = t.text[:len(t.text)-length]
}

func (t *TextEditor) GetText() string {
    return t.text
}

// InsertCommand は挿入コマンド
type InsertCommand struct {
    editor *TextEditor
    text   string
}

func (i *InsertCommand) Execute() error {
    i.editor.Insert(i.text)
    fmt.Printf("Inserted: '%s'\n", i.text)
    return nil
}

func (i *InsertCommand) Undo() error {
    i.editor.Delete(len(i.text))
    fmt.Printf("Undid insert: '%s'\n", i.text)
    return nil
}

// DeleteCommand は削除コマンド
type DeleteCommand struct {
    editor      *TextEditor
    length      int
    deletedText string
}

func (d *DeleteCommand) Execute() error {
    textLen := len(d.editor.text)
    if d.length > textLen {
        d.length = textLen
    }
    d.deletedText = d.editor.text[textLen-d.length:]
    d.editor.Delete(d.length)
    fmt.Printf("Deleted %d characters\n", d.length)
    return nil
}

func (d *DeleteCommand) Undo() error {
    d.editor.Insert(d.deletedText)
    fmt.Printf("Undid delete: '%s'\n", d.deletedText)
    return nil
}

// CommandHistory はコマンド履歴を管理
type CommandHistory struct {
    commands []Command
}

func (h *CommandHistory) Execute(cmd Command) error {
    if err := cmd.Execute(); err != nil {
        return err
    }
    h.commands = append(h.commands, cmd)
    return nil
}

func (h *CommandHistory) Undo() error {
    if len(h.commands) == 0 {
        return fmt.Errorf("nothing to undo")
    }

    lastCmd := h.commands[len(h.commands)-1]
    h.commands = h.commands[:len(h.commands)-1]

    return lastCmd.Undo()
}
```

---

このリファレンスは、主要な設計パターンをカバーしていますが、すべてのパターンを網羅しているわけではありません。実際のプロジェクトでは、これらのパターンを組み合わせたり、プロジェクトの特性に合わせて調整したりすることが重要です。
