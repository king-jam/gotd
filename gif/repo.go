package gif

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
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

// GIF is the model used for actions
type GIF struct {
	gorm.Model
	DeactivatedAt *time.Time
	URL           string `json:"url"`
	RequesterID   string
	RequestSrc    string
	//Tags          []Tag `gorm:"many2many:gif_tags;"`
}

// type Tag struct {
// 	gorm.Model
// 	Value string
// }

// Repo provides access to the database through an abstraction
// that allows for swapping out the datastore and not breaking
// application logic
type Repo struct {
	DB *gorm.DB
}

// NewGIFRepo initializes the repo with the ORM
func NewGIFRepo(orm *gorm.DB) (*Repo, error) {
	return &Repo{
		DB: orm,
	}, nil
}

// InitDB provides hooks to ensure tables and migrations are performed
func (r *Repo) InitDB() error {
	if !r.DB.HasTable(&GIF{}) {
		r.DB.CreateTable(&GIF{})
	}
	r.DB.AutoMigrate(&GIF{})
	// if !r.DB.HasTable(&Tag{}) {
	// 	r.DB.CreateTable(&Tag{})
	// }
	return nil
}

// Insert will add a gif into the database
func (r *Repo) Insert(gif *GIF) error {
	//gotd := TransformGifToDBGif(gif)
	if result := r.DB.Model(&GIF{}).Create(gif); result.Error != nil {
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

// Update performs updates to a record
func (r *Repo) Update(gif *GIF) error {
	//gotd := TransformGifToDBGif(gif)
	if result := r.DB.Model(&GIF{}).Updates(gif); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return ErrRecordNotFound
		}
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

// LatestGIF gets the latest entry from the database
func (r *Repo) LatestGIF() (*GIF, error) {
	gif := new(GIF)
	if result := r.DB.Model(&GIF{}).Last(gif); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, ErrRecordNotFound
		}
		return nil, ErrDatabaseGeneral(result.Error.Error())
	}
	return gif, nil
}