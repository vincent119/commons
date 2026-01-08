// Package ipx 提供 IP 位址相關的通用工具函式。
//
// 此套件包含以下功能：
//   - IP 驗證：IsValidIP、IsIPv4、IsIPv6、IsPublicIP
//   - IP 轉換：IPv4ToUint32、Uint32ToIPv4、ExpandIPv6
//   - 網段工具：IsIPInCIDR、GetNetworkInfo
//   - 地理位置：GetLocationByIP（可整合 GeoIP2）
//   - 客戶端 IP 偵測：GetClientIP（支援 X-Forwarded-For、X-Real-IP）
//   - 本機 IP 取得：GetLocalIPs
package ipx

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strings"
)

// =============================================================================
// IP 驗證工具
// =============================================================================

// IsValidIP 驗證字串是否為有效的 IP 位址（支援 IPv4 與 IPv6）。
//
// 範例：
//
//	IsValidIP("192.168.1.1")     // true
//	IsValidIP("::1")             // true
//	IsValidIP("invalid")         // false
func IsValidIP(ip string) bool {
	return net.ParseIP(strings.TrimSpace(ip)) != nil
}

// IsIPv4 判斷字串是否為有效的 IPv4 位址。
//
// 範例：
//
//	IsIPv4("192.168.1.1")    // true
//	IsIPv4("::1")            // false（這是 IPv6）
//	IsIPv4("256.1.1.1")      // false（超出範圍）
func IsIPv4(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	// To4() 回傳非 nil 表示是 IPv4
	return parsed.To4() != nil
}

// IsIPv6 判斷字串是否為有效的 IPv6 位址。
//
// 注意：IPv4-mapped IPv6（如 ::ffff:192.168.1.1）會被視為 IPv6。
//
// 範例：
//
//	IsIPv6("2001:db8::1")            // true
//	IsIPv6("::ffff:192.168.1.1")     // true（IPv4-mapped IPv6）
//	IsIPv6("192.168.1.1")            // false
func IsIPv6(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	// 若 To4() 回傳 nil，且 parsed 非 nil，則為純 IPv6
	return parsed.To4() == nil
}

// IsPublicIP 判斷 IP 是否為公網 IP（非私有、非保留、非迴環）。
//
// 範例：
//
//	IsPublicIP("8.8.8.8")        // true
//	IsPublicIP("192.168.1.1")    // false（私有）
//	IsPublicIP("127.0.0.1")      // false（迴環）
func IsPublicIP(ip string) bool {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return false
	}
	return !isPrivateIP(parsed)
}

// =============================================================================
// IP 轉換工具
// =============================================================================

// IPv4ToUint32 將 IPv4 位址字串轉換為 uint32 整數。
//
// 此函式適用於需要以數值方式比較或儲存 IP 的場景。
// 若輸入非有效 IPv4，回傳 0 與錯誤。
//
// 範例：
//
//	IPv4ToUint32("192.168.1.1")   // 3232235777, nil
//	IPv4ToUint32("0.0.0.0")       // 0, nil
//	IPv4ToUint32("::1")           // 0, error（非 IPv4）
func IPv4ToUint32(ip string) (uint32, error) {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return 0, fmt.Errorf("無效的 IP 位址: %s", ip)
	}

	ip4 := parsed.To4()
	if ip4 == nil {
		return 0, fmt.Errorf("非 IPv4 位址: %s", ip)
	}

	// 使用 big-endian 將 4 bytes 轉為 uint32
	return binary.BigEndian.Uint32(ip4), nil
}

// Uint32ToIPv4 將 uint32 整數轉換為 IPv4 位址字串。
//
// 範例：
//
//	Uint32ToIPv4(3232235777)   // "192.168.1.1"
//	Uint32ToIPv4(0)            // "0.0.0.0"
func Uint32ToIPv4(n uint32) string {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, n)
	return ip.String()
}

