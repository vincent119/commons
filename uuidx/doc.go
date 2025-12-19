// Package uuidx 提供 UUID 產生與驗證的便利函式。
//
// 此套件封裝 github.com/google/uuid，提供更簡潔的 API。
//
// # 產生 UUID
//
// 產生 UUID v4（隨機）：
//
//	id := uuidx.NewUUID()    // 或 uuidx.NewUUIDv4()
//
// 產生 UUID v5（命名空間 + 名稱）：
//
//	id := uuidx.NewUUIDv5(uuid.NameSpaceDNS, "example.com")
//
// # 驗證 UUID
//
// 驗證字串是否為有效 UUID 格式：
//
//	valid := uuidx.IsValidUUID("550e8400-e29b-41d4-a716-446655440000") // true
//	valid := uuidx.IsValidUUID("invalid")                              // false
package uuidx
