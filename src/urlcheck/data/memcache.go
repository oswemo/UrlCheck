package data

// Memcache data handling.

import (
	"urlcheck/utils"

	"fmt"

	"github.com/koding/multiconfig"
	"github.com/bradfitz/gomemcache/memcache"
)

// Memcached structure
type Memcached struct {
	Config *MemcachedConfig
	Client *memcache.Client
}

type MemcachedConfig struct {
	Servers    string `json:"servers"    default:"memcache:11211"`    // MEMCACHEDCONFIG_SERVERS
	Expiration int    `json:"expiration" default:"300"`               // MEMCACHEDCONFIG_EXPIRATION
}

// New returns a new memcached object
func NewMemcache() (*Memcached) {
	config := &MemcachedConfig{}

	loader := multiconfig.New()
	err := loader.Load(config)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load configuration")
		return nil
	}

	memcached := &Memcached{ Config: config }

	utils.LogInfo(utils.LogFields{"servers": config.Servers, "expiration": config.Expiration}, "Creating connection to memcached")
	memcached.Client = memcache.New(config.Servers)
	return memcached
}

// Get a cache value by key.
func (m Memcached) Get(hostname string, path string) (string, error) {
	key := m.cacheKey(hostname, path)

	item, err := m.Client.Get(key)

	// Handle MISS without having to include memcache code everywhere.
	if err != nil {
		if err == memcache.ErrCacheMiss {
			return "", NotFoundError
		}

		return "", err
	}

	return string(item.Value), nil
}

// Set a cache key value pair.
func (m Memcached) Set(hostname string, path string) error {
	key := m.cacheKey(hostname, path)
	return m.Client.Set(&memcache.Item{Key: key, Value: []byte("exists"), Expiration: int32(m.Config.Expiration) })

}

// cacheKey returns a formatted key for storing cache items.
func (m Memcached) cacheKey(hostname string, path string) string {
	return fmt.Sprintf("%s%s", hostname, path)
}
