package data

import (
    "errors"
    "urlcheck/models"
)

var NotFoundError = errors.New("URL not found")

type DBInterface interface {
    // Find a URL by hostname/port and path/query
    FindUrl(string, string) (*models.Urls, error)
}

type CacheInterface interface {

    // Get a cache value by key.
    Get(string, string) (string, error)

    // Set a cache key value pair.
    Set(string, string) (error)

    // Delete a cache key
    Delete(string, string)
}
