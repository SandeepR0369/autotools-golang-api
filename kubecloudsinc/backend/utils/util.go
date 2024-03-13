package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

// Modified ErrorResponse to match the desired structure
type ErrorResponse struct {
	Metadata struct {
		ID                string `json:"id"`
		Name              string `json:"name"`
		Status            int    `json:"status"`
		Method            string `json:"method"`
		AdditionalDetails struct {
			Description   string `json:"description"`
			StatusCode    int    `json:"statusCode"`
			Code          string `json:"code"`
			EsrxRequestID string `json:"esrxRequestId"`
			ErrorLocation string `json:"errorLocation"`
		} `json:"AdditionalDetails"`
	} `json:"metadata"`
}

// SendErrorResponse sends a JSON-encoded error response with the new structure
func SendErrorResponse(w http.ResponseWriter, r *http.Request, statusCode int, err error, errorID, errorCode, errorLocation string) {
	resp := ErrorResponse{}
	resp.Metadata.ID = errorID
	resp.Metadata.Name = http.StatusText(statusCode)
	resp.Metadata.Status = statusCode
	resp.Metadata.Method = r.Method
	resp.Metadata.AdditionalDetails.Description = err.Error()
	resp.Metadata.AdditionalDetails.StatusCode = statusCode
	resp.Metadata.AdditionalDetails.Code = errorCode
	resp.Metadata.AdditionalDetails.EsrxRequestID = "someUniqueRequestID" // This should be dynamically generated or passed as an argument
	resp.Metadata.AdditionalDetails.ErrorLocation = errorLocation

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		log.Printf("Error encoding the error response: %v", encodeErr)
	}
}
