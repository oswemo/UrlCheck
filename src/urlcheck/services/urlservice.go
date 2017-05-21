package services

import (
    "urlcheck/data"
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
    Safe    bool         `json:"safe"`
}

// FindUrl attempts to find a given hostname/port and path/query combination in
// the backend storage.
// If found, a UrlInfoStatus object is returned with a "true" value.
// If not found, a UrlInfoStatus object is returned with a "false" value.
//
func (u *UrlService) FindUrl() (UrlInfoStatus) {
    status := UrlInfoStatus{ Safe: true }

    _, err := u.Database.FindUrl(u.Hostname, u.Path); if err == nil {
        // No error, so the URL was found in the database.
        status.Safe = false
    }
    return status
}
