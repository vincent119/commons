package s3

import "testing"

func TestBuildS3Prefix(t *testing.T) {
	tests := []struct {
		name         string
		bucketPrefix string
		mediaPrefix  string
		want         string
	}{
		{"normal", "bucket", "media", "bucket/media/"},
		{"with_slashes", "bucket/", "/media/", "bucket/media/"},
		{"empty_media", "bucket", "", "bucket//"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildS3Prefix(tt.bucketPrefix, tt.mediaPrefix); got != tt.want {
				t.Errorf("BuildS3Prefix(%q, %q) = %q, want %q", tt.bucketPrefix, tt.mediaPrefix, got, tt.want)
			}
		})
	}
}

func TestBuildPrefix(t *testing.T) {
	tests := []struct {
		name  string
		parts []string
		want  string
	}{
		{"empty", nil, "/"},
		{"one_part", []string{"foo"}, "foo/"},
		{"multiple_parts", []string{"foo", "bar"}, "foo/bar/"},
		{"with_slashes", []string{"/foo/", "/bar/"}, "foo/bar/"},
		{"empty_parts", []string{"foo", "", "bar"}, "foo/bar/"},
		{"all_empty", []string{"", ""}, "/"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildPrefix(tt.parts...); got != tt.want {
				t.Errorf("BuildPrefix(%v) = %q, want %q", tt.parts, got, tt.want)
			}
		})
	}
}
