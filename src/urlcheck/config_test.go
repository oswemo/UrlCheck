package main

import "os"
import "fmt"
import "testing"
import "urlcheck/utils"

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		Port      int
		DBType    string
		CacheType string
		Valid     bool
	}{
		{
			Port:      8080,
			DBType:    "mongodb",
			CacheType: "memcached",
			Valid:     true,
		},
		{
			Port:      8080,
			DBType:    "mongodb",
			CacheType: "",
			Valid:     true,
		},
		{
			Port:      8080,
			DBType:    "",
			CacheType: "",
			Valid:     false,
		},
		{
			Port:      8080,
			DBType:    "FakeDB",
			CacheType: "",
			Valid:     false,
		},
		{
			Port:      8080,
			DBType:    "mongodb",
			CacheType: "FakeCache",
			Valid:     false,
		},
	}

	utils.SetFatal()
	for _, c := range testCases {
		config := &Config{
			Port:      c.Port,
			Debug:     true,
			DBType:    c.DBType,
			CacheType: c.CacheType,
		}

		err := validateConfig(config)
		if c.Valid && err != nil {
			t.Errorf("Validation failed.  Expected a valid configuration with DB [%s] and Cache [%s]", c.DBType, c.CacheType)
		}
		if !c.Valid && err == nil {
			t.Errorf("Validation failed.  Expected an invalid configuration with DB [%s] and Cache [%s]", c.DBType, c.CacheType)
		}
	}
}

func TestLoadConfig(t *testing.T) {
	testCases := []struct {
		Port      int
		DBType    string
		CacheType string
		Debug     bool
	}{
		{
			Port:      8080,
			DBType:    "mongodb",
			CacheType: "memcached",
			Debug:     true,
		},
		{
			Port:      9000,
			DBType:    "mongodb",
			CacheType: "",
			Debug:     false,
		},
	}

	utils.SetFatal()
	for _, c := range testCases {

		os.Setenv("CONFIG_PORT", fmt.Sprintf("%d", c.Port))
		os.Setenv("CONFIG_DBTYPE", c.DBType)
		os.Setenv("CONFIG_CACHETYPE", c.CacheType)
		os.Setenv("CONFIG_DEBUG", fmt.Sprintf("%t", c.Debug))

		config, err := loadConfig()
		if err != nil {
			t.Error("Configuration loading failed.", err)
		}

		if config.Port != c.Port ||
			config.DBType != c.DBType ||
			config.CacheType != c.CacheType ||
			config.Debug != c.Debug {
			t.Errorf("Configuration options to not match as expected.")
			fmt.Printf("%v\n", config)
		}
	}
}
