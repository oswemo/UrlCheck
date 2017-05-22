package main

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"urlcheck/data"
	"urlcheck/services"
	"urlcheck/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// validateUrl attempts to validate that the inputs are correct.
func validateUrl(hostname string, path string) error {
	// A hostname should contain the host and port.
	parts := strings.Split(hostname, ":")
	if len(parts) != 2 {
		return errors.New("Invalid hostname port combination given")
	}

	re, _ := regexp.Compile(`^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\-]*[A-Za-z0-9])$`)
	if !re.MatchString(parts[0]) {
		return errors.New("Hostname does not appear to be a valid format")
	}

	// Check the port range.
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return errors.New("Port does not appear to be valid")
	}

	if port < 1 || port > 65535 {
		return errors.New("Port is not within the valid range")
	}

	return nil
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

	err := validateUrl(hostname, querypath)

	if err != nil {
		response = APIResponse{
			Status: HttpStatus{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)},
			Data:   APIError{Error: err.Error()},
		}
	} else {
		// We had to enable "UseEncodedPath" on Gorilla MUX to fix an issue where path
		// matching failed if the querypath started with an encoded/escaped character (ie, %2F)
		// Because of this, we need to unescape/decode the path here.  We have error handling here, though
		// Gorilla MUX should return a 400 response on it's own if encoding is bad.
		decodedPath, err := url.PathUnescape(querypath)
		if err != nil {
			utils.LogError(utils.LogFields{}, err, "Failed to decode query path")

			response = APIResponse{
				Status: HttpStatus{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)},
				Data:   APIError{Error: err.Error()},
			}

		} else
		// Successful decoding of the query path.
		// Perform service call to check URL.
		{
			urlService := services.UrlService{
				Hostname: hostname,
				Path:     decodedPath,
				Database: data.NewMongoDB("mongo", "urlinfo", "urls"),
				Cache:    data.NewMemcache("memcache:11211", 30),
			}

			urlStatus, err := urlService.FindUrl()

			if err != nil {
				utils.LogError(utils.LogFields{"hostname": hostname, "path": decodedPath}, err, "Error getting URL")
				response = APIResponse{
					Status: HttpStatus{Code: http.StatusInternalServerError, Message: http.StatusText(http.StatusInternalServerError)},
					Data:   nil,
				}
			} else {
				response = APIResponse{
					Status: HttpStatus{Code: http.StatusOK, Message: http.StatusText(http.StatusOK)},
					Data:   urlStatus,
				}
			}
		}
	}

	http_respond(response, writer)
}

// addUrl adds a new host / path combination to the system.  The new entry will be
// added first to the database then to the cache.
func addUrl(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	hostname := vars["hostname"]
	querypath := vars["querypath"]

	var response APIResponse

	err := validateUrl(hostname, querypath)
	if err != nil {
		response = APIResponse{
			Status: HttpStatus{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)},
			Data:   APIError{Error: err.Error()},
		}
	} else {
		// We had to enable "UseEncodedPath" on Gorilla MUX to fix an issue where path
		// matching failed if the querypath started with an encoded/escaped character (ie, %2F)
		// Because of this, we need to unescape/decode the path here.  We have error handling here, though
		// Gorilla MUX should return a 400 response on it's own if encoding is bad.
		decodedPath, err := url.PathUnescape(querypath)
		if err != nil {
			utils.LogError(utils.LogFields{}, err, "Failed to decode query path")

			response = APIResponse{
				Status: HttpStatus{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)},
				Data:   APIError{Error: err.Error()},
			}

		} else
		// Successful decoding of the query path.
		// Perform service call to check URL.
		{
			urlService := services.UrlService{
				Hostname: hostname,
				Path:     decodedPath,
				Database: data.NewMongoDB("mongo", "urlinfo", "urls"),
				Cache:    data.NewMemcache("memcache:11211", 30),
			}

			err = urlService.AddUrl()
			if err != nil {
				response = APIResponse{
					Status: HttpStatus{Code: http.StatusBadRequest, Message: http.StatusText(http.StatusBadRequest)},
					Data:   APIError{Error: err.Error()},
				}
			} else {
				response = APIResponse{
					Status: HttpStatus{Code: http.StatusOK, Message: http.StatusText(http.StatusOK)},
					Data:   "Successfully added URL",
				}
			}
		}
	}

	http_respond(response, writer)
}

// Main function
func main() {
	config, err := LoadConfig()
	if err != nil {
		os.Exit(1)
	}

	if config.Debug {
		utils.SetDebug()
	}

	router := mux.NewRouter().StrictSlash(true)
	router.UseEncodedPath()

	// Query URL info.  query path is optional.
	router.HandleFunc("/urlinfo/1/{hostname}/{querypath:.*}", checkUrl).Methods("GET")

	// Add a new URL. query path is optional.
	router.HandleFunc("/urlinfo/1/{hostname}/{querypath:.*}", addUrl).Methods("PUT")

	utils.LogInfo(utils.LogFields{"port": config.Port}, "Starting HTTP service")
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), handlers.LoggingHandler(os.Stdout, router))
}
