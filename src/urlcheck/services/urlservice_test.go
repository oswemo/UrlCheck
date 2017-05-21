package services

import "testing"
import "errors"
import "urlcheck/data"
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

    return nil, errors.New("No matching url found.")
}

// TestFindUrl tests the FindUrl() method
func TestFindUrl(t *testing.T) {

    testCases := []struct{
        Hostname       string
        Path           string
        Database       data.DBInterface
        ExpectedResult bool
    }{
        {
            Hostname:       "www.example.com:80",
            Path:           "/foo/bar?query=something",
            Database:       MockDB{},
            ExpectedResult: true,
        },
        {
            Hostname:       "www.evildoers.net:8080",
            Path:           "/install/ransomwhere?bequiet=true",
            Database:       MockDB{},
            ExpectedResult: false,
        },
    }

    // Run through test cases.
    for _, c := range testCases {
        urlService := UrlService{
            Hostname: c.Hostname,
            Path:     c.Path,
            Database: c.Database,
        }

        urlStatus := urlService.FindUrl()

        if urlStatus.Safe != c.ExpectedResult {
            t.Errorf("%s and %s should be %t, but got %t", c.Hostname, c.Path, c.ExpectedResult, urlStatus.Safe)
        }
    }
}
