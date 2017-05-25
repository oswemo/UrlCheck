package main

import (
	"errors"
	"github.com/koding/multiconfig"

	"urlcheck/data"
	"urlcheck/utils"
)

// Configuration instance.
var config *Config

// Configuration options
type Config struct {
	Port  int  `json:"port"      default:"8010"`
	Debug bool `json:"debug"     default:"false"`

	DBType    string `json:"dbtype"    default:""`
	CacheType string `json:"cachetype" default:""`
}

// validateConfig attempts to validate configuration where possible.
func validateConfig(config *Config) error {
	// Require a database definition.
	if config.DBType == "" {
		return errors.New("Invalid configuration.  Database type required.")
	}

	// Validate the database type.
	_, err := data.SelectDB(config.DBType)
	if err != nil {
		return errors.New("Invalid configuration.  Invalid database type defined.")
	}

	// Validate the cache type.
	_, err = data.SelectCache(config.CacheType)
	if err != nil {
		return errors.New("Invalid configuration.  Invalid cache type defined.")
	}

	return nil
}

// Load the configuration from tag defaults and environment variables.
func loadConfig() (*Config, error) {
	configuration := &Config{}

	tagLoader := multiconfig.TagLoader{}
	err := tagLoader.Load(configuration)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load tag configuration")
		return nil, err
	}

	envLoader := multiconfig.EnvironmentLoader{Prefix: "CONFIG"}
	err = envLoader.Load(configuration)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load environment configuration")
		return nil, err
	}

	err = validateConfig(configuration)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to validate configuration")
		return nil, err
	}
	return configuration, nil
}

// ShowConfig shows all configuration environment variables.
func showConfig() {
	configuration := &Config{}

	envLoader := multiconfig.EnvironmentLoader{Prefix: "CONFIG"}
	err := envLoader.Load(configuration)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load environment configuration")
		return
	}

	envLoader.PrintEnvs(configuration)

}
