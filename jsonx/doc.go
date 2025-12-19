// Package jsonx 提供 JSON 處理相關的工具函式。
//
// # JSON 字串跳脫
//
// 處理 JSON 字串中的特殊字元：
//
//	escaped := jsonx.EscapeJSON("line1\nline2")
//	// "line1\\nline2"
//
// 跳脫處理的字元：
//   - 反斜線 \ → \\
//   - 雙引號 " → \"
//   - 換行符 \n → \\n
//   - 回車符 \r → \\r
//   - Tab \t → \\t
//
// 適用場景：
//   - 手動建構 JSON 字串
//   - Log 輸出格式化
//   - 字串安全處理
package jsonx
