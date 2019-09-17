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

type DB interface {
	InitDB() error
	Insert(*GIF) error
	DeleteGIFByID(id int) error
	Update(*GIF) error
	FindGIFByID(id uint) (*GIF, error)
	FindAllGifs() ([]GIF, error)
	LatestGIF() (*GIF, error)
}

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

type Repo struct {
	DB *gorm.DB
}

func NewGIFRepo(orm *gorm.DB) (*Repo, error) {
	return &Repo{
		DB: orm,
	}, nil
}

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

func (r *Repo) DeleteGIFByID(id int) error {
	return nil
}

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

func (r *Repo) FindGIFByID(id uint) (*GIF, error) {
	return &GIF{}, nil
}

func (r *Repo) FindAllGifs() ([]GIF, error) {
	return []GIF{}, nil
}

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

// func TransformGifToDBGif(gif *GIF) GIF {
// 	dbGif := GIF{
// 		Model: gorm.Model{
// 			ID:        gif.ID,
// 			CreatedAt: gif.CreatedAt,
// 			UpdatedAt: gif.UpdatedAt,
// 			DeletedAt: gif.DeletedAt,
// 		},
// 		GIF:           gif.GIF,
// 		RequestSrc:    gif.RequestSrc,
// 		RequesterID:   gif.RequesterID,
// 		Tags:          pq.StringArray(gif.Tags),
// 		DeactivatedAt: gif.DeactivatedAt,
// 	}
// 	return dbGif
// }

// func TransformDBGifToGif(dbGIF *GIF) GIF {
// 	gif := GIF{
// 		ID:            dbGIF.ID,
// 		CreatedAt:     dbGIF.CreatedAt,
// 		UpdatedAt:     dbGIF.UpdatedAt,
// 		DeletedAt:     dbGIF.DeletedAt,
// 		GIF:           dbGIF.GIF,
// 		RequestSrc:    dbGIF.RequestSrc,
// 		RequesterID:   dbGIF.RequesterID,
// 		Tags:          dbGIF.Tags,
// 		DeactivatedAt: dbGIF.DeactivatedAt,
// 	}
// 	return gif
// }

// 1. Import SQL dialect we are using (postgres)
// 2. create GORM DB instance with dialect
// 3. create GIF Service Repo with GORM DB instance
// 4. create GIF Service with GIF Service Repo
