package data

// MongoDB database support.
// TODO: Add SSL/TLS and authentication.

import (
	"urlcheck/models"
	"urlcheck/utils"

	"errors"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB_TIMEOUT = 2
)

// MongoDB Struct, defines options for MongoDB database connections.
type MongoDB struct {
	URL        string // MongoDB URL
	DBName     string // MongoDB DB name
	Collection string // Collection name

	Session *mgo.Session
}

// UrlSchema struct defines the layout of the MongoDB collection data.
type UrlSchema struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	HostPort  string        `json:"hostport"`
	PathQuery string        `json:"pathquery"`
}

// NewMongoDB returns a new instance of the MongoDB struct.
func NewMongoDB(url string, dbname string, collection string) MongoDB {
	mongo := MongoDB{
		URL:        url,
		DBName:     dbname,
		Collection: collection,
	}

	(&mongo).Connect()
	return mongo
}

// Connect connects to the MongoDB server
func (m *MongoDB) Connect() error {
	var err error

	// Set up the MongoDB session.  Timeout set for connection and subsequent queries
	// to limit requests having to wait for a response.
	utils.LogInfo(utils.LogFields{"url": m.URL, "timeout": DB_TIMEOUT}, "Creating connection to mongodb")
	timeout := time.Duration(DB_TIMEOUT) * time.Second
	m.Session, err = mgo.DialWithTimeout(m.URL, timeout)

	if err != nil {
		utils.LogError(utils.LogFields{"url": m.URL}, err, "Error connecting to Mongo")
	}

	return err
}

// FindUrl attempts to look up the URL in the MongoDB collection.
// If a matching entry is found, the point to a models.Urls object is returned.
// If no entry is found, an error object is returned.
func (m MongoDB) FindUrl(hostname string, path string) (*models.Urls, error) {

	if m.Session == nil {
		return nil, errors.New("No active connection to the database")
	}

	// Simple query for the data
	query := bson.M{
		"hostport":  hostname,
		"pathquery": path,
	}

	c := m.Session.DB(m.DBName).C(m.Collection)

	result := models.Urls{}
	err := c.Find(query).One(&result)
	if err != nil {

		// Check if it's NotFound so that we can notify the caller properly.
		if err == mgo.ErrNotFound {
			return nil, NotFoundError
		}

		// Return the actual error if other than not found.
		return nil, err
	}

	return &result, nil
}

// AddUrl adds a hostname/port and path/query combination to the database.
// An error is returned if anything goes wrong.
func (m MongoDB) AddUrl(hostname string, path string) error {
	if m.Session == nil {
		return errors.New("No active connection to the database")
	}

	doc := UrlSchema{HostPort: hostname, PathQuery: path}
	c := m.Session.DB(m.DBName).C(m.Collection)
	err := c.Insert(&doc)
	return err
}
