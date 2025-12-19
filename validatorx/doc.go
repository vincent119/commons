// Package validatorx 提供常用的格式驗證函式。
//
// 所有驗證函式回傳 bool，適合用於表單驗證或 API 輸入檢查。
//
// # Email 驗證
//
//	valid := validatorx.IsEmail("user@example.com") // true
//	valid := validatorx.IsEmail("invalid")          // false
//
// # 手機號驗證（台灣格式）
//
//	valid := validatorx.IsMobile("0912345678") // true
//
// # UUID 驗證
//
//	valid := validatorx.IsUUID("550e8400-e29b-41d4-a716-446655440000") // true
//
// # IP 位址驗證
//
//	valid := validatorx.IsIPv4("192.168.1.1")   // true
//	valid := validatorx.IsIPv6("2001:db8::1")   // true
//
// # URL 驗證
//
//	valid := validatorx.IsURL("https://example.com") // true
//
// # 日期時間驗證
//
//	valid := validatorx.IsDate("2025-12-19")       // true
//	valid := validatorx.IsTime("10:30:00")         // true
//
// # 密碼強度驗證
//
//	valid := validatorx.IsPassword("Abc123!@#") // true（需包含大小寫、數字、特殊字元）
package validatorx
