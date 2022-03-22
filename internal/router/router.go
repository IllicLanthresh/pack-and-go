package router

import (
	"database/sql"
	"net/http"

	"github.com/IllicLanthresh/pack-and-go/internal/controllers"
	"github.com/IllicLanthresh/pack-and-go/internal/interfaces"
	"github.com/IllicLanthresh/pack-and-go/internal/middlewares"
	"github.com/IllicLanthresh/pack-and-go/internal/repositories"
	"github.com/IllicLanthresh/pack-and-go/internal/services"
	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

func errorHandler(err error, ctx echo.Context) {
	ctx.Logger().Error(err)

	if httpError, ok := err.(*echo.HTTPError); ok {
		code := httpError.Code

		if err = ctx.String(code, httpError.Error()); err != nil {
			ctx.Logger().Error(err)
		}
	} else {
		code := http.StatusInternalServerError

		if err = ctx.NoContent(code); err != nil {
			ctx.Logger().Error(err)
		}
	}
}

func New() (*echo.Echo, error) {
	e := echo.New()
	e.HTTPErrorHandler = errorHandler
	e.HideBanner = true

	// Load and setup middlewares
	err := middlewares.Init(e)
	if err != nil {
		return nil, err
	}

	// Data sources
	db, err := sql.Open("sqlite3", "./db/data.db")
	if err != nil {
		return nil, err
	}

	// Repos
	var (
		tripDbRepo   = repositories.NewTripDbRepo(db)
		cityFileRepo = repositories.NewCityFileRepo("./db/cities.txt")
	)

	// Services
	var (
		tripService = services.NewTripService(tripDbRepo, cityFileRepo)
	)

	// Controllers
	for _, controller := range []interfaces.Controller{
		controllers.NewTripController(tripService),
	} {
		controller.AddRoutes(e)
	}

	return e, nil
}
