package data

import "urlcheck/models"

type DBInterface interface {
    FindUrl(string, string) (*models.Urls, error)
}

type CacheInterface interface {
    GetKey(string) (*string, error)
}
