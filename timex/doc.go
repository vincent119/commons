// Package timex 提供時區安全的時間處理工具函式。
//
// # 取得時間
//
// 取得 UTC 時間：
//
//	utc := timex.NowUTC()
//
// 取得某天的零點（指定時區）：
//
//	start := timex.StartOfDay(time.Now(), time.Local)
//
// # 時間截斷
//
// 截斷時間至指定粒度：
//
//	truncated := timex.TruncateTo(time.Now(), time.Hour)
//
// # 格式化與解析
//
// 格式化時間：
//
//	s := timex.FormatTime(time.Now(), "2006-01-02")
//
// 解析時間字串：
//
//	t, err := timex.ParseTime("2025-12-19", "2006-01-02")
//
// # 時間戳
//
// 各種格式的時間戳：
//
//	timex.TimeStamp()       // 本地時間戳
//	timex.TimeStampUTC()    // UTC 時間戳 "2025-12-19T10:30:00.000Z"
//	timex.DateStamp()       // 日期 "2025-12-19"
//	timex.UnixTimeStamp()   // Unix 秒數
//	timex.UnixMilliStamp()  // Unix 毫秒數
//	timex.TimeOnlyStamp()   // 僅時間 "10:30:00"
package timex
