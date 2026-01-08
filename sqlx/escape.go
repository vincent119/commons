package sqlx

import "strings"

// LikePosStart/End/Both 保留原有名稱（向後相容）。
// 注意：這裡的 start/end 是「匹配型態」而不是「% 放置位置」：
//   - LikePosStart: value%   （前綴匹配，字串以 value 開頭）
//   - LikePosEnd:   %value   （後綴匹配，字串以 value 結尾）
//   - LikePosBoth:  %value%  （包含匹配）
const (
	LikePosStart = "start" // value%
	LikePosEnd   = "end"   // %value
	LikePosBoth  = "both"  // %value%
)

// 提供更直覺的別名常數（不影響原常數）。
const (
	LikeMatchPrefix = LikePosStart // value%
	LikeMatchSuffix = LikePosEnd   // %value
	LikeMatchBoth   = LikePosBoth  // %value%
)

// LikeEscapeChar 定義 LIKE 轉義字元（預設使用反斜線 \）。
// 實務上建議搭配 WHERE ... LIKE ? ESCAPE '\' 使用，以確保跨 DB 行為一致。
const LikeEscapeChar = `\`

// EscapeLikeQuery 轉義 LIKE 查詢特殊字元（%, _, \）。
// 目的：讓輸入字串中的 %, _ 不再具有 LIKE 通配語意。
// 注意：跨 DB 時，請搭配 LikeEscapeClause() 產生的 ESCAPE 子句一起用。
func EscapeLikeQuery(s string) string {
	// 先處理反斜線，避免後續替換產生二次干擾
	s = strings.ReplaceAll(s, `\`, `\\`) // 將 \ 變成 \\

	// 將 LIKE 的通配符 % 與 _ 轉義
	s = strings.ReplaceAll(s, `%`, `\%`) // 將 % 變成 \%
	s = strings.ReplaceAll(s, `_`, `\_`) // 將 _ 變成 \_

	// 回傳轉義後字串
	return s
}

// BuildLikeQueryValue 產生 LIKE 查詢字串（含通配符）。
// position 建議用 LikePosStart/LikePosEnd/LikePosBoth 或 LikeMatch*。
//   - LikePosStart (value%): 前綴匹配（以 value 開頭）
//   - LikePosEnd   (%value): 後綴匹配（以 value 結尾）
//   - LikePosBoth  (%value%): 包含匹配
func BuildLikeQueryValue(value string, position string) string {
	// 先對輸入做 LIKE 特殊字元轉義
	escaped := EscapeLikeQuery(value)

	// 根據匹配型態加上通配符
	switch position {
	case LikePosStart:
		return escaped + `%` // value%
	case LikePosEnd:
		return `%` + escaped // %value
	case LikePosBoth:
		return `%` + escaped + `%` // %value%
	default:
		// 若未指定，回傳純轉義字串（不加通配符）
		return escaped
	}
}

// LikeEscapeClause 回傳 SQL 的 ESCAPE 子句字串。
// 用法示例：
//
//	q := "WHERE col LIKE ? " + sqlx.LikeEscapeClause()
//	args := []any{sqlx.BuildLikeQueryValue(input, sqlx.LikePosBoth)}
//
// 注意：不同 DB 對 ESCAPE 的支援程度不同，但多數主流（MySQL/Postgres/SQLite）可用。
// 若你們 DB 不支援，可在專案層決定不要加這段。
func LikeEscapeClause() string {
	// 回傳固定的 ESCAPE '\'
	return `ESCAPE '\'`
}

// EscapeSQLString 基礎 SQL 字串 escape。
// 注意：不能取代 prepared statement；僅建議用於 log 或「非使用者輸入」的固定字串拼接。
func EscapeSQLString(s string) string {
	// 將反斜線轉義，避免後續跳脫序列混亂
	s = strings.ReplaceAll(s, `\`, `\\`) // \ -> \\

	// 轉義單引號與雙引號（視 DB/語境可能不同，這裡提供最基本處理）
	s = strings.ReplaceAll(s, `'`, `\'`) // ' -> \'
	s = strings.ReplaceAll(s, `"`, `\"`) // " -> \"

	// 回傳處理後字串
	return s
}

// UnescapeBackslash 還原已轉義的反斜線（通常供 log 格式化使用）。
func UnescapeBackslash(s string) string {
	// 將 \\ 還原成 \
	return strings.ReplaceAll(s, `\\`, `\`)
}

// FormatSQLForLog 壓縮空白並移除雙重轉義，方便寫 log。
// 注意：這是「可讀性」用的格式化，不保證是可執行 SQL。
func FormatSQLForLog(sql string) string {
	// 將多餘空白折疊成單一空白
	sql = strings.Join(strings.Fields(sql), " ")

	// 移除 log 中常見的雙重反斜線轉義
	return UnescapeBackslash(sql)
}
