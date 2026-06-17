package dto

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Results any    `json:"results,omitempty"`
	Error   string `json:"error,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}
