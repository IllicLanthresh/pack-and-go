package interfaces

import "github.com/labstack/echo/v4"

type Controller interface {
	AddRoutes(*echo.Echo)
}
