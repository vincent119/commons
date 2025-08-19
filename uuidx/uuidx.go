package uuidx

import "github.com/google/uuid"

// NewUUID 產生新的 UUID v4 字串（隨機）
func NewUUID() string {
	return uuid.New().String()
}

// NewUUIDv4 產生新的 UUID v4 字串（隨機）
func NewUUIDv4() string {
	return uuid.NewString()
}

// NewUUIDv5 產生基於 namespace 與 name 的 UUID v5 字串
func NewUUIDv5(namespace uuid.UUID, name string) string {
	return uuid.NewSHA1(namespace, []byte(name)).String()
}

// IsValidUUID 驗證字串是否為合法 UUID
func IsValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}
