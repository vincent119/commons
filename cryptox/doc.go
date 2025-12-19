// Package cryptox 提供常用的加密與雜湊工具函式。
//
// # 雜湊函式
//
// MD5 雜湊（回傳十六進位字串）：
//
//	hash := cryptox.MD5Hash("password")
//
// SHA256 雜湊（回傳十六進位字串）：
//
//	hash := cryptox.SHA256Hash("data")
//
// # 安全提醒
//
// MD5 不應用於密碼儲存或安全敏感場景，建議使用 bcrypt 或 argon2。
// SHA256 適用於資料完整性驗證，但密碼儲存仍建議使用專用演算法。
package cryptox
