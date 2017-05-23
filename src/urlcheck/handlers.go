package main

// MUX handlers.

import (
	"errors"
	"github.com/gorilla/mux"
	"net/http"

	"urlcheck/data"
	"urlcheck/services"
	"urlcheck/utils"
)

// checkUrl attempts to validate whether a hostname and query path constitute a known
// malware site.  It is expected that the query path is a URI encoded string.
// An HTTP response is returned vai the writer parameter.
func checkUrl(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	hostname := vars["hostname"]
	querypath := vars["querypath"]

	response := APIResponse{}
	err := utils.ValidateUrl(hostname)
	if err != nil {
		response.BadRequest(err)
		http_respond(response, writer)
		return
	}

	decodedPath, err := utils.URLDecode(querypath)
	if err != nil {
		response.BadRequest(err)
		http_respond(response, writer)
		return
	}

	urlService := services.UrlService{
		Hostname: hostname,
		Path:     decodedPath,
		Database: data.NewMongoDB(),
		Cache:    data.NewMemcache(),
	}

	urlStatus, err := urlService.FindUrl()
	if err != nil {
		utils.LogError(utils.LogFields{"hostname": hostname, "path": decodedPath}, err, "Error getting URL")
		response.InternalError(errors.New("An error occurred"))
	} else {
		response.Success(urlStatus)
	}

	http_respond(response, writer)
}

// addUrl adds a new host / path combination to the system.  The new entry will be
// added first to the database then to the cache.
func addUrl(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	hostname := vars["hostname"]
	querypath := vars["querypath"]

	response := APIResponse{}
	err := utils.ValidateUrl(hostname)
	if err != nil {
		response.BadRequest(err)
		http_respond(response, writer)
		return
	}

	decodedPath, err := utils.URLDecode(querypath)
	if err != nil {
		response.BadRequest(err)
		http_respond(response, writer)
		return
	}

	urlService := services.UrlService{
		Hostname: hostname,
		Path:     decodedPath,
		Database: data.NewMongoDB(),
		Cache:    data.NewMemcache(),
	}

	err = urlService.AddUrl()
	if err != nil {
		response.BadRequest(err)
	} else {
		response.Success(StringData{Message: "Successfully added URL"})
	}

	http_respond(response, writer)
}
