package postgres

import (
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // per gorm
)

const minimumHistoryThresholdMins = 10 * time.Minute

type DBClient struct {
	DB *gorm.DB
}

// InitDatabase takes a connection string URL to pass into the Database
func InitDatabase(url *url.URL) (*DBClient, error) {
	db, err := gorm.Open(url.Scheme, url.String())
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(20)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(20)

	return &DBClient{
		DB: db,
	}, nil
}

// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	c.DB.Close()
}
