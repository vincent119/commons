package stringx

import "strings"

// ToSnake 將字串轉為 snake_case（簡化版）。
func ToSnake(s string) string {
	if s == "" {
		return s
	}

	var b strings.Builder
	b.Grow(len(s) * 2)

	lastWasUnderscore := false
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			if i > 0 && !lastWasUnderscore {
				b.WriteByte('_')
			}
			r += 'a' - 'A'
			b.WriteRune(r)
			lastWasUnderscore = false
			continue
		}

		if r == ' ' || r == '-' {
			if !lastWasUnderscore {
				b.WriteByte('_')
				lastWasUnderscore = true
			}
			continue
		}

		b.WriteRune(r)
		lastWasUnderscore = (r == '_')
	}

	out := b.String()
	return strings.TrimRight(out, "_")
}

// EscapeBackslash 將單反斜線替換為雙反斜線（通用字串處理）。
func EscapeBackslash(s string) string {
	return strings.ReplaceAll(s, "\\", "\\\\")
}

// UnescapeBackslash 還原已轉義的反斜線（通用字串處理）。
func UnescapeBackslash(s string) string {
	return strings.ReplaceAll(s, "\\\\", "\\")
}

// IsEmpty 檢查字串是否為空（忽略空白）。
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// Truncate 截斷字串到指定長度（以 byte 計，UTF-8 可能切到半個 rune）。
func Truncate(s string, maxLen int) string {
	if maxLen <= 0 {
		return ""
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
