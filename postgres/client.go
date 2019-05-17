package postgres

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // per gorm
)

var (
	// ErrRecordNotFound record not found error, happens when haven't find any matched data when looking up with a struct
	ErrRecordNotFound = errors.New("record not found")
)

// ErrDatabaseGeneral is a generic error wrapper for unexplained errors
type ErrDatabaseGeneral string

func (edg ErrDatabaseGeneral) Error() string {
	return fmt.Sprintf("General Database Error: %s", edg)
}

type DBClient struct {
	db *gorm.DB
}

type GOTD struct {
	gorm.Model
	GIF string
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

	if !db.HasTable(&GOTD{}) {
		db.CreateTable(&GOTD{})
	}

	return &DBClient{
		db: db,
	}, nil
}

func (c *DBClient) Insert(gif *GOTD) error {
	if result := c.db.Create(gif); result.Error != nil {
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (c *DBClient) Update(gif *GOTD) error {
	if result := c.db.Model(&GOTD{}).Updates(gif); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return ErrRecordNotFound
		}
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	c.db.Close()
}
