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

| 套件 | 說明 |
|------|------|
| `stringx` | 字串處理（snake_case、跳脫等）|
| `errorx` | 錯誤包裝與判斷 |
| `slicex` | 泛型切片操作（Contains、Filter、Map 等）|
| `timex` | 時區安全的時間操作 |
| `uuidx` | UUID 產生與驗證 |
| `cryptox` | MD5、SHA256 雜湊 |
| `validatorx` | 格式驗證（Email、手機、IP 等）|
| `ipx` | IP 位址工具（驗證、轉換、網段、GeoIP）|
| `sqlx` | SQL 查詢工具（LIKE 跳脫、字串跳脫）|
| `jsonx` | JSON 字串跳脫 |
| `pathx` | 路徑處理（分隔符正規化）|
| `httpx/resp` | HTTP 回應結構定義 |
| `structx` | 結構體轉 Map (StructToMap) |
| `graceful` | 優雅關機與生命週期管理 |
| `awsx/s3` | AWS S3 路徑工具 |

---

### ipx - IP 位址工具

提供完整的 IP 位址驗證、轉換、網段計算與 GeoIP 整合。

```go
import "github.com/vincent119/commons/ipx"

// IP 驗證
ipx.IsValidIP("192.168.1.1")     // true
ipx.IsIPv4("192.168.1.1")        // true
ipx.IsIPv6("2001:db8::1")        // true
ipx.IsPublicIP("8.8.8.8")        // true

// IP 轉換
n, _ := ipx.IPv4ToUint32("192.168.1.1")  // 3232235777
ip := ipx.Uint32ToIPv4(3232235777)       // "192.168.1.1"
expanded, _ := ipx.ExpandIPv6("::1")     // "0000:0000:...:0001"

// 網段工具
inCIDR, _ := ipx.IsIPInCIDR("192.168.1.100", "192.168.1.0/24") // true
info, _ := ipx.GetNetworkInfo("192.168.1.0/24")
// info.Network = "192.168.1.0"
// info.Broadcast = "192.168.1.255"
// info.TotalHosts = 254

// 地理位置
location := ipx.GetLocationByIP("127.0.0.1")     // "本地"
location := ipx.GetLocationByIP("192.168.1.1")   // "內部網絡"

// 取得客戶端 IP（從 HTTP headers）
clientIP := ipx.GetClientIP(headers)

// 取得本機 IP
localIPs := ipx.GetLocalIPs()  // "192.168.1.100,10.0.0.5"
```

---

### sqlx - SQL 查詢工具

提供 LIKE 查詢跳脫與 SQL 字串處理。

```go
import "github.com/vincent119/commons/sqlx"

// LIKE 查詢跳脫
escaped := sqlx.EscapeLikeQuery("50%_off")  // "50\%\_off"

// 建構 LIKE 查詢值（含通配符）
like := sqlx.BuildLikeQueryValue("test", sqlx.LikePosBoth)  // "%test%"

// 位置常數
// - LikePosStart: 前綴匹配 "value%"
// - LikePosEnd:   後綴匹配 "%value"
// - LikePosBoth:  包含匹配 "%value%"

// 搭配 ESCAPE 子句
query := "WHERE name LIKE ? " + sqlx.LikeEscapeClause()

// SQL 字串跳脫（不能取代 prepared statement）
escaped := sqlx.EscapeSQLString("O'Reilly")  // "O\'Reilly"

// Log 格式化
formatted := sqlx.FormatSQLForLog("SELECT * FROM   users")
```

**主要函式：**

- `EscapeLikeQuery(s string) string` - 轉義 LIKE 特殊字元
- `BuildLikeQueryValue(value, position string) string` - 建構 LIKE 查詢值
- `LikeEscapeClause() string` - 回傳 ESCAPE 子句
- `EscapeSQLString(s string) string` - 基礎 SQL 跳脫
- `FormatSQLForLog(sql string) string` - Log 格式化

---

### jsonx - JSON 處理工具

提供 JSON 字串跳脫。

```go
import "github.com/vincent119/commons/jsonx"

// JSON 字串跳脫
escaped := jsonx.EscapeJSON("line1\nline2")  // "line1\\nline2"
```

**主要函式：**

- `EscapeJSON(s string) string` - 跳脫 JSON 特殊字元（\, ", \n, \r, \t）

---

### pathx - 路徑處理工具

提供跨平台路徑處理。

