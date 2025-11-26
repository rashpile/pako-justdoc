package model

// SuccessResponse for POST operations
type SuccessResponse struct {
	Status   string `json:"status"`   // "created" or "updated"
	Channel  string `json:"channel"`
	Document string `json:"document"`
}

// ErrorResponse for all error cases
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}