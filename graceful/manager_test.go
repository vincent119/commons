package graceful

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"
)

func TestRun_GracefulShutdown(t *testing.T) {
	// 模擬一個立即返回的任務（模擬工作完成）
	task := func(ctx context.Context) error {
		return nil
	}

	// 我們想要驗證當任務結束時，清理函式是否被呼叫。
	cleanupCalled := false
	cleanup := func(ctx context.Context) error {
		cleanupCalled = true
		return nil
	}

	err := Run(task, WithCleanup(cleanup))
	if err != nil {
		t.Errorf("預期無錯誤，但得到 %v", err)
	}
	if !cleanupCalled {
		t.Error("任務正常退出時應該呼叫清理函式")
	}

	// 測試錯誤案例
	taskErr := func(ctx context.Context) error {
		return errors.New("boom")
	}

	err = Run(taskErr, WithCleanup(cleanup))
	if err == nil || err.Error() != "boom" {
		t.Errorf("預期 boom 錯誤，但得到 %v", err)
	}
	if !cleanupCalled {
		t.Error("即使任務發生錯誤，也應該呼叫清理函式")
	}
}

func TestRun_CleanupTimeout(t *testing.T) {
	// 使用預設的 logger (或者可以創建一個不輸出的 logger)
	logger := slog.Default()

	task := func(ctx context.Context) error {
		return nil
	}

	slowCleanup := func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return nil
		}
	}

	start := time.Now()
	// 設定清理超時 1ms，但清理需要 100ms
	err := Run(task, WithLogger(logger), WithCleanup(slowCleanup), WithTimeout(1*time.Millisecond))
	duration := time.Since(start)

	if duration > 50*time.Millisecond {
		t.Errorf("清理應該超時，但花費了 %v", duration)
	}

	// 清理回傳的錯誤取決於實作。
	// 在 manager.go 中：
	// if cErr := c(shutdownCtx); cErr != nil { log; append err }
	// slowCleanup 回傳 ctx.Err() (DeadlineExceeded)
	if err == nil {
		t.Error("預期因清理失敗而產生錯誤")
	}
}

func TestWithCloser(t *testing.T) {
	m := &mockCloser{}
	task := func(ctx context.Context) error { return nil }

	Run(task, WithCloser(m))

	if !m.closed {
		t.Error("closer 應該被關閉")
	}
}

type mockCloser struct {
	closed bool
}

func (m *mockCloser) Close() error {
	m.closed = true
	return nil
}

func TestRun_CleanupOrder(t *testing.T) {
	var executionOrder []int

	task := func(ctx context.Context) error { return nil }

	cleanup1 := func(ctx context.Context) error {
		executionOrder = append(executionOrder, 1)
		return nil
	}
	cleanup2 := func(ctx context.Context) error {
		executionOrder = append(executionOrder, 2)
		return nil
	}

	// 註冊順序: 1, then 2
	// 預期執行順序: 2, then 1 (LIFO)
	Run(task, WithCleanup(cleanup1), WithCleanup(cleanup2))

	if len(executionOrder) != 2 {
		t.Fatalf("預期執行 2 個 cleanup，但執行了 %d 個", len(executionOrder))
	}
	if executionOrder[0] != 2 || executionOrder[1] != 1 {
		t.Errorf("Cleanup 執行順序錯誤。預期 [2, 1]，得到 %v", executionOrder)
	}
}
