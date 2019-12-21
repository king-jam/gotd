// Package postgres provides a gorm DB instance for postgres
package postgres

import (
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // per gorm
)

// NewClient takes a connection string to pass into the Database
func NewClient(connectionString string) (*gorm.DB, error) {
	dbURL, err := url.Parse(connectionString)
	if err != nil {
		return nil, fmt.Errorf("invalid database URL format: %s", err)
	}

	db, err := gorm.Open(dbURL.Scheme, dbURL.String())
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(20)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(20)

	return db, nil
}