// ExpandIPv6 將 IPv6 位址展開為完整的 8 組表示法。
//
// 此函式會將縮寫的 IPv6（如 ::1）展開為完整格式（如 0000:0000:0000:0000:0000:0000:0000:0001）。
// 若輸入非有效 IPv6，回傳空字串與錯誤。
//
// 範例：
//
//	ExpandIPv6("::1")                    // "0000:0000:0000:0000:0000:0000:0000:0001", nil
//	ExpandIPv6("2001:db8::1")            // "2001:0db8:0000:0000:0000:0000:0000:0001", nil
//	ExpandIPv6("192.168.1.1")            // "", error
func ExpandIPv6(ip string) (string, error) {
	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return "", fmt.Errorf("無效的 IP 位址: %s", ip)
	}

	// 確保是 IPv6（To16 會將 IPv4 也轉為 IPv6 格式，需排除）
	if parsed.To4() != nil {
		return "", fmt.Errorf("非 IPv6 位址: %s", ip)
	}

	ip6 := parsed.To16()
	if ip6 == nil {
		return "", fmt.Errorf("無法轉換為 IPv6: %s", ip)
	}

	// 格式化為完整的 8 組 16 進位表示
	groups := make([]string, 8)
	for i := 0; i < 8; i++ {
		groups[i] = fmt.Sprintf("%02x%02x", ip6[i*2], ip6[i*2+1])
	}

	return strings.Join(groups, ":"), nil
}

// =============================================================================
// 網段相關工具
// =============================================================================

// IsIPInCIDR 判斷指定 IP 是否在給定的 CIDR 網段內。
//
// 範例：
//
//	IsIPInCIDR("192.168.1.100", "192.168.1.0/24")   // true, nil
//	IsIPInCIDR("10.0.0.1", "192.168.1.0/24")        // false, nil
//	IsIPInCIDR("invalid", "192.168.1.0/24")         // false, error
func IsIPInCIDR(ip, cidr string) (bool, error) {
	parsedIP := net.ParseIP(strings.TrimSpace(ip))
	if parsedIP == nil {
		return false, fmt.Errorf("無效的 IP 位址: %s", ip)
	}

	_, ipNet, err := net.ParseCIDR(strings.TrimSpace(cidr))
	if err != nil {
		return false, fmt.Errorf("無效的 CIDR 格式: %s", cidr)
	}

	return ipNet.Contains(parsedIP), nil
}

// NetworkInfo 網段資訊結構。
type NetworkInfo struct {
	// Network 網路位址（如 192.168.1.0）
	Network string `json:"network"`

	// Broadcast 廣播位址（如 192.168.1.255），僅適用於 IPv4
	Broadcast string `json:"broadcast,omitempty"`

	// FirstHost 第一個可用主機位址
	FirstHost string `json:"first_host"`

	// LastHost 最後一個可用主機位址
	LastHost string `json:"last_host"`

	// TotalHosts 可用主機數量
	TotalHosts uint64 `json:"total_hosts"`

	// PrefixLength 前綴長度（如 24）
	PrefixLength int `json:"prefix_length"`

	// Netmask 子網路遮罩（如 255.255.255.0），僅適用於 IPv4
	Netmask string `json:"netmask,omitempty"`
}

