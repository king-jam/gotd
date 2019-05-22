package postgres

import (
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // per gorm
)

const minimumHistoryThresholdMins = 10 * time.Minute

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
	Db *gorm.DB
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
		Db: db,
	}, nil
}

// func (c *DBClient) Insert(gif *GifHistory) error {
// 	if result := c.db.Create(gif); result.Error != nil {
// 		return ErrDatabaseGeneral(result.Error.Error())
// 	}

// 	//Debugging
// 	mrGif, _ := c.LatestGIF()
// 	log.Printf("New History ID: %d", mrGif.ID)
// 	log.Printf("Tags: +%v", mrGif.Tags)

// 	return nil
// }

// // Update will update a gif from the database
// func (c *DBClient) Update(gif *GifHistory) error {
// 	if result := c.db.Model(&GifHistory{}).Updates(gif); result.Error != nil {
// 		if gorm.IsRecordNotFoundError(result.Error) {
// 			return ErrRecordNotFound
// 		}
// 		return ErrDatabaseGeneral(result.Error.Error())
// 	}
// 	return nil
// }

// func (c *DBClient) FindGIFByID(id uint) (*GifHistory, error) {
// 	return &GifHistory{}, nil
// }

// func (c *DBClient) FindAllGifs() ([]GifHistory, error) {
// 	return []GifHistory{}, nil
// }

// // LatestGIF will return the latest gif from database
// func (c *DBClient) LatestGIF() (*GifHistory, error) {
// 	gif := new(GifHistory)
// 	if result := c.db.Model(&GifHistory{}).Last(gif); result.Error != nil {
// 		if gorm.IsRecordNotFoundError(result.Error) {
// 			return nil, ErrRecordNotFound
// 		}
// 		return nil, ErrDatabaseGeneral(result.Error.Error())
// 	}
// 	return gif, nil
// }

// func (c *DBClient) DeleteGIFByID(id int) error {
// 	return nil
// }

// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	c.Db.Close()
}
