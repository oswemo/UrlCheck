package data

import "testing"

func TestSelectDB(t *testing.T) {
	testCases := []struct {
		DBType string
		Valid  bool
	}{
		{DBType: "mongodb", Valid: true},
		{DBType: "testing", Valid: false},
	}

	for _, c := range testCases {
		result := SelectDB(c.DBType)
		if result == nil && c.Valid {
			t.Errorf("Expected DB type %s to be valid", c.DBType)
		}
		if result != nil && !c.Valid {
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

	for _, c := range testCases {
		result := SelectCache(c.CacheType)
		if result == nil && c.Valid {
			t.Errorf("Expected cache type %s to be valid", c.CacheType)
		}
		if result != nil && !c.Valid {
			t.Errorf("Expected cache type %s to be invalid", c.CacheType)
		}
	}
}