// GetNetworkInfo 取得指定 CIDR 網段的詳細資訊。
//
// 回傳網路位址、廣播位址、可用主機數等資訊。
//
// 範例：
//
//	info, _ := GetNetworkInfo("192.168.1.0/24")
//	// info.Network = "192.168.1.0"
//	// info.Broadcast = "192.168.1.255"
//	// info.TotalHosts = 254
func GetNetworkInfo(cidr string) (*NetworkInfo, error) {
	_, ipNet, err := net.ParseCIDR(strings.TrimSpace(cidr))
	if err != nil {
		return nil, fmt.Errorf("無效的 CIDR 格式: %s", cidr)
	}

	// 取得前綴長度
	prefixLen, totalBits := ipNet.Mask.Size()

	info := &NetworkInfo{
		Network:      ipNet.IP.String(),
		PrefixLength: prefixLen,
	}

	// 計算可用主機數量
	hostBits := totalBits - prefixLen
	if hostBits >= 64 {
		// 超過 uint64 可表示範圍，設為最大值
		info.TotalHosts = math.MaxUint64
	} else if hostBits > 1 {
		// 扣除網路位址與廣播位址
		info.TotalHosts = (1 << hostBits) - 2
	} else {
		// /31 或 /32 網段
		info.TotalHosts = 1 << hostBits
	}

	// IPv4 特定處理
	if ip4 := ipNet.IP.To4(); ip4 != nil {
		// 計算廣播位址與子網路遮罩
		broadcast := make(net.IP, 4)
		netmask := make(net.IP, 4)

		for i := 0; i < 4; i++ {
			broadcast[i] = ip4[i] | ^ipNet.Mask[i]
			netmask[i] = ipNet.Mask[i]
		}

		info.Broadcast = broadcast.String()
		info.Netmask = netmask.String()

		// 計算第一個與最後一個可用主機
		if hostBits > 1 {
			firstHost := make(net.IP, 4)
			lastHost := make(net.IP, 4)
			copy(firstHost, ip4)
			copy(lastHost, broadcast)
			firstHost[3]++ // 網路位址 + 1
			lastHost[3]--  // 廣播位址 - 1
			info.FirstHost = firstHost.String()
			info.LastHost = lastHost.String()
		} else {
			info.FirstHost = ip4.String()
			info.LastHost = broadcast.String()
		}
	} else {
		// IPv6：計算第一個與最後一個可用主機
		ip6 := ipNet.IP.To16()
		if ip6 != nil {
			// 第一個主機（網路位址 + 1）
			firstHost := make(net.IP, 16)
			copy(firstHost, ip6)
			for i := 15; i >= 0; i-- {
				firstHost[i]++
				if firstHost[i] != 0 {
					break
				}
			}
			info.FirstHost = firstHost.String()

			// 最後一個主機（需計算末端位址 - 1）
			lastHost := make(net.IP, 16)
			for i := 0; i < 16; i++ {
				lastHost[i] = ip6[i] | ^ipNet.Mask[i]
			}
			// 減 1
			for i := 15; i >= 0; i-- {
				if lastHost[i] > 0 {
					lastHost[i]--
					break
				}
				lastHost[i] = 0xff
			}
			info.LastHost = lastHost.String()
		}
	}

	return info, nil
}

// =============================================================================
// 地理位置工具
// =============================================================================

// GeoIPProvider 定義 GeoIP 服務提供者介面。
//
// 實作此介面可整合不同的 GeoIP 服務（如 MaxMind GeoIP2、IP-API 等）。
type GeoIPProvider interface {
	// Lookup 根據 IP 位址查詢地理位置資訊
	Lookup(ip string) (*GeoLocation, error)
}

// GeoLocation 地理位置資訊結構。
type GeoLocation struct {
	// IP 查詢的 IP 位址
	IP string `json:"ip"`

	// Country 國家名稱
	Country string `json:"country"`

	// CountryCode 國家代碼（ISO 3166-1 alpha-2）
	CountryCode string `json:"country_code"`

	// Region 區域/省份名稱
	Region string `json:"region"`

	// City 城市名稱
	City string `json:"city"`

	// Latitude 緯度
	Latitude float64 `json:"latitude"`

	// Longitude 經度
	Longitude float64 `json:"longitude"`

	// ISP 網路服務提供商
	ISP string `json:"isp,omitempty"`

	// Organization 組織名稱
	Organization string `json:"organization,omitempty"`
}

// defaultGeoIPProvider 預設的 GeoIP 提供者（內部使用）
var defaultGeoIPProvider GeoIPProvider

// SetGeoIPProvider 設定全域的 GeoIP 服務提供者。
//
// 使用此函式可整合外部 GeoIP 服務。設定後，GetLocationByIP
// 將會使用該提供者進行查詢。
//
// 範例（整合 MaxMind GeoIP2）：
//
//	type MaxMindProvider struct {
//	    reader *geoip2.Reader
//	}
//
//	func (p *MaxMindProvider) Lookup(ip string) (*GeoLocation, error) {
//	    // 實作查詢邏輯...
//	}
//
//	provider := &MaxMindProvider{reader: reader}
//	net.SetGeoIPProvider(provider)
func SetGeoIPProvider(provider GeoIPProvider) {
	defaultGeoIPProvider = provider
}

