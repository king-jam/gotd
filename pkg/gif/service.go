package gif

import (
	"context"

	"github.com/king-jam/gotd/pkg/api/models"
)

// Service represent the gif usecases
type Service interface {
	Set(ctx context.Context, gif *models.GIF) error
	Latest(ctx context.Context) (*models.GIF, error)
}
