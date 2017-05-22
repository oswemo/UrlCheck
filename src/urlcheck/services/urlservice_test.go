package services

import "testing"
import "urlcheck/data"
import "urlcheck/utils"
import "urlcheck/models"

// Data mock to test without requiring MongoDB directly.
type MockDB struct { }
func (m MockDB) FindUrl(hostname string, path string) (*models.Urls, error) {

    // Test data for mocking the DB
    TestData := []models.Urls{
        models.Urls{ Hostname: "foo.example.com:80", Path: "/some/path?dostuff=example" },
        models.Urls{ Hostname: "www.evildoers.net:8080", Path: "/install/ransomwhere?bequiet=true" },
    }

    for _, url := range TestData {
        if url.Hostname == hostname && url.Path == path {
            return &url, nil
        }
    }

    return nil, data.NotFoundError
}

// Cache mock to test without requiring Memcached directly.
type MockCache struct { }
func (m MockCache) Get(hostname string, path string) (string, error) {
    // Test data for mocking the DB
    TestData := []models.Urls{
        models.Urls{ Hostname: "www.example.com:80", Path: "/some/path?dostuff=example" },
    }

    for _, url := range TestData {
        if url.Hostname == hostname && url.Path == path {
            return "exists", nil
        }
    }

    return "", data.NotFoundError
}

// Set a cache key value pair.
func (m MockCache) Set(hostname string, path string) (error) {
    // key := m.cacheKey(hostname, path)
    return nil
}

// Delete a cache key
func (m MockCache) Delete(hostname string, path string) {
    // key := m.cacheKey(hostname, path)
}

// TestFindUrl tests the FindUrl() method
func TestFindUrl(t *testing.T) {

    utils.SetDebug()

    testCases := []struct{
        Hostname       string
        Path           string
        Database       data.DBInterface
        Cache          data.CacheInterface
        IsSafe         bool
    }{
        {
            Hostname:       "www.example.com:80",
            Path:           "/some/path?dostuff=example",
            Database:       MockDB{},
            Cache:          MockCache{},
            IsSafe:         false,
        },
        {
            Hostname:       "www.example.com:80",
            Path:           "/foo/bar?query=something",
            Database:       MockDB{},
            Cache:          MockCache{},
            IsSafe:         true,
        },
        {
            Hostname:       "www.evildoers.net:8080",
            Path:           "/install/ransomwhere?bequiet=true",
            Database:       MockDB{},
            Cache:          MockCache{},
            IsSafe:         false,
        },
    }

    // Run through test cases.
    for _, c := range testCases {
        urlService := UrlService{
            Hostname: c.Hostname,
            Path:     c.Path,
            Database: c.Database,
            Cache:    c.Cache,
        }

        urlStatus, err := urlService.FindUrl()
        if err != nil {
            t.Errorf("An unexpected error occurred.")
        }

        if urlStatus.Safe != c.IsSafe {
            t.Errorf("%s and %s should be %t, but got %t", c.Hostname, c.Path, c.IsSafe, urlStatus.Safe)
        }
    }
}

// TestAddUrl tests the AddUrl() method
func TestAddUrl(t *testing.T) {
}
