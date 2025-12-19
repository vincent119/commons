package ipx

import (
	"net"
	"testing"
)

// =============================================================================
// IP 驗證工具測試
// =============================================================================

func TestIsValidIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"有效 IPv4", "192.168.1.1", true},
		{"有效 IPv4 - 邊界", "0.0.0.0", true},
		{"有效 IPv4 - 最大值", "255.255.255.255", true},
		{"有效 IPv6", "2001:db8::1", true},
		{"有效 IPv6 - 迴環", "::1", true},
		{"有效 IPv6 - 完整格式", "2001:0db8:0000:0000:0000:0000:0000:0001", true},
		{"無效 - 空字串", "", false},
		{"無效 - 非 IP", "invalid", false},
		{"無效 - 超出範圍", "256.1.1.1", false},
		{"有效 - 含空白", "  192.168.1.1  ", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidIP(tt.ip)
			if result != tt.expected {
				t.Errorf("IsValidIP(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"有效 IPv4", "192.168.1.1", true},
		{"IPv6 - 非 IPv4", "::1", false},
		{"IPv6 - 非 IPv4", "2001:db8::1", false},
		{"無效 IP", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIPv4(tt.ip)
			if result != tt.expected {
				t.Errorf("IsIPv4(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestIsIPv6(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"有效 IPv6", "2001:db8::1", true},
		{"有效 IPv6 - 迴環", "::1", true},
		{"IPv4 - 非 IPv6", "192.168.1.1", false},
		{"無效 IP", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsIPv6(tt.ip)
			if result != tt.expected {
				t.Errorf("IsIPv6(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestIsPublicIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		{"公網 IP", "8.8.8.8", true},
		{"公網 IP - Google DNS", "8.8.4.4", true},
		{"私有 IP - 10.x.x.x", "10.0.0.1", false},
		{"私有 IP - 172.16.x.x", "172.16.0.1", false},
		{"私有 IP - 192.168.x.x", "192.168.1.1", false},
		{"私有 IP - CGNAT", "100.64.0.1", false},
		{"私有 IP - TEST-NET-1", "192.0.2.1", false},
		{"私有 IP - TEST-NET-2", "198.51.100.1", false},
		{"私有 IP - TEST-NET-3", "203.0.113.1", false},
		{"私有 IP - 基準測試", "198.18.0.1", false},
		{"迴環位址", "127.0.0.1", false},
		{"無效 IP", "invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPublicIP(tt.ip)
			if result != tt.expected {
				t.Errorf("IsPublicIP(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// IP 轉換工具測試
// =============================================================================

func TestIPv4ToUint32(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected uint32
		wantErr  bool
	}{
		{"標準 IP", "192.168.1.1", 3232235777, false},
		{"最小值", "0.0.0.0", 0, false},
		{"最大值", "255.255.255.255", 4294967295, false},
		{"10.0.0.1", "10.0.0.1", 167772161, false},
		{"無效 IP", "invalid", 0, true},
		{"IPv6 - 非 IPv4", "::1", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IPv4ToUint32(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("IPv4ToUint32(%q) error = %v, wantErr %v", tt.ip, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("IPv4ToUint32(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestUint32ToIPv4(t *testing.T) {
	tests := []struct {
		name     string
		n        uint32
		expected string
	}{
		{"標準 IP", 3232235777, "192.168.1.1"},
		{"最小值", 0, "0.0.0.0"},
		{"最大值", 4294967295, "255.255.255.255"},
		{"10.0.0.1", 167772161, "10.0.0.1"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Uint32ToIPv4(tt.n)
			if result != tt.expected {
				t.Errorf("Uint32ToIPv4(%v) = %v, want %v", tt.n, result, tt.expected)
			}
		})
	}
}

func TestExpandIPv6(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
		wantErr  bool
	}{
		{"迴環位址", "::1", "0000:0000:0000:0000:0000:0000:0000:0001", false},
		{"標準縮寫", "2001:db8::1", "2001:0db8:0000:0000:0000:0000:0000:0001", false},
		{"完整格式", "2001:0db8:0000:0000:0000:0000:0000:0001", "2001:0db8:0000:0000:0000:0000:0000:0001", false},
		{"IPv4 - 非 IPv6", "192.168.1.1", "", true},
		{"無效 IP", "invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpandIPv6(tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExpandIPv6(%q) error = %v, wantErr %v", tt.ip, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("ExpandIPv6(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 網段相關工具測試
// =============================================================================

func TestIsIPInCIDR(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		cidr     string
		expected bool
		wantErr  bool
	}{
		{"在網段內", "192.168.1.100", "192.168.1.0/24", true, false},
		{"在網段內 - 邊界", "192.168.1.0", "192.168.1.0/24", true, false},
		{"在網段內 - 廣播", "192.168.1.255", "192.168.1.0/24", true, false},
		{"不在網段內", "10.0.0.1", "192.168.1.0/24", false, false},
		{"無效 IP", "invalid", "192.168.1.0/24", false, true},
		{"無效 CIDR", "192.168.1.1", "invalid", false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := IsIPInCIDR(tt.ip, tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsIPInCIDR(%q, %q) error = %v, wantErr %v", tt.ip, tt.cidr, err, tt.wantErr)
				return
			}
			if result != tt.expected {
				t.Errorf("IsIPInCIDR(%q, %q) = %v, want %v", tt.ip, tt.cidr, result, tt.expected)
			}
		})
	}
}

func TestGetNetworkInfo(t *testing.T) {
	tests := []struct {
		name           string
		cidr           string
		wantNetwork    string
		wantBroadcast  string
		wantTotalHosts uint64
		wantPrefixLen  int
		wantFirstHost  string
		wantLastHost   string
		wantErr        bool
	}{
		{
			name:           "/24 網段",
			cidr:           "192.168.1.0/24",
			wantNetwork:    "192.168.1.0",
			wantBroadcast:  "192.168.1.255",
			wantTotalHosts: 254,
			wantPrefixLen:  24,
			wantFirstHost:  "192.168.1.1",
			wantLastHost:   "192.168.1.254",
			wantErr:        false,
		},
		{
			name:           "/16 網段",
			cidr:           "10.0.0.0/16",
			wantNetwork:    "10.0.0.0",
			wantBroadcast:  "10.0.255.255",
			wantTotalHosts: 65534,
			wantPrefixLen:  16,
			wantFirstHost:  "10.0.0.1",
			wantLastHost:   "10.0.255.254",
			wantErr:        false,
		},
		{
			name:    "無效 CIDR",
			cidr:    "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := GetNetworkInfo(tt.cidr)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkInfo(%q) error = %v, wantErr %v", tt.cidr, err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}
			if info.Network != tt.wantNetwork {
				t.Errorf("Network = %v, want %v", info.Network, tt.wantNetwork)
			}
			if info.Broadcast != tt.wantBroadcast {
				t.Errorf("Broadcast = %v, want %v", info.Broadcast, tt.wantBroadcast)
			}
			if info.TotalHosts != tt.wantTotalHosts {
				t.Errorf("TotalHosts = %v, want %v", info.TotalHosts, tt.wantTotalHosts)
			}
			if info.PrefixLength != tt.wantPrefixLen {
				t.Errorf("PrefixLength = %v, want %v", info.PrefixLength, tt.wantPrefixLen)
			}
			if info.FirstHost != tt.wantFirstHost {
				t.Errorf("FirstHost = %v, want %v", info.FirstHost, tt.wantFirstHost)
			}
			if info.LastHost != tt.wantLastHost {
				t.Errorf("LastHost = %v, want %v", info.LastHost, tt.wantLastHost)
			}
		})
	}
}

// =============================================================================
// 地理位置工具測試
// =============================================================================

// mockGeoIPProvider 測試用的 GeoIP 提供者
type mockGeoIPProvider struct{}

func (m *mockGeoIPProvider) Lookup(ip string) (*GeoLocation, error) {
	return &GeoLocation{
		IP:          ip,
		Country:     "台灣",
		CountryCode: "TW",
		Region:      "台北市",
		City:        "信義區",
		Latitude:    25.0330,
		Longitude:   121.5654,
	}, nil
}

func TestGetLocationByIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected string
	}{
		{"空字串", "", ""},
		{"迴環位址 IPv4", "127.0.0.1", "本地"},
		{"迴環位址 IPv6", "::1", "本地"},
		{"私有 IP - 192.168.x.x", "192.168.1.1", "內部網絡"},
		{"私有 IP - 10.x.x.x", "10.0.0.1", "內部網絡"},
		{"私有 IP - CGNAT", "100.64.0.1", "內部網絡"},
		{"私有 IP - TEST-NET-1", "192.0.2.1", "內部網絡"},
		{"無效 IP", "invalid", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLocationByIP(tt.ip)
			if result != tt.expected {
				t.Errorf("GetLocationByIP(%q) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

func TestGetLocationByIP_WithProvider(t *testing.T) {
	// 設定 mock provider
	SetGeoIPProvider(&mockGeoIPProvider{})
	defer SetGeoIPProvider(nil) // 清理

	result := GetLocationByIP("8.8.8.8")
	expected := "台灣 台北市 信義區"
	if result != expected {
		t.Errorf("GetLocationByIP(\"8.8.8.8\") = %v, want %v", result, expected)
	}
}

func TestGetGeoLocation_NoProvider(t *testing.T) {
	// 確保無 provider
	SetGeoIPProvider(nil)

	_, err := GetGeoLocation("8.8.8.8")
	if err == nil {
		t.Error("GetGeoLocation without provider should return error")
	}
}

// =============================================================================
// 客戶端 IP 偵測測試
// =============================================================================

func TestGetClientIP(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string][]string
		expected string
	}{
		{"空 headers", nil, "127.0.0.1"},
		{"無相關 header", map[string][]string{"other": {"value"}}, "127.0.0.1"},
		{
			"X-Forwarded-For 單一 IP",
			map[string][]string{"X-Forwarded-For": {"203.0.113.195"}},
			"203.0.113.195",
		},
		{
			"X-Forwarded-For 多個 IP",
			map[string][]string{"X-Forwarded-For": {"203.0.113.195, 70.41.3.18, 150.172.238.178"}},
			"203.0.113.195",
		},
		{
			"X-Forwarded-For 小寫",
			map[string][]string{"x-forwarded-for": {"203.0.113.195"}},
			"203.0.113.195",
		},
		{
			"X-Real-IP",
			map[string][]string{"X-Real-IP": {"203.0.113.195"}},
			"203.0.113.195",
		},
		{
			"X-Forwarded-For 優先於 X-Real-IP",
			map[string][]string{
				"X-Forwarded-For": {"203.0.113.195"},
				"X-Real-IP":       {"70.41.3.18"},
			},
			"203.0.113.195",
		},
		{
			"X-Forwarded-For 無效時使用 X-Real-IP",
			map[string][]string{
				"X-Forwarded-For": {"invalid"},
				"X-Real-IP":       {"70.41.3.18"},
			},
			"70.41.3.18",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetClientIP(tt.headers)
			if result != tt.expected {
				t.Errorf("GetClientIP() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// =============================================================================
// 本機 IP 取得測試
// =============================================================================

func TestGetLocalIPs(t *testing.T) {
	result := GetLocalIPs()
	// 只驗證函式不會 panic，實際 IP 因環境而異
	t.Logf("GetLocalIPs() = %v", result)
}

// =============================================================================
// 私有 IP 判斷測試
// =============================================================================

func TestIsPrivateIP(t *testing.T) {
	tests := []struct {
		name     string
		ip       string
		expected bool
	}{
		// RFC1918 私有網段
		{"私有 - 10.0.0.0/8", "10.0.0.1", true},
		{"私有 - 10.255.255.255", "10.255.255.255", true},
		{"私有 - 172.16.0.0/12", "172.16.0.1", true},
		{"私有 - 172.31.255.255", "172.31.255.255", true},
		{"非私有 - 172.32.0.0", "172.32.0.0", false},
		{"私有 - 192.168.0.0/16", "192.168.0.1", true},

		// RFC6598 CGNAT
		{"私有 - CGNAT 起始", "100.64.0.0", true},
		{"私有 - CGNAT 中間", "100.100.100.100", true},
		{"私有 - CGNAT 結束", "100.127.255.255", true},
		{"非私有 - CGNAT 之前", "100.63.255.255", false},
		{"非私有 - CGNAT 之後", "100.128.0.0", false},

		// RFC5737 TEST-NET
		{"私有 - TEST-NET-1", "192.0.2.1", true},
		{"私有 - TEST-NET-2", "198.51.100.1", true},
		{"私有 - TEST-NET-3", "203.0.113.1", true},

		// RFC2544 基準測試
		{"私有 - 基準測試起始", "198.18.0.1", true},
		{"私有 - 基準測試結束", "198.19.255.255", true},

		// 其他保留
		{"私有 - link-local", "169.254.1.1", true},
		{"私有 - loopback", "127.0.0.1", true},
		{"私有 - IETF 協議", "192.0.0.1", true},

		// 公網 IP
		{"公網 IP - Google DNS", "8.8.8.8", false},
		{"公網 IP - Cloudflare", "1.1.1.1", false},

		// IPv6
		{"IPv6 - ULA", "fc00::1", true},
		{"IPv6 - link-local", "fe80::1", true},
		{"IPv6 - loopback", "::1", true},
		{"IPv6 - 公網", "2001:db8::1", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed := net.ParseIP(tt.ip)
			if parsed == nil {
				t.Fatalf("Failed to parse IP: %s", tt.ip)
			}
			result := isPrivateIP(parsed)
			if result != tt.expected {
				t.Errorf("isPrivateIP(%s) = %v, want %v", tt.ip, result, tt.expected)
			}
		})
	}
}

// =============================================================================
// IPv4/Uint32 轉換往返測試
// =============================================================================

func TestIPv4Uint32RoundTrip(t *testing.T) {
	testIPs := []string{
		"0.0.0.0",
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"255.255.255.255",
		"8.8.8.8",
		"1.2.3.4",
	}

	for _, ip := range testIPs {
		t.Run(ip, func(t *testing.T) {
			n, err := IPv4ToUint32(ip)
			if err != nil {
				t.Fatalf("IPv4ToUint32(%s) error: %v", ip, err)
			}
			result := Uint32ToIPv4(n)
			if result != ip {
				t.Errorf("RoundTrip: %s -> %d -> %s", ip, n, result)
			}
		})
	}
}
