package cryptox

import "testing"

func TestMD5Hash(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"hello", "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"test", "test", "098f6bcd4621d373cade4e832627b4f6"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5Hash(tt.in); got != tt.want {
				t.Errorf("MD5Hash(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestSHA256Hash(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256Hash(tt.in); got != tt.want {
				t.Errorf("SHA256Hash(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
