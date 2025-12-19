// Package responsex 提供 HTTP API 的標準回應結構定義。
//
// 此套件定義了常用的 API 回應結構，適用於 Swagger/OpenAPI 文檔生成。
//
// # 錯誤回應
//
// 標準的 API 錯誤回應結構：
//
//	errResp := responsex.Error{
//	    Code:    401,
//	    Message: "unauthorized",
//	}
//
// # 健康檢查
//
// 健康檢查端點回應：
//
//	health := responsex.Health{
//	    Status: "ok",
//	}
//
// # Swagger 標籤
//
// 所有結構都包含 json 與 example 標籤，方便 Swagger 文檔生成：
//
//	type Error struct {
//	    Code    int    `json:"code" example:"401"`
//	    Message string `json:"message" example:"unauthorized"`
//	}
package resp
