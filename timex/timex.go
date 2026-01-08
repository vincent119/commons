package timex

import (
	"fmt"
	"time"
)

// NowUTC 取得目前 UTC 時間。
// 使用 time.Now().UTC()，避免隱含時區引發序列化/跨服務差異。
func NowUTC() time.Time {
	return time.Now().UTC()
}

// StartOfDay 回傳指定時區下某時刻的「零點」時間（當地日界）。
// 1) 先將 t 轉到指定時區 loc
// 2) 取當地年月日
// 3) 建立當地零點，再轉回 UTC（便於儲存/比較）
func StartOfDay(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)                               // 轉到目標時區
	y, m, d := local.Date()                          // 取當地年月日
	zeroLocal := time.Date(y, m, d, 0, 0, 0, 0, loc) // 當地零點
	return zeroLocal.UTC()                           // 標準化成 UTC
}

// TruncateTo 將時間截斷至指定粒度（如分鐘/小時），以 UTC 作業避免跨時區差異。
func TruncateTo(t time.Time, d time.Duration) time.Time {
	return t.UTC().Truncate(d)
}

// FormatTime 格式化時間為字串。
func FormatTime(t time.Time, layout string) string {
	return t.Format(layout)
}

// ParseTime 解析字串為時間。
func ParseTime(str, layout string) (time.Time, error) {
	return time.Parse(layout, str)
}

// TimeStamp 取得目前時間的時間戳。
func TimeStamp() string {
	return time.Now().Format("2006-01-02T15:04:05.000Z07:00")
}

// TimeStampUTC 取得目前時間的 UTC 時間戳。
func TimeStampUTC() string {
	return time.Now().UTC().Format("2006-01-02T15:04:05.000Z")
}

// DateStamp 取得目前日期（YYYY-MM-DD）。
func DateStamp() string {
	return time.Now().Format("2006-01-02")
}

// UnixTimeStamp 取得目前時間的 Unix 秒數字串。
func UnixTimeStamp() string {
	return fmt.Sprintf("%d", time.Now().Unix())
}

// UnixMilliStamp 取得目前時間的 Unix 毫秒數字串。
func UnixMilliStamp() string {
	return fmt.Sprintf("%d", time.Now().UnixMilli())
}

// TimeOnlyStamp 取得目前時間（HH:MM:SS）。
func TimeOnlyStamp() string {
	return time.Now().Format("15:04:05")
}

// WithZoneTimeStamp 回傳指定時區的當下時間戳（ISO 8601 格式）。
func WithZoneTimeStamp(loc *time.Location) string {
	return time.Now().In(loc).Format("2006-01-02T15:04:05.000Z07:00")
}
