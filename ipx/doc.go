// Package ipx 提供 IP 位址相關的通用工具函式。
//
// # IP 驗證
//
// 驗證 IP 位址格式與類型：
//
//	ipx.IsValidIP("192.168.1.1")  // true
//	ipx.IsIPv4("192.168.1.1")     // true
//	ipx.IsIPv6("2001:db8::1")     // true
//	ipx.IsPublicIP("8.8.8.8")     // true
//
// # IP 轉換
//
// IPv4 與 uint32 互轉：
//
//	n, _ := ipx.IPv4ToUint32("192.168.1.1")  // 3232235777
//	ip := ipx.Uint32ToIPv4(3232235777)       // "192.168.1.1"
//
// IPv6 展開：
//
//	expanded, _ := ipx.ExpandIPv6("::1")
//	// "0000:0000:0000:0000:0000:0000:0000:0001"
//
// # 網段工具
//
// 判斷 IP 是否在 CIDR 網段內：
//
//	inCIDR, _ := ipx.IsIPInCIDR("192.168.1.100", "192.168.1.0/24") // true
//
// 取得網段詳細資訊：
//
//	info, _ := ipx.GetNetworkInfo("192.168.1.0/24")
//	// info.Network = "192.168.1.0"
//	// info.Broadcast = "192.168.1.255"
//	// info.TotalHosts = 254
//
// # 地理位置
//
// 簡化地理位置判斷：
//
//	ipx.GetLocationByIP("127.0.0.1")     // "本地"
//	ipx.GetLocationByIP("192.168.1.1")   // "內部網絡"
//
// 整合 GeoIP 服務（需實作 GeoIPProvider 介面）：
//
//	ipx.SetGeoIPProvider(myProvider)
//	loc, _ := ipx.GetGeoLocation("8.8.8.8")
//
// # 客戶端 IP 偵測
//
// 從 HTTP headers 取得真實客戶端 IP：
//
//	clientIP := ipx.GetClientIP(headers)
//
// # 私有網段支援
//
// 支援以下 RFC 定義的私有與保留網段：
//   - RFC1918: 10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
//   - RFC6598: 100.64.0.0/10 (CGNAT)
//   - RFC5737: 192.0.2.0/24, 198.51.100.0/24, 203.0.113.0/24 (TEST-NET)
//   - RFC2544: 198.18.0.0/15 (基準測試)
//   - RFC3927: 169.254.0.0/16 (link-local)
//   - RFC4193: fc00::/7 (IPv6 ULA)
package ipx