// GetGeoLocation 取得指定 IP 的詳細地理位置資訊。
//
// 需先透過 SetGeoIPProvider 設定 GeoIP 服務提供者，
// 否則回傳錯誤。
//
// 範例：
//
//	loc, err := GetGeoLocation("8.8.8.8")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("國家: %s, 城市: %s\n", loc.Country, loc.City)
func GetGeoLocation(ip string) (*GeoLocation, error) {
	if defaultGeoIPProvider == nil {
		return nil, fmt.Errorf("未設定 GeoIP 服務提供者，請先呼叫 SetGeoIPProvider")
	}

	parsed := net.ParseIP(strings.TrimSpace(ip))
	if parsed == nil {
		return nil, fmt.Errorf("無效的 IP 位址: %s", ip)
	}

	return defaultGeoIPProvider.Lookup(ip)
}

// GetLocationByIP 根據 IP 位址回傳簡化的地理位置描述。
//
// 此函式會依序判斷：
//  1. 迴環位址（127.0.0.1, ::1）→ "本地"
//  2. 私有 IP → "內部網絡"
//  3. 若已設定 GeoIP 提供者，查詢詳細位置
//  4. 預設回傳 "未知位置"
//
// 範例：
//
//	GetLocationByIP("127.0.0.1")     // "本地"
//	GetLocationByIP("192.168.1.1")   // "內部網絡"
//	GetLocationByIP("8.8.8.8")       // "美國" 或 "未知位置"（視 GeoIP 設定）
func GetLocationByIP(ip string) string {
	if ip == "" {
		return ""
	}

	ipAddr := net.ParseIP(strings.TrimSpace(ip))
	if ipAddr == nil {
		return ""
	}

	// 判斷迴環位址
	if ip == "127.0.0.1" || ip == "::1" {
		return "本地"
	}

	// 判斷私有 IP
	if isPrivateIP(ipAddr) {
		return "內部網絡"
	}

	// 若有設定 GeoIP 提供者，嘗試查詢
	if defaultGeoIPProvider != nil {
		loc, err := defaultGeoIPProvider.Lookup(ip)
		if err == nil && loc != nil {
			// 組合位置描述
			parts := make([]string, 0, 3)
			if loc.Country != "" {
				parts = append(parts, loc.Country)
			}
			if loc.Region != "" {
				parts = append(parts, loc.Region)
			}
			if loc.City != "" {
				parts = append(parts, loc.City)
			}
			if len(parts) > 0 {
				return strings.Join(parts, " ")
			}
		}
	}

	return "未知位置"
}

// =============================================================================
// 私有 IP 判斷
// =============================================================================

// privateIPv4Blocks IPv4 私有與保留網段列表。
//
// 包含以下網段：
//   - 10.0.0.0/8       (RFC1918 私有網段)
//   - 172.16.0.0/12    (RFC1918 私有網段)
//   - 192.168.0.0/16   (RFC1918 私有網段)
//   - 100.64.0.0/10    (RFC6598 CGNAT 電信商共享 NAT)
//   - 169.254.0.0/16   (RFC3927 link-local)
//   - 127.0.0.0/8      (迴環位址)
//   - 192.0.0.0/24     (RFC6890 IETF 協議分配)
//   - 192.0.2.0/24     (RFC5737 TEST-NET-1 文檔範例)
//   - 198.51.100.0/24  (RFC5737 TEST-NET-2 文檔範例)
//   - 203.0.113.0/24   (RFC5737 TEST-NET-3 文檔範例)
//   - 198.18.0.0/15    (RFC2544 網路設備基準測試)
var privateIPv4Blocks = []string{
	"10.0.0.0/8",      // RFC1918 私有網段
	"172.16.0.0/12",   // RFC1918 私有網段
	"192.168.0.0/16",  // RFC1918 私有網段
	"100.64.0.0/10",   // RFC6598 CGNAT（電信商共享 NAT）
	"169.254.0.0/16",  // RFC3927 link-local
	"127.0.0.0/8",     // 迴環位址
	"192.0.0.0/24",    // RFC6890 IETF 協議分配
	"192.0.2.0/24",    // RFC5737 TEST-NET-1（文檔範例）
	"198.51.100.0/24", // RFC5737 TEST-NET-2（文檔範例）
	"203.0.113.0/24",  // RFC5737 TEST-NET-3（文檔範例）
	"198.18.0.0/15",   // RFC2544 網路設備基準測試
}

