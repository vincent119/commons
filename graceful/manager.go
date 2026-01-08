package graceful

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

// Task 代表一個長時間運行的任務，需監聽 context 取消訊號。
// 當 context 完成或遇到致命錯誤時應返回。
type Task func(ctx context.Context) error

// Cleaner 定義在關機期間釋放資源的函式。
// 重要：Cleaner 必須尊重 ctx 的超時設定。若遇到 ctx.Done()，應立即返回 ctx.Err()，
// 否則會阻塞後續的清理工作。
type Cleaner func(ctx context.Context) error

// Run 執行給定的任務並在收到系統訊號時處理優雅關機。
// 它監聽 SIGINT 和 SIGTERM。
func Run(task Task, opts ...Option) error {
	o := defaultOptions()
	for _, opt := range opts {
		opt(o)
	}

	// 1. 設定訊號 Context
	// NotifyContext 建立一個父 context 的副本，當收到列出的訊號、
	// 呼叫返回的 stop 函式，或父 context 的 Done 通道關閉時（視何者先發生），
	// 該副本會被標記為完成（其 Done 通道關閉）。
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// 2. 執行任務
	o.logger.Info("starting task")
	startTime := time.Now()

	// 執行任務。預期任務會阻塞直到完成或 ctx 完成。
	err := task(ctx)

	// 記錄任務退出
	duration := time.Since(startTime)
	if err != nil {
		o.logger.Error("task exited with error", "error", err, "duration", duration)
	} else {
		o.logger.Info("task exited successfully", "duration", duration)
	}

	// 3. 執行清理
	// 我們建立一個新的 context 進行清理，因為訊號 context 已經完成。
	o.logger.Info("starting shutdown cleanup", "timeout", o.shutdownTimeout)

	shutdownCtx, cancel := context.WithTimeout(context.Background(), o.shutdownTimeout)
	defer cancel()

	var cleanupErrors []error
	// 執行清理 (LIFO: 後進先出)
	// 確保依賴關係較高的資源（通常較晚註冊）先被釋放
	for i := len(o.cleaners) - 1; i >= 0; i-- {
		c := o.cleaners[i]
		if cErr := c(shutdownCtx); cErr != nil {
			o.logger.Error("cleanup failed", "error", cErr)
			cleanupErrors = append(cleanupErrors, cErr)
		}
	}

	if len(cleanupErrors) > 0 {
		// 如果任務本身有錯，也包含在返回的錯誤中，避免錯誤吞沒
		if err != nil {
			cleanupErrors = append([]error{err}, cleanupErrors...)
		}
		return errors.Join(cleanupErrors...)
	}

	o.logger.Info("shutdown complete")
	return err
}

// HTTPTask 將 http.Server 包裝為 graceful.Task。
// 它在 goroutine 中啟動伺服器並等待 context 完成。
// 注意：您通常需要註冊一個清理函式來關閉伺服器，例如：
// graceful.WithCleanup(func(ctx context.Context) error { return srv.Shutdown(ctx) })
// 或者若有偏好，可使用處理此邏輯的輔助函式。
//
// 然而，為了保持簡單與可組合性，此 helper 僅處理「運行」部分。
// 使用者負責確保在清理階段呼叫 srv.Shutdown。
func HTTPTask(srv *http.Server) Task {
	return func(ctx context.Context) error {
		errCh := make(chan error, 1)
		go func() {
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				errCh <- err
			}
			close(errCh)
		}()

		select {
		case err := <-errCh:
			// 伺服器失敗或停止
			return err
		case <-ctx.Done():
			// 收到訊號。返回以允許開始清理階段。
			// 此時伺服器仍在運行！
			// 清理階段（應包含 srv.Shutdown）將會停止它。
			return nil
		}
	}
}
