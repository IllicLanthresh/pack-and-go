package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/IllicLanthresh/pack-and-go/internal/util"

	"github.com/IllicLanthresh/pack-and-go/internal/interfaces"
	"github.com/IllicLanthresh/pack-and-go/internal/models"
)

type tripService struct {
	tripRepo interfaces.TripRepo
	cityRepo interfaces.CityRepo
}

func NewTripService(tripRepo interfaces.TripRepo, cityRepo interfaces.CityRepo) *tripService {
	return &tripService{tripRepo: tripRepo, cityRepo: cityRepo}
}

const alreadyCompleteError = util.Error("city already complete")

func (s tripService) completeCityData(ctx context.Context, incompleteCity *models.City) error {
	if !incompleteCity.Incomplete {
		return alreadyCompleteError
	}

	completeCity, err := s.cityRepo.ReadById(ctx, incompleteCity.ID)
	*incompleteCity = *completeCity
	return err
}

func (s tripService) completeCitiesData(ctx context.Context, trip *models.Trip) error {
	err := s.completeCityData(ctx, &trip.Origin)
	if err != nil && err != alreadyCompleteError {
		return fmt.Errorf("error querying origin city: %w", err)
	}

	err = s.completeCityData(ctx, &trip.Destination)
	if err != nil && err != alreadyCompleteError {
		return fmt.Errorf("error querying destination city: %w", err)
	}
	return err
}

func (s tripService) GetAllTrips(ctx context.Context) ([]*models.Trip, error) {
	trips, err := s.tripRepo.ListAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error querying trip data: %w", err)
	}

	for _, trip := range trips {
		err = s.completeCitiesData(ctx, trip)
		if err != nil {
			return nil, err
		}
	}

	return trips, nil
}

func (s tripService) CreateTrip(ctx context.Context, trip *models.Trip) (int64, error) {
	origin, err := s.cityRepo.ReadById(ctx, trip.Origin.ID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("error getting origin city: %w", err))
	}

	destination, err := s.cityRepo.ReadById(ctx, trip.Destination.ID)
	if err != nil {
		return 0, echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("error getting destination city: %w", err))
	}

	trip.Origin = *origin
	trip.Destination = *destination

	id, err := s.tripRepo.Create(ctx, trip)
	if err != nil {
		return id, fmt.Errorf("error saving trip data: %w", err)
	}
	return id, nil
}

func (s tripService) GetTripById(ctx context.Context, id int64) (*models.Trip, error) {
	trip, err := s.tripRepo.ReadById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error querying trip data: %w", err)
	}

	if trip == nil {
		return nil, nil
	}

	err = s.completeCitiesData(ctx, trip)

	return trip, err
}
