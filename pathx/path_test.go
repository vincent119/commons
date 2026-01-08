package pathx

import "testing"

func TestNormalizePathSeparator(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{"empty", "", ""},
		{"unix", "path/to/file", "path/to/file"},
		{"windows", "path\\to\\file", "path/to/file"},
		{"mixed", "path/to\\file", "path/to/file"},
		{"double_backslash", "path\\\\to\\\\file", "path//to//file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NormalizePathSeparator(tt.path); got != tt.want {
				t.Errorf("NormalizePathSeparator(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
