package main

import (
	"fmt"
	"net/http"
	"os"

	"urlcheck/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// Main function
func main() {
	var err error

	config, err = LoadConfig()
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
