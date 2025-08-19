package timex

import (
	"time"
)

// NowUTC 取得目前 UTC 時間。
// 使用 time.Now().UTC()，避免隱含時區引發序列化/跨服務差異。
func NowUTC() time.Time {
	return time.Now().UTC()
}

// StartOfDay 回傳指定時區下某時刻的「零點」時間（當地日界）
// 1) 先將 t 轉到指定時區 loc
// 2) 取當地年月日
// 3) 建立當地零點，再轉回 UTC（便於儲存/比較）
func StartOfDay(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)                                        // 轉到目標時區
	y, m, d := local.Date()                                   // 取當地年月日
	zeroLocal := time.Date(y, m, d, 0, 0, 0, 0, loc)          // 當地零點
	return zeroLocal.UTC()                                    // 標準化成 UTC
}

// TruncateTo 將時間截斷至指定粒度（如分鐘/小時），以 UTC 作業避免跨時區差異。
func TruncateTo(t time.Time, d time.Duration) time.Time {
	return t.UTC().Truncate(d)
}