package main

import (
	"github.com/koding/multiconfig"

	"urlcheck/data"
	"urlcheck/utils"
)

// Configuration instance.
var config *Config

// Configuration options
type Config struct {
	Port     int                  `json:"port"  default:"8010"`
	Debug    bool                 `json:"debug" default:"false"`

	Mongo    data.MongoDBConfig   `json:"mongo"`
	Memcache data.MemcachedConfig `json:"memcache"`
}

// Load the configuration from tag defaults and environment variables.
func LoadConfig() (*Config, error) {
	configuration := &Config{}

	loader := multiconfig.New()
	err := loader.Load(configuration)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load configuration")
		return nil, err
	}

	return configuration, nil
}
