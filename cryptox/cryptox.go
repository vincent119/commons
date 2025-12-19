package cryptox

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
)

// MD5Hash 回傳字串的 MD5 雜湊
func MD5Hash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}

// SHA256Hash 回傳字串的 SHA256 雜湊
func SHA256Hash(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
