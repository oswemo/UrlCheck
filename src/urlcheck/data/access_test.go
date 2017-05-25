package data

import "testing"
import "urlcheck/utils"

func TestSelectDB(t *testing.T) {
	testCases := []struct {
		DBType string
		Valid  bool
	}{
		{DBType: "mongodb", Valid: true},
		{DBType: "testing", Valid: false},
	}

	utils.SetFatal()

	for _, c := range testCases {
		_, err := SelectDB(c.DBType)
		if err != nil && c.Valid {
			t.Errorf("Expected DB type %s to be valid", c.DBType)
		}
		if err == nil && !c.Valid {
			t.Errorf("Expected DB type %s to be invalid", c.DBType)
		}
	}
}

func TestSelectCache(t *testing.T) {
	testCases := []struct {
		CacheType string
		Valid     bool
	}{
		{CacheType: "memcached", Valid: true},
		{CacheType: "testing", Valid: false},
	}

	utils.SetFatal()
	for _, c := range testCases {
		_, err := SelectCache(c.CacheType)
		if err != nil && c.Valid {
			t.Errorf("Expected cache type %s to be valid", c.CacheType)
		}
		if err == nil && !c.Valid {
			t.Errorf("Expected cache type %s to be invalid", c.CacheType)
		}
	}
}
