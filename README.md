# Commons

一個輕量的 Go 工具庫，提供常用功能的泛型與高效能實作。

## 功能特色

- 🚀 **高效能**：針對常見場景優化，降低記憶體分配
- 🔧 **泛型支援**：充分利用 Go 1.18+ 泛型特性
- 📦 **模組化**：各子套件獨立，按需引用
- ✅ **測試完整**：包含單元測試與效能基準測試

## 安裝

```bash
go get github.com/vincent119/commons
```

## 套件列表

### stringx - 字串處理

提供高效的字串轉換與處理功能。

```go
import "github.com/vincent119/commons/stringx"

// 轉換為 snake_case
s := stringx.ToSnake("UserID") // "user_i_d"
```

**主要函式：**

- `ToSnake(s string) string` - 將字串轉為 snake_case

### errorx - 錯誤處理

簡化錯誤包裝與判斷的工具函式。

```go
import "github.com/vincent119/commons/errorx"

// 包裝錯誤
err := errorx.Wrap(someErr, "操作失敗")

// 判斷錯誤類型
if errorx.Is(err, target) { ... }

// 取得底層錯誤
cause := errorx.Cause(err)
```

**主要函式：**

- `Wrap(err error, msg string) error` - 包裝錯誤並加上訊息
- `Is(err, target error) bool` - 判斷錯誤鏈是否包含特定錯誤
- `As[T any](err error, target *T) bool` - 嘗試將錯誤轉型
- `Cause(err error) error` - 取出最底層錯誤

### slicex - 切片操作

提供泛型切片操作函式，類似 JavaScript 的陣列方法。

```go
import "github.com/vincent119/commons/slicex"

// 檢查元素是否存在
found := slicex.Contains([]int{1, 2, 3}, 2) // true

// 尋找元素索引
idx := slicex.IndexOf([]string{"a", "b", "c"}, "b") // 1

// 過濾元素
evens := slicex.Filter([]int{1, 2, 3, 4}, func(n int) bool {
    return n%2 == 0
}) // [2, 4]

// 映射轉換
strs := slicex.Map([]int{1, 2, 3}, func(n int) string {
    return fmt.Sprintf("%d", n)
}) // ["1", "2", "3"]
```

**主要函式：**

- `Contains[T comparable](s []T, v T) bool` - 檢查是否包含元素
- `IndexOf[T comparable](s []T, v T) int` - 回傳元素索引
- `Filter[T any](s []T, f func(T) bool) []T` - 過濾元素
- `Map[T any, R any](s []T, f func(T) R) []R` - 映射轉換

### timex - 時間處理

提供時區安全的時間操作函式。

```go
import "github.com/vincent119/commons/timex"

// 取得 UTC 時間
utc := timex.NowUTC()

// 取得某天的零點（指定時區）
start := timex.StartOfDay(time.Now(), time.Local)

// 時間截斷
truncated := timex.TruncateTo(time.Now(), time.Hour)

// 時間戳
ts := timex.TimeStampUTC() // "2025-12-18T10:30:00.000Z"
date := timex.DateStamp()  // "2025-12-18"
```

**主要函式：**

- `NowUTC() time.Time` - 取得目前 UTC 時間
- `StartOfDay(t time.Time, loc *time.Location) time.Time` - 回傳指定時區的零點
- `TruncateTo(t time.Time, d time.Duration) time.Time` - 截斷時間至指定粒度
- `FormatTime(t time.Time, layout string) string` - 格式化時間
- `ParseTime(str, layout string) (time.Time, error)` - 解析時間字串
- `TimeStamp() string` - 取得目前時間戳
- `TimeStampUTC() string` - 取得 UTC 時間戳
- `DateStamp() string` - 取得目前日期

### uuidx - UUID 工具

封裝 google/uuid 的便利函式。

```go
import "github.com/vincent119/commons/uuidx"

// 產生新的 UUID
id := uuidx.NewUUID()

// 驗證 UUID 格式
valid := uuidx.IsValidUUID("550e8400-e29b-41d4-a716-446655440000")
```

**主要函式：**

- `NewUUID() string` - 產生新的 UUID v4
- `NewUUIDv4() string` - 產生新的 UUID v4
- `NewUUIDv5(namespace uuid.UUID, name string) string` - 產生 UUID v5
- `IsValidUUID(s string) bool` - 驗證 UUID 格式

### cryptox - 加密工具

提供常用的雜湊函式。

```go
import "github.com/vincent119/commons/cryptox"

// MD5 雜湊
hash := cryptox.MD5Hash("password")

// SHA256 雜湊
sha := cryptox.SHA256Hash("data")
```

**主要函式：**

- `MD5Hash(s string) string` - 回傳 MD5 雜湊
- `SHA256Hash(s string) string` - 回傳 SHA256 雜湊

### validatorx - 驗證工具

提供常用的格式驗證函式。

```go
import "github.com/vincent119/commons/validatorx"

// Email 驗證
valid := validatorx.IsEmail("user@example.com")

// 手機號驗證（台灣）
valid := validatorx.IsMobile("0912345678")

// UUID 驗證
valid := validatorx.IsUUID("550e8400-e29b-41d4-a716-446655440000")

// IPv4 驗證
valid := validatorx.IsIPv4("192.168.1.1")
```

**主要函式：**

- `IsEmail(email string) bool` - 驗證 Email 格式
- `IsMobile(mobile string) bool` - 驗證台灣手機號格式
- `IsUUID(u string) bool` - 驗證 UUID 格式
- `IsIPv4(ip string) bool` - 驗證 IPv4 格式
- `IsIPv6(ip string) bool` - 驗證 IPv6 格式

### modelx - 通用模型

提供常用的回應結構定義。

```go
import "github.com/vincent119/commons/modelx"

// 錯誤回應
errResp := modelx.ErrorResponse{
    Code:    401,
    Message: "unauthorized",
}

// 健康檢查回應
health := modelx.ResponseHealthCheck{
    Status: "ok",
}
```

## 完整範例

```go
package main

import (
    "fmt"
    "github.com/vincent119/commons/stringx"
    "github.com/vincent119/commons/timex"
    "github.com/vincent119/commons/slicex"
    "github.com/vincent119/commons/uuidx"
    "github.com/vincent119/commons/validatorx"
)

func main() {
    // 字串處理
    snake := stringx.ToSnake("UserProfileData")
    fmt.Println(snake) // "user_profile_data"

    // 時間處理
    now := timex.NowUTC()
    timestamp := timex.TimeStampUTC()
    fmt.Println(timestamp)

    // 切片操作
    nums := []int{1, 2, 3, 4, 5}
    evens := slicex.Filter(nums, func(n int) bool { return n%2 == 0 })
    fmt.Println(evens) // [2, 4]

    // UUID 產生
    id := uuidx.NewUUID()
    fmt.Println(id)

    // 驗證
    isEmail := validatorx.IsEmail("test@example.com")
    fmt.Println(isEmail) // true
}
```

## 開發指令

```bash
# 整理依賴
make tidy

# 程式碼檢查
make lint

# 執行測試
make test

# 效能基準測試
make bench
```

## 系統需求

- Go 1.25+

## 依賴

- `github.com/google/uuid` v1.6.0

## 授權

MIT License

## 貢獻

歡迎提交 Pull Request 或回報 Issue！
