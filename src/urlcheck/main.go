package main

import (
	"flag"
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

	// We aren't using flags with multiconfig as it messes with "go test -v" and
	// we don't really need it since are primarily using environment variables.
	// We'll map -envs manually to a function to output the variables.
	var showEnvs = flag.Bool("envs", false, "Show environment variables")
	flag.Parse()

	if *showEnvs {
		showConfig()
		os.Exit(0)
	}

	config, err = loadConfig()
	if err != nil {
		os.Exit(1)
	}

	utils.LogInfo(utils.LogFields{"database": config.DBType, "cache": config.CacheType, "port": config.Port}, "Configuration loaded")

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
