// Package pathx 提供路徑處理相關的工具函式。
//
// # 路徑分隔符正規化
//
// 將 Windows 風格的反斜線轉換為正斜線，實現跨平台一致性：
//
//	path := pathx.NormalizePathSeparator("a\\b\\c")
//	// "a/b/c"
//
// 適用場景：
//   - 跨平台路徑處理
//   - URL 路徑建構
//   - 檔案系統路徑統一
package pathx
