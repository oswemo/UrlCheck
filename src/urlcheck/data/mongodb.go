package data

// MongoDB database support.
// TODO: Add SSL/TLS and authentication.

import (
	"urlcheck/models"
	"urlcheck/utils"

	"errors"
	"time"

	"github.com/koding/multiconfig"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB Struct, defines options for MongoDB database connections.
type MongoDB struct {
	Config  *MongoDBConfig
	Session *mgo.Session
}

// MongoDBConfig defines configuration options for multiconfig
type MongoDBConfig struct {
	URL        string `json:"url"        default:"mongodb"` // MONGODB_URL
	Database   string `json:"database"   default:"urlinfo"` // MONGODB_DATABASE
	Collection string `json:"collection" default:"urls"`    // MONGODB_COLLECTION
	Timeout    int    `json:"timeout"    default:"2"`       // MONGODB_TIMEOUT
}

// UrlSchemaMongoDB struct defines the layout of the MongoDB collection data.
type UrlSchemaMongoDB struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	HostPort  string        `json:"hostport"`
	PathQuery string        `json:"pathquery"`
}

// NewMongoDB returns a new instance of the MongoDB struct.
func NewMongoDB() *MongoDB {

	config := &MongoDBConfig{}

	tagLoader := multiconfig.TagLoader{}
	err := tagLoader.Load(config)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load tag configuration for MongoDB")
		return nil
	}

	envLoader := multiconfig.EnvironmentLoader{Prefix: "MONGODB"}
	err = envLoader.Load(config)
	if err != nil {
		utils.LogError(utils.LogFields{}, err, "Failed to load environment configuration for MongoDB")
		return nil
	}

	mongodb := &MongoDB{Config: config}
	mongodb.Connect()
	return mongodb
}

// Connect connects to the MongoDB server
func (m *MongoDB) Connect() error {
	var err error

	// Set up the MongoDB session.  Timeout set for connection and subsequent queries
	// to limit requests having to wait for a response.
	utils.LogInfo(utils.LogFields{"url": m.Config.URL, "timeout": m.Config.Timeout}, "Creating connection to MongoDB")
	timeout := time.Duration(m.Config.Timeout) * time.Second
	m.Session, err = mgo.DialWithTimeout(m.Config.URL, timeout)

	if err != nil {
		utils.LogError(utils.LogFields{"url": m.Config.URL}, err, "Error connecting to MongoDB")
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

	c := m.Session.DB(m.Config.Database).C(m.Config.Collection)

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

	doc := UrlSchemaMongoDB{HostPort: hostname, PathQuery: path}
	c := m.Session.DB(m.Config.Database).C(m.Config.Collection)
	err := c.Insert(&doc)
	return err
}
