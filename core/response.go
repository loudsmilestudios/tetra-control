package core

import "encoding/json"

// Response is the base struct for all responses by TetraControl
type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewErrorResponse is a quick helper function for generating error responses
func NewErrorResponse(message string) []byte {
	data, _ := json.Marshal(Response{
		Success: false,
		Message: message,
	})
	return data
}

// NewUnknownResponse is a quick helper function for generating generic error responses
func NewUnknownResponse() []byte {
	data, _ := json.Marshal(Response{
		Success: false,
		Message: "An unknown error occured.",
	})
	return data
}
