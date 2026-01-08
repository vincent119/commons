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
	task := func(_ context.Context) error {
		return nil
	}

	// 我們想要驗證當任務結束時，清理函式是否被呼叫。
	cleanupCalled := false
	cleanup := func(_ context.Context) error {
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
	taskErr := func(_ context.Context) error {
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

	task := func(_ context.Context) error {
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
	task := func(_ context.Context) error { return nil }

	if err := Run(task, WithCloser(m)); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if !m.closed {
		t.Error("closer should be closed")
	}
}

func TestWithClosers(t *testing.T) {
	m1 := &mockCloser{}
	m2 := &mockCloser{}
	m3 := &mockCloser{}
	task := func(_ context.Context) error { return nil }

	if err := Run(task, WithClosers(m1, m2, m3)); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if !m1.closed || !m2.closed || !m3.closed {
		t.Error("all closers should be closed")
	}
}

func TestWithClosers_LIFO_Order(t *testing.T) {
	var closeOrder []int

	c1 := &orderedCloser{id: 1, order: &closeOrder}
	c2 := &orderedCloser{id: 2, order: &closeOrder}
	c3 := &orderedCloser{id: 3, order: &closeOrder}

	task := func(_ context.Context) error { return nil }

	// 註冊順序: 1, 2, 3
	// 預期執行順序: 3, 2, 1 (LIFO)
	if err := Run(task, WithClosers(c1, c2, c3)); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if len(closeOrder) != 3 {
		t.Fatalf("expected 3 closers to be closed, but only closed %d", len(closeOrder))
	}
	if closeOrder[0] != 3 || closeOrder[1] != 2 || closeOrder[2] != 1 {
		t.Errorf("Closer execution order error. expected [3, 2, 1], got %v", closeOrder)
	}
}

func TestWithClosers_NilHandling(t *testing.T) {
	m := &mockCloser{}
	task := func(_ context.Context) error { return nil }

	// 傳入 nil 不應該 panic
	if err := Run(task, WithClosers(m, nil, m)); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if !m.closed {
		t.Error("non-nil closer should be closed")
	}
}

type mockCloser struct {
	closed bool
}

type orderedCloser struct {
	id    int
	order *[]int
}

func (m *mockCloser) Close() error {
	m.closed = true
	return nil
}

func (o *orderedCloser) Close() error {
	*o.order = append(*o.order, o.id)
	return nil
}

func TestRun_CleanupOrder(t *testing.T) {
	var executionOrder []int

	task := func(_ context.Context) error { return nil }

	cleanup1 := func(_ context.Context) error {
		executionOrder = append(executionOrder, 1)
		return nil
	}
	cleanup2 := func(_ context.Context) error {
		executionOrder = append(executionOrder, 2)
		return nil
	}

	// 註冊順序: 1, then 2
	// 預期執行順序: 2, then 1 (LIFO)
	if err := Run(task, WithCleanup(cleanup1), WithCleanup(cleanup2)); err != nil {
		t.Errorf("Run() error = %v", err)
	}

	if len(executionOrder) != 2 {
		t.Fatalf("expected 2 cleaners to be executed, but only executed %d", len(executionOrder))
	}
	if executionOrder[0] != 2 || executionOrder[1] != 1 {
		t.Errorf("Cleanup execution order error. expected [2, 1], got %v", executionOrder)
	}
}
