package postgres

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" // per gorm
	"github.com/lib/pq"
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
	db *gorm.DB
}

// type CurrentGOTD struct {
// 	gorm.Model
// 	GIF string `json:"url"`
// }

type GifHistory struct {
	gorm.Model
	DeactivatedAt time.Time
	GIF           string `json:"url"`
	RequesterID   string
	RequestSrc    string
	Tags          pq.StringArray `gorm:"type:varchar(64)[]"`
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

	// if !db.HasTable(&CurrentGOTD{}) {
	// 	db.CreateTable(&CurrentGOTD{})
	// }

	if !db.HasTable(&GifHistory{}) {
		db.CreateTable(&GifHistory{})
	}

	return &DBClient{
		db: db,
	}, nil
}

func (c *DBClient) Update(gif *GifHistory) error {
	if result := c.db.Model(&GifHistory{}).Updates(gif); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return ErrRecordNotFound
		}
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (c *DBClient) Insert(gif *GifHistory) error {
	if result := c.db.Create(gif); result.Error != nil {
		return ErrDatabaseGeneral(result.Error.Error())
	}

	//Debugging
	mrGif, _ := c.LatestGIF()
	log.Printf("New History ID: %d", mrGif.ID)
	log.Printf("Tags: +%v", mrGif.Tags)

	return nil
}

// func (c *DBClient) UpdateGIF(gif *CurrentGOTD) error {
// 	current, err := c.LatestGIF()
// 	if err != nil {
// 		if err == ErrRecordNotFound {
// 			err = c.Insert(gif)
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		return err
// 	}
// 	duration := time.Since(current.UpdatedAt)
// 	if duration > minimumHistoryThresholdMins {
// 		tags, err := giphy.GetGIFTags(current.GIF)
// 		if err != nil {
// 			log.Print(err)
// 		}
// 		prevGif := GifHistory{
// 			GIF:         current.GIF,
// 			ElapsedTime: duration,
// 			Tags:        tags,
// 		}
// 		err = c.AddGifHistory(&prevGif)
// 		if err != nil {
// 			return err
// 		}
// 	}
// 	// otherwise just update it
// 	gif.ID = current.ID
// 	err = c.Update(gif)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func (c *DBClient) LatestGIF() (*GifHistory, error) {
	gif := new(GifHistory)
	if result := c.db.Model(&GifHistory{}).Last(gif); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, ErrRecordNotFound
		}
		return nil, ErrDatabaseGeneral(result.Error.Error())
	}
	return gif, nil
}

// Close wraps the db close function for easy cleanup
func (c *DBClient) Close() {
	c.db.Close()
}
