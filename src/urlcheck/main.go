package main

import (
	"os"
	"fmt"
	"flag"
	"net/url"
	"net/http"
	"encoding/json"

	"urlcheck/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	// Default port to attach HTTP service to.
	DEFAULT_PORT = "8001"
)

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
	Data interface{}  `json:"data"`
}

// checkUrl attempts to validate whether a hostname and query path constitute a known
// malware site.  It is expected that the query path is a URI encoded string.
// An HTTP response is returned vai the writer parameter.
//
func checkUrl(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	hostname := vars["hostname"]
	querypath := vars["querypath"]

	var response APIResponse

	// We had to enable "UseEncodedPath" on Gorilla MUX to fix an issue where path
	// matching failed if the querypath started with an encoded/escaped character (ie, %2F)
	// Because of this, we need to unescape/decode the path here.  We have error handling here, though
	// Gorilla MUX should return a 400 response on it's own if encoding is bad.
	decodedPath, err := url.PathUnescape(querypath) ; if err != nil {
		utils.LogError(utils.LogFields{ }, err, "Failed to decode query path")

		response = APIResponse{
			Status: HttpStatus{ Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest) },
			Data  : APIError{ Error: err.Error() },
		}

	} else
	// Successful decoding of the query path.
	// Perform service call to check URL.
	{
		utils.LogInfo(utils.LogFields{ "hostname": hostname, "path": decodedPath }, "Checking URL status")
	}

	// Marshal the response to be sent back to the user.  If marshaling fails for
	// some reason, we log an error and return a simple 500 error.
	json_result, err := json.Marshal(response) ; if err != nil {
		utils.LogError(utils.LogFields{ }, err, "Failed to process response")
		writer.WriteHeader(500)
		writer.Write([]byte("An error occurred responding to your request.\n"))
		return
	}

	// Return the marshaled response.
	writer.WriteHeader(response.Status.Code)
	writer.Write(json_result)
}

// Main function
func main() {
	port := flag.String("port", DEFAULT_PORT, "Listening port")
	debug := flag.Bool("debug", false, "Debug output")
	flag.Parse()

	if *debug {
        utils.SetDebug()
    }

	router := mux.NewRouter().StrictSlash(true)
	router.UseEncodedPath()

    // Query URL info.  query path is optional.
	router.HandleFunc("/urlinfo/1/{hostname}/{querypath:.*}", checkUrl).Methods("GET")

	utils.LogInfo(utils.LogFields{ "port": *port}, "Starting HTTP service")
	http.ListenAndServe(fmt.Sprintf(":%s", *port), handlers.LoggingHandler(os.Stdout, router))
}
