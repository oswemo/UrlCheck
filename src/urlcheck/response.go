package main

// Response handling (HTTP)

import (
	"urlcheck/utils"

	"encoding/json"
	"net/http"
)

// Empty data type for responses that don't contain data.
type EmptyData struct{}

// String data type for simple responses of a single string.
type StringData struct {
	Message string `json:"message"`
}

// HttpStatus is a simple wrapper containing the HTTP status code and an optional
// message.
type HttpStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// APIError is used to wrap any errors in the Data section of the APIResponse
type APIError struct {
	Error string `json:"error"`
}

// APIResponse is used to wrap all API responses in a standard layout.
type APIResponse struct {
	// HTTP status, notes HTTP status code and optional message.
	Status HttpStatus `json:"status"`

	// Response object.
	Data interface{} `json:"data"`
}

// BadRequest generates an APIResponse with embedded error messaging and a Bad Requet status
func (a *APIResponse) BadRequest(err error) {
	a.Status = HttpStatus{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)}
	a.Data = APIError{Error: err.Error()}
}

// InternalError generates an APIResponse with embedded error messaging and a Internal Error status
func (a *APIResponse) InternalError(err error) {
	a.Status = HttpStatus{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)}
	a.Data = APIError{Error: err.Error()}
}

// Success generates an APIResponse with embedded error messaging and a Success (OK) status
func (a *APIResponse) Success(data interface{}) {
	a.Status = HttpStatus{Code: http.StatusOK, Message: http.StatusText(http.StatusOK)}
	a.Data = data
}

// http_respond marhals an APIResponse object and writes the output to the
// http requestor.  If marshaling fails, a simple 500 error is returned.
func http_respond(response APIResponse, writer http.ResponseWriter) {
	// Rather than sending a "null" object back in JSON which is sort of ugly,
	// we'll send back an empty object any time the Data object is null.
	if response.Data == nil {
		response.Data = EmptyData{}
	}

	// Marshal the response to be sent back to the user.  If marshaling fails for
	// some reason, we log an error and return a simple 500 error.
	json_result, err := json.Marshal(response)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to process response")
		writer.WriteHeader(500)
		writer.Write([]byte("An error occurred responding to your request.\n"))
		return
	}

	// Return the marshaled response.
	writer.WriteHeader(response.Status.Code)
	writer.Write(json_result)
}
