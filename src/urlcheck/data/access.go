package data

import (
	"errors"

	"urlcheck/models"
	"urlcheck/utils"
)

var NotFoundError = errors.New("URL not found")
var AddFailureError = errors.New("Error adding new URL")
var AlreadyExistsError = errors.New("URL already exists")

type DBInterface interface {
	// Find a URL by hostname/port and path/query
	FindUrl(string, string) (*models.Urls, error)

	// Add a new URL to the system.
	AddUrl(string, string) error
}

type CacheInterface interface {

	// Get a cache value by key.
	Get(string, string) (string, error)

	// Set a cache key value pair.
	Set(string, string) error
}

// Return the selected database backend
func SelectDB(dbType string) DBInterface {
	switch dbType {
	default:
		utils.LogError(utils.LogFields{"dbtype": dbType}, errors.New("Invalid DB type"), "")
	case "mongodb":
		return NewMongoDB()
	}

	return nil
}

// Return the selected cache backend
func SelectCache(cacheType string) CacheInterface {
	switch cacheType {
	default:
		utils.LogError(utils.LogFields{"cachetype": cacheType}, errors.New("Invalid cache type"), "")
	case "memcached":
		return NewMemcached()
	}

	return nil
}
