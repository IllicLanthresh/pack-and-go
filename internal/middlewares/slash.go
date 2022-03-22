package middlewares

import (
	"github.com/labstack/echo/v4/middleware"
)

func init() {
	middlewares[preRoutingMiddleware] = append(middlewares[preRoutingMiddleware],
		middleware.RemoveTrailingSlash(),
	)
}
