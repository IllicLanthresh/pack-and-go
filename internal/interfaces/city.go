package interfaces

import (
	"context"

	"github.com/IllicLanthresh/pack-and-go/internal/models"
)

type CityRepo interface {
	ReadById(ctx context.Context, id int64) (*models.City, error)
	ReadByName(ctx context.Context, name string) (*models.City, error)
}
