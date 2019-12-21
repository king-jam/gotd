// Package repository provides a postgres implementation of the GIF repo interface
package repository

import (
	"github.com/jinzhu/gorm"
	"github.com/king-jam/gotd/pkg/api/models"
	"github.com/king-jam/gotd/pkg/gif"
)

// Repo provides access to the database through an abstraction
// that allows for swapping out the datastore and not breaking
// application logic
type Repo struct {
	DB *gorm.DB
}

// New initializes the repo with the ORM
func New(orm *gorm.DB) (*Repo, error) {
	return &Repo{
		DB: orm,
	}, nil
}

// InitDB provides hooks to ensure tables and migrations are performed
func (r *Repo) InitDB() error {
	if !r.DB.HasTable(&models.GIF{}) {
		r.DB.CreateTable(&models.GIF{})
	}

	r.DB.AutoMigrate(&models.GIF{})

	return nil
}

// Insert will add a gif into the database
func (r *Repo) Insert(g *models.GIF) error {
	if result := r.DB.Model(&models.GIF{}).Create(g); result.Error != nil {
		return gif.ErrDatabaseGeneral(result.Error.Error())
	}

	return nil
}

// Update performs updates to a record
func (r *Repo) Update(g *models.GIF) error {
	if result := r.DB.Model(&models.GIF{}).Updates(g); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return gif.ErrRecordNotFound
		}

		return gif.ErrDatabaseGeneral(result.Error.Error())
	}

	return nil
}

// Last gets the latest entry from the database
func (r *Repo) Last() (*models.GIF, error) {
	g := new(models.GIF)
	if result := r.DB.Model(&models.GIF{}).Last(g); result.Error != nil {
		if gorm.IsRecordNotFoundError(result.Error) {
			return nil, gif.ErrRecordNotFound
		}

		return nil, gif.ErrDatabaseGeneral(result.Error.Error())
	}

	return g, nil
}
