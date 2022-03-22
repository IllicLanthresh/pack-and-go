package controllers

import (
	"net/http"
	"strconv"

	"github.com/IllicLanthresh/pack-and-go/internal/models"

	"github.com/IllicLanthresh/pack-and-go/internal/interfaces"
	"github.com/labstack/echo/v4"
)

type tripController struct {
	tripService interfaces.TripService
}

func NewTripController(tripService interfaces.TripService) *tripController {
	return &tripController{tripService: tripService}
}

type tripApiResponse struct {
	Id          int64           `json:"id"` // I assume this was actually missing in the README
	Origin      string          `json:"origin"`
	Destination string          `json:"destination"`
	Dates       models.Weekdays `json:"dates"`
	Price       float64         `json:"price"`
}

type tripApiRequest struct {
	OriginId      int64           `json:"originId"`
	DestinationId int64           `json:"destinationId"`
	Dates         models.Weekdays `json:"dates"`
	Price         float64         `json:"price"`
}

type createTripApiResponse struct {
	Id int64 `json:"id"`
}

func (c tripController) AddRoutes(e *echo.Echo) {
	e.GET("/trip", c.allTrips())
	e.POST("/trip", c.createTrip())
	e.GET("/trip/:id", c.singleTrip())
}

func (c tripController) allTrips() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		trips, err := c.tripService.GetAllTrips(ctx.Request().Context())
		if err != nil {
			return err
		}

		tripsResponse := make([]tripApiResponse, 0)
		for _, trip := range trips {
			tripsResponse = append(
				tripsResponse,
				tripApiResponse{
					Id:          trip.Id,
					Origin:      trip.Origin.Name,
					Destination: trip.Destination.Name,
					Dates:       trip.Dates,
					Price:       float64(trip.Price) / 100,
				},
			)
		}

		return ctx.JSON(http.StatusOK, tripsResponse)
	}
}

func (c tripController) createTrip() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		var request tripApiRequest
		err := ctx.Bind(&request)
		if err != nil {
			return err
		}

		id, err := c.tripService.CreateTrip(ctx.Request().Context(), &models.Trip{
			Origin: models.City{
				ID: request.OriginId,
			},
			Destination: models.City{
				ID: request.DestinationId,
			},
			Price: int64(request.Price * 100),
			Dates: request.Dates,
		})
		if err != nil {
			return err
		}

		return ctx.JSON(http.StatusOK, createTripApiResponse{
			Id: id,
		})
	}
}

func (c tripController) singleTrip() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil || id < 0 {
			return ctx.NoContent(http.StatusBadRequest)
		}

		trip, err := c.tripService.GetTripById(ctx.Request().Context(), id)
		if err != nil {
			return err
		}
		if trip == nil {
			return ctx.JSON(http.StatusNotFound, struct{}{})
		}

		return ctx.JSON(http.StatusOK, tripApiResponse{
			Id:          trip.Id,
			Origin:      trip.Origin.Name,
			Destination: trip.Destination.Name,
			Dates:       trip.Dates,
			Price:       float64(trip.Price) / 100,
		})
	}
}
