package timex

import (
	"testing"
	"time"
)

func TestStartOfDay(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Taipei")
	// 2025-08-19 10:00:00+08 測試：零點應回到同一天 00:00:00+08（再轉 UTC）
	in := time.Date(2025, 8, 19, 10, 0, 0, 0, loc)
	got := StartOfDay(in, loc)

	// 期望的當地零點（轉 UTC）
	want := time.Date(2025, 8, 18, 16, 0, 0, 0, time.UTC) // 台北比 UTC +8
	if !got.Equal(want) {
		t.Fatalf("StartOfDay() got %v; want %v", got, want)
	}
}
