package stringx


import "strings"

// BuildS3Prefix 建立 S3 路徑前綴
func BuildS3Prefix(bucketPrefix, mediaPrefix string) string {

	bucketPrefix = strings.TrimSuffix(bucketPrefix, "/")
	mediaPrefix = strings.Trim(mediaPrefix, "/")
	return bucketPrefix + "/" + mediaPrefix + "/"
}