package modelx



// ErrorResponse represents a standard error response
type ErrorResponse struct {
    Code    int    `json:"code" example:"401"`
    Message string `json:"message" example:"unauthorized"`
}

// ResponseHealthCheck represents health check response
type ResponseHealthCheck struct {
    Status string `json:"status" example:"ok"`
}