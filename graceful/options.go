package graceful

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"time"
)

// Option 定義用來設定 graceful shutdown manager 的函式。
type Option func(*options)

type options struct {
	shutdownTimeout time.Duration
	logger          *slog.Logger
	cleaners        []Cleaner
}

// defaultOptions 回傳預設選項。
func defaultOptions() *options {
	return &options{
		shutdownTimeout: 30 * time.Second,
		logger:          slog.Default(),
		cleaners:        make([]Cleaner, 0),
	}
}

// WithTimeout 設定關機流程的超時時間。
// 若清理任務執行的時間超過此時長，可能會被取消。
// 預設為 30 秒。
func WithTimeout(d time.Duration) Option {
	return func(o *options) {
		o.shutdownTimeout = d
	}
}

// WithLogger 設定 manager 使用的 logger。
// 接受 *slog.Logger。
func WithLogger(l *slog.Logger) Option {
	return func(o *options) {
		if l != nil {
			o.logger = l
		}
	}
}

// WithCleanup 增加一個在關機期間執行的清理函式。
// 清理函式將按照加入的順序執行。
func WithCleanup(c Cleaner) Option {
	return func(o *options) {
		if c != nil {
			o.cleaners = append(o.cleaners, c)
		}
	}
}

// WithCloser 增加一個在關機期間關閉的 io.Closer。
// Close 方法將會在 Cleaner 包裝器內被呼叫。
// 注意：由於 io.Closer 不接受 context，若 Close 阻塞超過 shutdown timeout，
// 管理器會放棄等待並返回超時錯誤，但底層的 Close 操作仍會在背景執行直到完成或進程結束。
func WithCloser(c io.Closer) Option {
	return func(o *options) {
		if c != nil {
			o.cleaners = append(o.cleaners, func(ctx context.Context) error {
				done := make(chan error, 1)
				go func() {
					done <- c.Close()
				}()

				select {
				case err := <-done:
					return err
				case <-ctx.Done():
					return fmt.Errorf("closer (%T) timed out: %w", c, ctx.Err())
				}
			})
		}
	}
}
