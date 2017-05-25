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

// Return the selected database backend or an error if the type is invalid.
func SelectDB(dbType string) (DBInterface, error) {
	var db DBInterface
	var err error

	switch dbType {
	case "mongodb":
		db, err = NewMongoDB()
	case "sqldb":
		db, err = NewSQLDB()
	}

	if err != nil {
		return nil, err
	}

	if db == nil {
		utils.LogDebug(utils.LogFields{"dbtype": dbType}, "Invalid DB type")
		return nil, errors.New("Invalid database type")
	}

	return db, nil
}

// Return the selected cache backend or an error if the type is invalid.
func SelectCache(cacheType string) (CacheInterface, error) {

	switch cacheType {
	case "":
		utils.LogDebug(utils.LogFields{}, "Continuing with cache disabled")
		return nil, nil

	case "memcached":
		return NewMemcached(), nil
	}

	utils.LogDebug(utils.LogFields{"cachetype": cacheType}, "Invalid cache type")
	return nil, errors.New("Invalid cache type")
}
