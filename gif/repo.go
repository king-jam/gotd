package gif

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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

type dbGIF struct {
	gorm.Model
	DeactivatedAt time.Time
	GIF           string
	RequesterID   string
	RequestSrc    string
	Tags          pq.StringArray `gorm:"type:varchar(64)[]"`
}

type Repo struct {
	DB *gorm.DB
}

func NewGIFRepo(orm *gorm.DB) (*Repo, error) {
	return &Repo{
		DB: orm,
	}, nil
}

func (r *Repo) InitDB() error {
	if !r.DB.HasTable(&dbGIF{}) {
		r.DB.CreateTable(&dbGIF{})
	}
	return nil
}

// Insert will add a gif into the database
func (r *Repo) Insert(gif *GIF) error {
	gotd := TransformGif(gif)
	if result := r.DB.Create(&gotd); result.Error != nil {
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (r *Repo) DeleteGIFByID(id int) error {
	return nil
}

func (r *Repo) Update(gif *GIF) error {
	gotd := TransformGif(gif)
	if result := r.DB.Model(&dbGIF{}).Updates(gotd); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return ErrRecordNotFound
		}
		return ErrDatabaseGeneral(result.Error.Error())
	}
	return nil
}

func (r *Repo) FindGIFByID(id uint) (*dbGIF, error) {
	return &dbGIF{}, nil
}

func (r *Repo) FindAllGifs() ([]dbGIF, error) {
	return []dbGIF{}, nil
}

func (r *Repo) LatestGIF() (*dbGIF, error) {
	gif := new(dbGIF)
	if result := r.DB.Model(&dbGIF{}).Last(gif); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, ErrRecordNotFound
		}
		return nil, ErrDatabaseGeneral(result.Error.Error())
	}
	return gif, nil
}

func TransformGif(gif *GIF) dbGIF {
	dbGif := dbGIF{
		Model: gorm.Model{
			ID:        gif.ID,
			CreatedAt: gif.CreatedAt,
			UpdatedAt: gif.UpdatedAt,
			DeletedAt: gif.DeletedAt,
		},
		GIF:           gif.GIF,
		RequestSrc:    gif.RequestSrc,
		RequesterID:   gif.RequesterID,
		Tags:          pq.StringArray(gif.Tags),
		DeactivatedAt: gif.DeactivatedAt,
	}
	return dbGif
}

// 1. Import SQL dialect we are using (postgres)
// 2. create GORM DB instance with dialect
// 3. create GIF Service Repo with GORM DB instance
// 4. create GIF Service with GIF Service Repo
