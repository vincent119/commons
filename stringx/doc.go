// Package stringx 提供高效的字串處理工具函式。
//
// # 大小寫轉換
//
// 將 CamelCase 轉為 snake_case：
//
//	s := stringx.ToSnake("UserID") // "user_i_d"
//
// # SQL 跳脫
//
// 跳脫 SQL 字串中的特殊字元：
//
//	escaped := stringx.EscapeSQLString("O'Reilly") // "O''Reilly"
//
// 建構 LIKE 查詢值：
//
//	like := stringx.BuildLikeQueryValue("test", stringx.LikeContains)
//	// "%test%"
//
// # 路徑處理
//
// 統一路徑分隔符：
//
//	path := stringx.NormalizePathSeparator("a\\b\\c") // "a/b/c"
//
// # 字串工具
//
// 判斷是否為空：
//
//	stringx.IsEmpty("")      // true
//	stringx.IsEmpty("  ")    // true
//
// 截斷字串：
//
//	s := stringx.Truncate("hello world", 5) // "hello"
//
// JSON 跳脫：
//
//	escaped := stringx.EscapeJSON("line1\nline2")
package stringx
