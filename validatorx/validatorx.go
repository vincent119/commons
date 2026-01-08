package validatorx

import "regexp"

// IsEmail 驗證 email 格式
func IsEmail(email string) bool {
	re := regexp.MustCompile(`^[\w\.\-]+@[\w\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// IsMobile 驗證手機號格式（台灣簡易版，09 開頭共 10 碼）。
func IsMobile(mobile string) bool {
	re := regexp.MustCompile(`^09\d{8}$`)
	return re.MatchString(mobile)
}

// IsUUID 驗證 UUID v4 格式：8-4-4-4-12 的十六進位字串。
func IsUUID(u string) bool {
	re := regexp.MustCompile(`^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[1-5][a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$`)
	return re.MatchString(u)
}

// IsIPv4 驗證 IPv4 格式（0-255.0-255.0-255.0-255）。
func IsIPv4(ip string) bool {
	re := regexp.MustCompile(`^(25[0-5]|2[0-4]\d|[0-1]?\d?\d)(\.(25[0-5]|2[0-4]\d|[0-1]?\d?\d)){3}$`)
	return re.MatchString(ip)
}

// IsIPv6 驗證 IPv6 簡易格式（完整支援需 net.ParseIP）。
func IsIPv6(ip string) bool {
	re := regexp.MustCompile(`^([0-9a-fA-F]{1,4}:){7}[0-9a-fA-F]{1,4}$`)
	return re.MatchString(ip)
}

// IsURL 驗證 URL 格式（http/https）。
func IsURL(url string) bool {
	re := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
	return re.MatchString(url)
}

// IsDate 驗證日期格式（YYYY-MM-DD）。
func IsDate(date string) bool {
	re := regexp.MustCompile(`^(19|20)\\d{2}-(0[1-9]|1[0-2])-(0[1-9]|[12]\\d|3[01])$`)
	return re.MatchString(date)
}

// IsTime 驗證時間格式（HH:MM:SS，24 小時制）。
func IsTime(timeStr string) bool {
	re := regexp.MustCompile(`^([01]\\d|2[0-3]):[0-5]\\d:[0-5]\\d$`)
	return re.MatchString(timeStr)
}

// IsPassword 驗證密碼強度：至少 8 碼，需包含大小寫字母與數字。
func IsPassword(pwd string, max int) bool {
	if len(pwd) < max {
		return false
	}
	hasLower := false
	hasUpper := false
	hasDigit := false
	for _, c := range pwd {
		switch {
		case c >= 'a' && c <= 'z':
			hasLower = true
		case c >= 'A' && c <= 'Z':
			hasUpper = true
		case c >= '0' && c <= '9':
			hasDigit = true
		}
	}
	return hasLower && hasUpper && hasDigit
}
