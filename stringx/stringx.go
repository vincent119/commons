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