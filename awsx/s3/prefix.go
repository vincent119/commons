package s3

import "strings"

// BuildS3Prefix 建立 S3 路徑前綴
func BuildS3Prefix(bucketPrefix, mediaPrefix string) string {

	bucketPrefix = strings.TrimSuffix(bucketPrefix, "/")
	mediaPrefix = strings.Trim(mediaPrefix, "/")
	return bucketPrefix + "/" + mediaPrefix + "/"
}

// BuildPrefix 建立路徑前綴
func BuildPrefix(parts ...string) string {
	var cleaned []string
	for _, p := range parts {
		p = strings.Trim(p, "/")
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}
	return strings.Join(cleaned, "/") + "/"
}
