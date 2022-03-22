package middlewares

import (
	"time"

	"github.com/labstack/echo/v4/middleware"
)

func init() {
	middlewares[postRoutingMiddleware] = append(middlewares[postRoutingMiddleware],
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{
			Timeout: 10 * time.Second,
		}),
	)
}
