package plethora_api

// Success response
// swagger:response ok
type swaggerOKResponse struct{}

// Error response
// swagger:response err
type swaggerErrorResponse struct{}

// Error response with message
// swagger:response errMsg
type swaggerErrorMessageResponse struct {
	Message string `json:"message"`
}
