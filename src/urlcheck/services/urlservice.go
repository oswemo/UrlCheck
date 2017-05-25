package services

import (
	"errors"

	"urlcheck/data"
	"urlcheck/utils"
)

// UrlService definition
type UrlService struct {
	Hostname string
	Path     string

	Database data.DBInterface
	Cache    data.CacheInterface
}

// UrlInfoStatus is the response object for requests to the IsSafe method.
type UrlInfoStatus struct {
	Safe bool `json:"safe"`
}

func NewUrlService(hostname string, path string, dbType string, cacheType string) (*UrlService, error) {

	database, err := data.SelectDB(dbType)
	if err != nil {
		return nil, err
	}

	cache, err := data.SelectCache(cacheType)
	if err != nil {
		return nil, err
	}

	service := &UrlService{
		Hostname: hostname,
		Path:     path,
		Database: database,
		Cache:    cache,
	}

	return service, nil
}

// FindUrl attempts to find a given hostname/port and path/query combination in
// the backend storage.
// If found, a UrlInfoStatus object is returned with a "true" value.
// If not found, a UrlInfoStatus object is returned with a "false" value.
//
func (u *UrlService) FindUrl() (*UrlInfoStatus, error) {
	status := &UrlInfoStatus{Safe: true}

	var err error

	// Assuming cache endpoint is defined, look up in cache.
	if u.Cache != nil {
		_, err = u.Cache.Get(u.Hostname, u.Path)

		// if the key is in the cache, then further work is not needed.
		if err == nil {
			utils.LogDebug(utils.LogFields{"hostname": u.Hostname, "path": u.Path, "safe": status.Safe}, "A match was found in the cache")
			status.Safe = false
			return status, nil
		}
	}

	// The key was not found in the cache, cache is not being used, or cache lookup failed with an error.
	if u.Database != nil {
		_, err := u.Database.FindUrl(u.Hostname, u.Path)

		// A matching URL was found in the database.
		if err == nil {
			utils.LogDebug(utils.LogFields{"hostname": u.Hostname, "path": u.Path, "safe": status.Safe}, "A match was found in the database")
			status.Safe = false

			//Update cache.
			if u.Cache != nil {
				u.Cache.Set(u.Hostname, u.Path)
			}

			return status, nil
		}

		// Database returned a NotFoundError
		if err == data.NotFoundError {
			utils.LogDebug(utils.LogFields{"hostname": u.Hostname, "path": u.Path, "safe": status.Safe}, "No matching URL found")
		}

		// Database returned some other (connection failure, etc)
		if err != data.NotFoundError {
			utils.LogError(utils.LogFields{"hostname": u.Hostname, "path": u.Path}, err, "Database error looking up URL")
			return nil, errors.New("Error looking up URL")
		}
	}

	return status, nil
}

// AddUrl adds a new hostname/port and path/query combination to the database and cache.
// Returns an error if something goes wrong.
func (u *UrlService) AddUrl() error {
	if u.Database != nil {
		_, err := u.Database.FindUrl(u.Hostname, u.Path)
		if err == nil {
			return data.AlreadyExistsError
		}

		err = u.Database.AddUrl(u.Hostname, u.Path)
		if err != nil {
			utils.LogError(utils.LogFields{"hostname": u.Hostname, "path": u.Path}, err, "Error adding new URL")
			return data.AddFailureError
		}
	}

	if u.Cache != nil {
		//Update cache.
		u.Cache.Set(u.Hostname, u.Path)
	}

	return nil
}