```go
import "github.com/vincent119/commons/pathx"

// 路徑分隔符正規化
path := pathx.NormalizePathSeparator("a\\b\\c")  // "a/b/c"
```

**主要函式：**

- `NormalizePathSeparator(path string) string` - 將 \ 轉換為 /

---

### awsx/s3 - AWS S3 路徑工具

提供 S3 路徑前綴建構工具。

```go
import "github.com/vincent119/commons/awsx/s3"

// 建構 S3 路徑前綴
prefix := s3.BuildS3Prefix("bucket/prefix", "media/images")
// "bucket/prefix/media/images/"

// 通用路徑前綴（支援多段）
prefix := s3.BuildPrefix("uploads", "2025", "12")
// "uploads/2025/12/"
```

---

### stringx - 字串處理

提供高效的字串轉換與處理功能。

```go
import "github.com/vincent119/commons/stringx"

// 轉換為 snake_case
s := stringx.ToSnake("UserID")  // "user_i_d"

// 反斜線處理
escaped := stringx.EscapeBackslash("a\\b")    // "a\\\\b"
unescaped := stringx.UnescapeBackslash("a\\\\b")  // "a\\b"

// 字串工具
stringx.IsEmpty("")       // true
stringx.Truncate("hello world", 5)  // "hello"
```

**主要函式：**

- `ToSnake(s string) string` - 將字串轉為 snake_case
- `EscapeBackslash(s string) string` - 將 \ 轉為 \\
- `UnescapeBackslash(s string) string` - 將 \\ 還原為 \
- `IsEmpty(s string) bool` - 判斷是否為空
- `Truncate(s string, maxLen int) string` - 截斷字串

---

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

---

### slicex - 切片操作

提供泛型切片操作函式。

```go
import "github.com/vincent119/commons/slicex"

slicex.Contains([]int{1, 2, 3}, 2)  // true
slicex.IndexOf([]string{"a", "b"}, "b")  // 1
slicex.Filter([]int{1, 2, 3, 4}, func(n int) bool { return n%2 == 0 })  // [2, 4]
slicex.Map([]int{1, 2}, func(n int) string { return fmt.Sprint(n) })  // ["1", "2"]
```

---

### timex - 時間處理

提供時區安全的時間操作函式。

```go
import "github.com/vincent119/commons/timex"

timex.NowUTC()                           // UTC 時間
timex.StartOfDay(time.Now(), time.Local) // 當天零點
timex.TimeStampUTC()                     // "2025-12-19T10:30:00.000Z"
timex.DateStamp()                        // "2025-12-19"
```

---

### uuidx - UUID 工具

```go
import "github.com/vincent119/commons/uuidx"

uuidx.NewUUID()      // 產生 UUID v4
uuidx.IsValidUUID("550e8400-e29b-41d4-a716-446655440000")  // true
```

---

### cryptox - 加密工具

```go
import "github.com/vincent119/commons/cryptox"

cryptox.MD5Hash("password")   // MD5 雜湊
cryptox.SHA256Hash("data")    // SHA256 雜湊
```

---

### validatorx - 驗證工具

```go
import "github.com/vincent119/commons/validatorx"

validatorx.IsEmail("user@example.com")  // true
validatorx.IsMobile("0912345678")       // true
validatorx.IsIPv4("192.168.1.1")        // true
```

---

### httpx/resp - HTTP 回應結構

```go
import "github.com/vincent119/commons/httpx/resp"

resp.Error{Code: 401, Message: "unauthorized"}
resp.Health{Status: "ok"}
```


---

### graceful - 優雅關機

提供應用程式生命週期管理，包含訊號監聽、資源釋放與超時控制 (支援 log/slog)。

[詳細文件](./graceful/README.md)

```go
import "github.com/vincent119/commons/graceful"

graceful.Run(task,
    graceful.WithLogger(slog.Default()), // 使用標準 slog
    graceful.WithCleanup(func(ctx context.Context) error { ... }),
)
```

---

## 開發指令

```bash
make help          # 顯示所有可用指令
make tidy          # 整理依賴
make fmt           # 格式化程式碼
make check         # 程式碼檢查（vet + lint）
make test          # 執行測試
make coverage      # 顯示覆蓋率報告
make coverage-html # 產生 HTML 覆蓋率報告
make bench         # 效能基準測試
make clean         # 清理產生的檔案
```

## 系統需求

- Go 1.21+

## 依賴

- `github.com/google/uuid` v1.6.0

## 授權

MIT License

## 貢獻

歡迎提交 Pull Request 或回報 Issue！
