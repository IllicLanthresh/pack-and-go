package interfaces

import (
	"context"

	"github.com/IllicLanthresh/pack-and-go/internal/models"
)

type TripRepo interface {
	ListAll(ctx context.Context) ([]*models.Trip, error)
	Create(ctx context.Context, trip *models.Trip) (int64, error)
	ReadById(ctx context.Context, id int64) (*models.Trip, error)
}
