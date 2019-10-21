package gif

import (
	"errors"
	"fmt"

	"github.com/king-jam/gotd/pkg/api/models"
)

// Repository represent the gif usecases
type Repository interface {
	Insert(gif *models.GIF) error
	Update(gif *models.GIF) error
	Last() (*models.GIF, error)
}

var (
	// ErrRecordNotFound happens when we haven't found any matched data
	ErrRecordNotFound = errors.New("record not found")
)

// ErrDatabaseGeneral is a generic error wrapper for unexplained errors
type ErrDatabaseGeneral string

func (edg ErrDatabaseGeneral) Error() string {
	return fmt.Sprintf("General Database Error: %s", edg)
}
