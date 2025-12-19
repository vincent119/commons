package responsex

// Error represents a standard API error response
type Error struct {
	Code    int    `json:"code" example:"401"`
	Message string `json:"message" example:"unauthorized"`
}
