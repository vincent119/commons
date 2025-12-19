// Package aws 提供 AWS 服務相關的工具函式。
//
// 目前支援：
//   - S3 路徑前綴建構
//
// # S3 路徑工具
//
// BuildS3Prefix 建構 S3 儲存路徑前綴：
//
//	prefix := s3.BuildS3Prefix("bucket/prefix", "media/images")
//	// prefix = "bucket/prefix/media/images/"
//
// BuildPrefix 通用路徑前綴建構（支援多段）：
//
//	prefix := s3.BuildPrefix("uploads", "2025", "12")
//	// prefix = "uploads/2025/12/"
package s3
