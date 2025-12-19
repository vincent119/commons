// Package sqlx 提供 SQL 查詢相關的工具函式。
//
// # LIKE 查詢
//
// 轉義 LIKE 特殊字元（%, _, \）：
//
//	escaped := sqlx.EscapeLikeQuery("50%_off")
//	// "50\%\_off"
//
// 建構 LIKE 查詢值（含通配符）：
//
//	like := sqlx.BuildLikeQueryValue("test", sqlx.LikePosBoth)
//	// "%test%"
//
// 位置常數（匹配語意，而非 % 放置位置）：
//   - LikePosStart: 前綴匹配（"value%"）
//   - LikePosEnd:   後綴匹配（"%value"）
//   - LikePosBoth:  包含匹配（"%value%"）
//
// 建議搭配 ESCAPE 子句使用：
//
//	query := "WHERE name LIKE ? " + sqlx.LikeEscapeClause()
//
// # SQL 字串跳脫
//
// 基礎 SQL 字串 escape（注意：不能取代 prepared statement）：
//
//	escaped := sqlx.EscapeSQLString("O'Reilly")
//	// "O\'Reilly"
//
// # Log 格式化
//
// 壓縮空白並移除雙重轉義，方便寫 log：
//
//	formatted := sqlx.FormatSQLForLog("SELECT * FROM   users")
//	// "SELECT * FROM users"
package sqlx
