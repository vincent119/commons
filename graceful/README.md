# Graceful Run & Shutdown Library (`graceful`)

`graceful` 是一個用於 Go 應用程式的標準生命週期管理套件。它提供了一致的介面來處理：

- 系統訊號監聽 (SIGINT, SIGTERM)
- Context 管理與傳遞
- 啟動與關機任務
- 資源清理機制的超時控制 (Graceful Shutdown)
- 結構化日誌支援 (log/slog)

## 安裝

```bash
go get github.com/vincent119/commons/graceful
```

## 使用範例

### 基礎用法 (HTTP Server)

```go
package main

import (
    "context"
    "log/slog"
    "net/http"
    "os"
    "time"

    "github.com/vincent119/commons/graceful"
)

func main() {
    // 初始化 Logger (可選，預設使用 slog.Default())
    logger := slog.New(slog.NewTextHandler(os.Stderr, nil))

    // 建立資源
    // db, _ := sql.Open(...)

    srv := &http.Server{Addr: ":8080"}

    // 啟動 Graceful Run
    err := graceful.Run(
        // 1. 主要任務：可以是 HTTP Task 或任何長時間運行的函式
        graceful.HTTPTask(srv),

        // 2. 設定 Logger (支援 *slog.Logger)
        graceful.WithLogger(logger),

        // 3. 設定 Shutdown Timeout
        graceful.WithTimeout(10*time.Second),

        // 4. 註冊清理函式 (Cleanup) - 執行順序為後進先出 (LIFO)
        // 關閉 HTTP Server (必要)
        graceful.WithCleanup(func(ctx context.Context) error {
            logger.Info("shutting down server...")
            return srv.Shutdown(ctx)
        }),

        // 5. 註冊資源關閉 (Closer)
        // 自動呼叫 Close()，適合資料庫連線等資源
        // graceful.WithCloser(db),
    )

    if err != nil {
        logger.Error("application exited with error", "error", err)
        os.Exit(1)
    }
}
```

### 通用 Worker 用法

支援任何符合 `func(ctx context.Context) error` 的任務。

```go
func MyWorker(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            // 收到訊號，退出迴圈
            return nil
        default:
            // 執行工作
            fmt.Println("working...")
            time.Sleep(1 * time.Second)
        }
    }
}

func main() {
    graceful.Run(MyWorker)
}
```

### 功能選項 (Options)

| Option | 說明 | 預設值 |
| :--- | :--- | :--- |
| `WithTimeout(d)` | Shutdown 階段的總體超時時間 | 30s |
| `WithLogger(l)` | 設定 logger (支援 `*slog.Logger`) | `slog.Default()` |
| `WithCleanup(f)` | 註冊清理函式 (LIFO 順序執行) | 無 |
| `WithCloser(c)` | 註冊 `io.Closer` 資源 | 無 |

## 注意事項

1. **清理順序**：Cleanup 函式採用 **LIFO (後進先出)** 順序執行。建議先註冊最底層資源（如 DB），再註冊上層服務（如 HTTP Server），確保關機時先停止服務再關閉資料庫。
2. **錯誤合併**：若主任務與清理工作皆發生錯誤，`graceful.Run` 會使用 `errors.Join` 返回所有錯誤。
3. **超時控制**：每個 Cleaner 必須尊重 `ctx` 的超時訊號 (`ctx.Done()`)，避免阻塞整體關機流程。
