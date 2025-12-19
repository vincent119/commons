package pathx


import "strings"

func NormalizePathSeparator(path string) string {
	// 將 Windows 風格的反斜線替換為正斜線
	return strings.ReplaceAll(path, "\\", "/")
}
