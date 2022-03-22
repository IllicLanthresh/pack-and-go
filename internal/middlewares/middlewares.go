package middlewares

import (
	"errors"

	"github.com/labstack/echo/v4"
)

const (
	preRoutingMiddleware = 1 << iota
	postRoutingMiddleware
)

func getMethodFor(e *echo.Echo, middlewareType middlewareType) (func(...echo.MiddlewareFunc), error) {
	switch middlewareType {
	case preRoutingMiddleware:
		return e.Pre, nil
	case postRoutingMiddleware:
		return e.Use, nil
	default:
		return nil, errors.New("unsupported middleware type")
	}
}

func Init(e *echo.Echo) error {
	for t, ms := range middlewares {
		method, err := getMethodFor(e, t)
		if err != nil {
			return err
		}
		for _, m := range ms {
			method(m)
		}
	}
	return nil
}

type middlewareType int

var middlewares = make(map[middlewareType][]echo.MiddlewareFunc)
