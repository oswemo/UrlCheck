package data

// SQL Database support

import (
    "errors"

    "urlcheck/utils"
    "urlcheck/models"

    "github.com/jinzhu/gorm"
    "github.com/koding/multiconfig"

    _ "github.com/jinzhu/gorm/dialects/sqlite"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/mssql"
)

type SQLDB struct {
    Config     *SQLDBConfig
    Connection *gorm.DB
}

type SQLDBConfig struct {
    Dialect    string     // Dialect of SQL (mysql, postgres, mssql, or sqlite)
    Url        string     // Connection URL
}

// NewSQLDB creates and returns a new SQLDB object
func NewSQLDB() (*SQLDB, error) {
    config := &SQLDBConfig{}

    tagLoader := multiconfig.TagLoader{}
    err := tagLoader.Load(config)
    if err != nil {
        utils.LogError(utils.LogFields{}, err, "Failed to load tag configuration for SQLDB")
        return nil, err
    }

    envLoader := multiconfig.EnvironmentLoader{Prefix: "SQLDB"}
    err = envLoader.Load(config)
    if err != nil {
        utils.LogError(utils.LogFields{}, err, "Failed to load environment configuration for SQLDB")
        return nil, err
    }

    sqldb := &SQLDB{Config: config}

    err = sqldb.validateConfig() ; if err != nil {
        return nil, err
    }

    sqldb.connect()
    sqldb.migrate()
    return sqldb, nil
}

// ValidateConfig attempts to validate the configuration.
// Returns an error if validation fails.
func (s *SQLDB) validateConfig() (error) {
    if s.Config.Dialect == "" {
        return errors.New("Configuration error.  SQLDB_DIALECT is required.")
    }
    if s.Config.Url == "" {
        return errors.New("Configuration error.  SQLDB_URL is required.")
    }
    return nil
}


// Connect creates a connection to the SQL database.
func (s *SQLDB) connect() (error) {
    db, err := gorm.Open(s.Config.Dialect, s.Config.Url) ; if err != nil {
        return err
    }

    s.Connection = db
    return nil
}

// Migrate performs SQL schema migrations.
func (s *SQLDB) migrate() (error) {
    if s.Connection == nil {
        s.connect()
    }

    s.Connection.AutoMigrate(&models.Urls{})
    return nil
}


// FindUrl attempts to look up the URL in the SQL database.
// If a matching entry is found, the point to a models.Urls object is returned.
// If no entry is found, an error object is returned.
func (s SQLDB) FindUrl(hostname string, path string) (*models.Urls, error) {
    query := &models.Urls{ Hostname: hostname, Path: path }
    urlinfo := &models.Urls{}

    result := s.Connection.Where(query).First(urlinfo)
    if result.Error != nil {

        // Check if it's NotFound so that we can notify the caller properly.
        if result.RecordNotFound() {
            return nil, NotFoundError
        }

        // Return the actual error if other than not found.
        return nil, result.Error
    }


    return urlinfo, nil
}

// Add a new URL to the system.
func (s SQLDB) AddUrl(hostname string, path string) (error) {
    query := &models.Urls{ Hostname: hostname, Path: path }

    result := s.Connection.Save(query)
    if result.Error != nil {
        return result.Error
    }
    return nil
}
