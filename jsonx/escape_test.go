package jsonx

import "testing"

func TestEscapeJSON(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"normal", "hello world", "hello world"},
		{"quotes", `hello "world"`, `hello \"world\"`},
		{"backslashes", `C:\Windows\System32`, `C:\\Windows\\System32`},
		{"newlines", "line1\nline2", "line1\\nline2"},
		{"tabs", "col1\tcol2", "col1\\tcol2"},
		{"carriage_return", "row1\rrow2", "row1\\rrow2"},
		{"mixed", "a\tb\nc\"d\\e", "a\\tb\\nc\\\"d\\\\e"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EscapeJSON(tt.in); got != tt.want {
				t.Errorf("EscapeJSON(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
