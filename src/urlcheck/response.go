package main

// Response handling (HTTP)

import (
	"urlcheck/utils"

	"encoding/json"
	"net/http"
)

// Empty data type for responses that don't contain data.
type EmptyData struct{}

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