// privateIPv6Blocks IPv6 私有與保留網段列表。
//
// 包含以下網段：
//   - fc00::/7  (RFC4193 ULA 唯一本地位址)
//   - fe80::/10 (link-local)
//   - ::1/128   (迴環位址)
var privateIPv6Blocks = []string{
	"fc00::/7",  // RFC4193 ULA（唯一本地位址）
	"fe80::/10", // link-local
	"::1/128",   // 迴環位址
}

// isPrivateIP 判斷 IP 是否為私有或保留位址。
//
// 支援 IPv4 與 IPv6，會檢查是否落在 privateIPv4Blocks 或
// privateIPv6Blocks 定義的網段內。
func isPrivateIP(ip net.IP) bool {
	// IPv4 檢查
	if ip4 := ip.To4(); ip4 != nil {
		for _, block := range privateIPv4Blocks {
			_, ipnet, err := net.ParseCIDR(block)
			if err != nil {
				continue
			}
			if ipnet.Contains(ip4) {
				return true
			}
		}
		return false
	}

	// IPv6 檢查
	for _, block := range privateIPv6Blocks {
		_, ipnet, err := net.ParseCIDR(block)
		if err != nil {
			continue
		}
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}

// =============================================================================
// 客戶端 IP 偵測
// =============================================================================

// GetClientIP 從 HTTP headers map 中取得客戶端真實 IP。
//
// 此函式會依序檢查以下 header：
//  1. X-Forwarded-For：格式為 "client, proxy1, proxy2"，取第一個有效 IP
//  2. X-Real-IP：通常由 Nginx/ALB/Ingress 設定
//  3. 若都無法取得，回傳 "127.0.0.1"
//
// 注意：header key 不區分大小寫。
//
// 範例：
//
//	headers := map[string][]string{
//	    "X-Forwarded-For": {"203.0.113.195, 70.41.3.18"},
//	}
//	GetClientIP(headers) // "203.0.113.195"
func GetClientIP(headers map[string][]string) string {
	if headers == nil {
		return "127.0.0.1"
	}

	// 將 key 統一轉成小寫，避免大小寫敏感造成取不到
	lower := make(map[string][]string, len(headers))
	for k, v := range headers {
		lower[strings.ToLower(k)] = v
	}

	// 1) X-Forwarded-For：格式為 "client, proxy1, proxy2"
	if xff, ok := lower["x-forwarded-for"]; ok && len(xff) > 0 {
		parts := strings.Split(xff[0], ",")
		for _, p := range parts {
			candidate := strings.TrimSpace(p)
			if candidate == "" {
				continue
			}
			if ip := net.ParseIP(candidate); ip != nil {
				return candidate
			}
		}
	}

	// 2) X-Real-IP：通常由 Nginx/ALB/Ingress 設定
	if xri, ok := lower["x-real-ip"]; ok && len(xri) > 0 {
		candidate := strings.TrimSpace(xri[0])
		if candidate != "" {
			if ip := net.ParseIP(candidate); ip != nil {
				return candidate
			}
		}
	}

	// 3) fallback：沒有任何 header 或格式不正確
	return "127.0.0.1"
}

// =============================================================================
// 本機 IP 取得
// =============================================================================

// GetLocalIPs 取得本機所有非迴環的 IPv4 位址。
//
// 回傳以逗號分隔的 IP 列表字串，適合用於 log 或查詢顯示。
// 若無法取得任何 IP，回傳空字串。
//
// 範例：
//
//	GetLocalIPs() // "192.168.1.100,10.0.0.5"
func GetLocalIPs() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	ips := make([]string, 0, 4)
	for _, addr := range addrs {
		ipnet, ok := addr.(*net.IPNet)
		if !ok || ipnet.IP == nil {
			continue
		}
		// 排除迴環位址
		if ipnet.IP.IsLoopback() {
			continue
		}
		// 只取 IPv4
		ip4 := ipnet.IP.To4()
		if ip4 == nil {
			continue
		}
		ips = append(ips, ip4.String())
	}

	return strings.Join(ips, ",")
}
