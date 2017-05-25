package services

import "testing"
import "errors"
import "urlcheck/data"
import "urlcheck/utils"
import "urlcheck/models"

// Initial DBData data
var DBData = []models.Urls{
	models.Urls{Hostname: "www.example.com:80", Path: "/some/path?dostuff=example"},
	models.Urls{Hostname: "www.evildoers.net:8080", Path: "/install/ransomwhere?bequiet=true"},
}

// Initial Cache data
var CacheData = []models.Urls{
	models.Urls{Hostname: "www.example.com:80", Path: "/some/path?dostuff=example"},
}

// Data mock to test without requiring MongoDB directly.
type MockDB struct{}

func (m MockDB) FindUrl(hostname string, path string) (*models.Urls, error) {

	for _, url := range DBData {
		if url.Hostname == hostname && url.Path == path {
			return &url, nil
		}
	}

	return nil, data.NotFoundError
}

// AddUrl Mocks the AddUrl database function.  In this case, we actually add
// the url to the slice so that we can test multiple attempts at adding the same
// URL
func (m MockDB) AddUrl(hostname string, path string) error {

	for _, url := range DBData {
		if url.Hostname == hostname && url.Path == path {
			return errors.New("URL already exists")
		}
	}

	dataLen := len(DBData)
	newData := make([]models.Urls, dataLen+1)
	copy(newData, DBData)
	DBData = newData
	DBData[dataLen] = models.Urls{Hostname: hostname, Path: path}

	return nil
}

// Cache mock to test without requiring Memcached directly.
type MockCache struct{}

func (m MockCache) Get(hostname string, path string) (string, error) {
	for _, url := range CacheData {
		if url.Hostname == hostname && url.Path == path {
			return "exists", nil
		}
	}

	return "", data.NotFoundError
}

// Set a cache key value pair.
func (m MockCache) Set(hostname string, path string) error {
	dataLen := len(CacheData)
	newData := make([]models.Urls, dataLen+1)
	copy(newData, CacheData)
	CacheData = newData
	CacheData[dataLen] = models.Urls{Hostname: hostname, Path: path}
	return nil
}

// TestFindUrl tests the FindUrl() method
func TestFindUrl(t *testing.T) {

	testCases := []struct {
		Hostname string
		Path     string
		Database data.DBInterface
		Cache    data.CacheInterface
		IsSafe   bool
	}{
		{
			Hostname: "www.example.com:80",
			Path:     "/some/path?dostuff=example",
			Database: MockDB{},
			Cache:    MockCache{},
			IsSafe:   false,
		},
		{
			Hostname: "www.example.com:80",
			Path:     "/foo/bar?query=something",
			Database: MockDB{},
			Cache:    MockCache{},
			IsSafe:   true,
		},
		{
			Hostname: "www.evildoers.net:8080",
			Path:     "/install/ransomwhere?bequiet=true",
			Database: MockDB{},
			Cache:    MockCache{},
			IsSafe:   false,
		},
	}

	utils.SetFatal()

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
	utils.SetDebug()

	testCases := []struct {
		Hostname string
		Path     string
		Database data.DBInterface
		Cache    data.CacheInterface
		Success  bool
	}{
		{
			Hostname: "www.example.com:80",
			Path:     "/some/path?dostuff=example",
			Database: MockDB{},
			Cache:    MockCache{},
			Success:  false,
		},
		{
			Hostname: "www.example.com:80",
			Path:     "/foo/bar?query=something",
			Database: MockDB{},
			Cache:    MockCache{},
			Success:  true,
		},
		{
			Hostname: "www.example.com:80",
			Path:     "/foo/bar?query=something",
			Database: MockDB{},
			Cache:    MockCache{},
			Success:  false,
		},
		{
			Hostname: "www.evildoers.net:8080",
			Path:     "/install/ransomwhere?bequiet=true",
			Database: MockDB{},
			Cache:    MockCache{},
			Success:  false,
		},
	}

	utils.SetFatal()

	// Run through test cases.
	for _, c := range testCases {
		urlService := UrlService{
			Hostname: c.Hostname,
			Path:     c.Path,
			Database: c.Database,
			Cache:    c.Cache,
		}

		err := urlService.AddUrl()
		if err != nil && c.Success {
			t.Errorf("Adding %s %s failed, but should have succeeded.", c.Hostname, c.Path)
		}
		if err == nil && !c.Success {
			t.Errorf("Adding %s %s succeeded, but should have failed.", c.Hostname, c.Path)
		}
	}

}
