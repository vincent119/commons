package stringx

import "strings"

// ToSnake 將字串轉為 snake_case（簡化版）。
// 設計目標：
//   - 針對常見 ASCII/大小寫場景最佳化，降低分配次數；
//   - 空字串快速回傳；
//   - 空白/連字號正規化為底線；
//   - 大寫前置於非首位時插入底線；
// 注意：縮寫（如 "ID"）會被拆成 i_d，若你要併成 id，請後續提供 ToSnakeSmart 以免破壞現有語義。
func ToSnake(s string) string {
	// 空字串直接回傳，避免不必要邏輯
	if s == "" {
		return s
	}

	// 使用 strings.Builder，預估容量以減少擴容（效能考量）
	var b strings.Builder
	b.Grow(len(s) * 2)

	lastWasUnderscore := false // 追蹤上一個是否為底線，避免連續底線
	for i, r := range s { // 逐 rune 處理，可兼容非 ASCII 字元
		// 大寫字母：若非首位，前置底線；再做 ASCII 快速小寫轉換
		if r >= 'A' && r <= 'Z' {
			if i > 0 && !lastWasUnderscore {
				b.WriteByte('_')
			}
			// ASCII 小寫化
			r += 'a' - 'A'
			b.WriteRune(r)
			lastWasUnderscore = false
			continue
		}

		// 空白或連字號規整為底線
		if r == ' ' || r == '-' {
			if !lastWasUnderscore {
				b.WriteByte('_')
				lastWasUnderscore = true
			}
			continue
		}

		// 其他字元原樣輸出
		b.WriteRune(r)
		lastWasUnderscore = (r == '_')
	}

	// 去尾端底線（若最後被規整為底線）
	out := b.String()
	if strings.HasSuffix(out, "_") {
		return strings.TrimRight(out, "_")
	}
	return out
}


// EscapeBackslash 處理字串中的反斜線，將單反斜線替換為雙反斜線
func EscapeBackslash(s string) string {
	return strings.ReplaceAll(s, "\\", "\\\\")
}

// UnescapeBackslash 還原已轉義的反斜線
func UnescapeBackslash(s string) string {
	return strings.ReplaceAll(s, "\\\\", "\\")
}

// EscapeSQLString 安全處理 SQL 字符串中的特殊字符
func EscapeSQLString(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\") // 替換反斜線為雙反斜線
	s = strings.ReplaceAll(s, "'", "\\'")   // 替換單引號
	s = strings.ReplaceAll(s, "\"", "\\\"") // 替換雙引號
	return s
}

// NormalizePathSeparator 標準化路徑分隔符（跨平台）
func NormalizePathSeparator(path string) string {
	// 將 Windows 風格的反斜線替換為正斜線
	return strings.ReplaceAll(path, "\\", "/")
}

// BuildLikeQueryValue 產生用於 LIKE 查詢的字串值，包含通配符
func BuildLikeQueryValue(value string, position string) string {
	// 處理輸入值中的特殊字元
	escaped := EscapeLikeQuery(value)

	// 根據位置增加通配符
	switch position {
	case "start":
		return escaped + "%"
	case "end":
		return "%" + escaped
	case "both":
		return "%" + escaped + "%"
	default:
		return escaped
	}
}

// EscapeLikeQuery 處理 LIKE 查詢中的特殊字元
func EscapeLikeQuery(s string) string {
	// 處理 SQL LIKE 查詢中的特殊字元 (%, _, \)
	s = strings.ReplaceAll(s, "\\", "\\\\") // 替換反斜線為雙反斜線
	s = strings.ReplaceAll(s, "%", "\\%")   // 替換百分號
	s = strings.ReplaceAll(s, "_", "\\_")   // 替換下劃線
	return s
}

// IsEmpty 檢查字串是否為空
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Truncate 截斷字串到指定長度
func Truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}

// EscapeJSON 處理JSON字串中的特殊字符
func EscapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}

// FormatSQLForLog 格式化SQL語句以便於日誌記錄
func FormatSQLForLog(sql string) string {
	// 移除多餘空白
	sql = strings.Join(strings.Fields(sql), " ")
	// 移除日誌中的雙重轉義
	return UnescapeBackslash(sql)
}
