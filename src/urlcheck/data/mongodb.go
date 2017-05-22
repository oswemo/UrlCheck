package data

// MongoDB database support.
// TODO: Add SSL/TLS and authentication.

import (
    "urlcheck/models"

    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

// MongoDB Struct, defines options for MongoDB database connections.
type MongoDB struct {
    URL        string  // MongoDB URL
    DBName     string  // MongoDB DB name
    Collection string  // Collection name
}

// FindUrl attempts to look up the URL in the MongoDB collection.
// If a matching entry is found, the point to a models.Urls object is returned.
// If no entry is found, an error object is returned.
func (m MongoDB) FindUrl(hostname string, path string) (*models.Urls, error) {
    session, err := mgo.Dial(m.URL) ; if err != nil {
        return nil, err
    }

    query := bson.M{
        "hostport": hostname,
        "pathquery": path,
    }

    c := session.DB(m.DBName).C(m.Collection)

    result := models.Urls{}
    err = c.Find(query).One(&result) ; if err != nil {
        return nil, err
    }

    return &result, nil
}
