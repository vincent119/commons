package resp

// Health represents health check response
type Health struct {
	Status string `json:"status" example:"ok"`
}
