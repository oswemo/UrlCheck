package data

// Memcache data handling.

import (
    "urlcheck/utils"

    "fmt"
    "github.com/bradfitz/gomemcache/memcache"
)

// Memcached structure
type Memcached struct {
    Client     *memcache.Client
    Servers    string
	Expiration int32
}

// New returns a new memcached object
func NewMemcache(servers string, expiration int32) (Memcached) {
    memcached := Memcached{
        Servers: servers,
        Expiration: expiration,
    }

    utils.LogInfo(utils.LogFields{"servers": servers, "expiration": expiration}, "Creating connection to memcached")
    memcached.Client = memcache.New(servers)
    return memcached
}

// Get a cache value by key.
func (m Memcached) Get(hostname string, path string) (string, error) {
    key := m.cacheKey(hostname, path)

    item, err := m.Client.Get(key)

    // Handle MISS without having to include memcache code everywhere.
    if err != nil {
        if  err == memcache.ErrCacheMiss {
            return "", NotFoundError
        }

        return "", err
    }

    return string(item.Value), nil
}

// Set a cache key value pair.
func (m Memcached) Set(hostname string, path string) (error) {
    // key := m.cacheKey(hostname, path)
    return nil
}

// Delete a cache key
func (m Memcached) Delete(hostname string, path string) {
    // key := m.cacheKey(hostname, path)
}

// cacheKey returns a formatted key for storing cache items.
func (m Memcached) cacheKey(hostname string, path string) (string) {
    return fmt.Sprintf("%s%s", hostname, path)
}
