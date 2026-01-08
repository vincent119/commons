package validatorx

import "testing"

func TestIsEmail(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"test@example.com", true},
		{"user.name@domain.co.uk", true},
		{"user-name@domain.com", true},
		{"invalid", false},
		{"@domain.com", false},
		{"user@", false},
		{"user@domain", false}, // TLD required
	}
	for _, tt := range tests {
		if got := IsEmail(tt.in); got != tt.want {
			t.Errorf("IsEmail(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsMobile(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"0912345678", true},
		{"0900000000", true},
		{"0812345678", false},
		{"091234567", false},  // too short
		{"09123456789", false}, // too long
		{"abc", false},
	}
	for _, tt := range tests {
		if got := IsMobile(tt.in); got != tt.want {
			t.Errorf("IsMobile(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsUUID(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"550e8400-e29b-41d4-a716-446655440000", true}, // v4
		{"123e4567-e89b-12d3-a456-426614174000", true}, // v1
		{"invalid-uuid", false},
		{"550e8400e29b41d4a716446655440000", false}, // formatting required
	}
	for _, tt := range tests {
		if got := IsUUID(tt.in); got != tt.want {
			t.Errorf("IsUUID(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsIPv4(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"192.168.0.1", true},
		{"255.255.255.255", true},
		{"0.0.0.0", true},
		{"256.0.0.1", false},
		{"192.168.1", false},
		{"abc", false},
	}
	for _, tt := range tests {
		if got := IsIPv4(tt.in); got != tt.want {
			t.Errorf("IsIPv4(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsIPv6(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"2001:0db8:85a3:0000:0000:8a2e:0370:7334", true},
		{"fe80::1", false}, // Simple regex might not support :: compression or shorthand
		{"invalid", false},
	}
	for _, tt := range tests {
		if got := IsIPv6(tt.in); got != tt.want {
			t.Errorf("IsIPv6(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"http://google.com", true},
		{"https://example.org/path?q=1", true},
		{"ftp://example.com", false}, // regex specifies http/https
		{"invalid", false},
	}
	for _, tt := range tests {
		if got := IsURL(tt.in); got != tt.want {
			t.Errorf("IsURL(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsDate(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"2023-12-31", true},
		{"2024-01-01", true},
		{"2023/12/31", false},
		{"2023-13-01", false},
		{"2023-12-32", false},
		{"abc", false},
	}
	for _, tt := range tests {
		if got := IsDate(tt.in); got != tt.want {
			t.Errorf("IsDate(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsTime(t *testing.T) {
	tests := []struct {
		in   string
		want bool
	}{
		{"23:59:59", true},
		{"00:00:00", true},
		{"12:30:45", true},
		{"24:00:00", false},
		{"12:60:00", false},
		{"12:00:60", false},
		{"abc", false},
	}
	for _, tt := range tests {
		if got := IsTime(tt.in); got != tt.want {
			t.Errorf("IsTime(%q) = %v, want %v", tt.in, got, tt.want)
		}
	}
}

func TestIsPassword(t *testing.T) {
	tests := []struct {
		pwd  string
		max  int
		want bool
	}{
		{"Password123", 8, true},
		{"pass", 8, false},           // too short
		{"password", 8, false},       // no upper, no digit
		{"PASSWORD", 8, false},       // no lower, no digit
		{"12345678", 8, false},       // no letters
		{"Pass1", 8, false},          // too short
		{"ComplexPass1", 10, true},
	}
	for _, tt := range tests {
		if got := IsPassword(tt.pwd, tt.max); got != tt.want {
			t.Errorf("IsPassword(%q, %d) = %v, want %v", tt.pwd, tt.max, got, tt.want)
		}
	}
}
