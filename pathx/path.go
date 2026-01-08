package pathx

import "strings"

// NormalizePathSeparator 將路徑分隔符標準化為 Unix 風格（正斜線）。
func NormalizePathSeparator(path string) string {
	// 將 Windows 風格的反斜線替換為正斜線
	return strings.ReplaceAll(path, "\\", "/")
}
