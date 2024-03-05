package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorResponse represents the structure of the error response
type ErrorResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Method      string `json:"method"`
	Status      int    `json:"status"`
	StatusCode  int    `json:"statusCode"`
	Details     struct {
		RequestID     string `json:"esrxRequestId"`
		ErrorLocation string `json:"errorLocation"`
	} `json:"details"`
}

// SendErrorResponse sends a JSON-encoded error response
func SendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err error, errorID, errorCode, errorLocation string) {
	resp := ErrorResponse{
		ID:          errorID,
		Name:        http.StatusText(statusCode),
		Code:        errorCode,
		Description: err.Error(),
		Method:      r.Method,
		Status:      statusCode,
		StatusCode:  statusCode,
		Details: struct {
			RequestID     string `json:"esrxRequestId"`
			ErrorLocation string `json:"errorLocation"`
		}{
			RequestID:     "someUniqueRequestID", // This should be a unique ID for the request
			ErrorLocation: errorLocation,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		log.Printf("Error encoding the error response: %v", encodeErr)
	}
}
