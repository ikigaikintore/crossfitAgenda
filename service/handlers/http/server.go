package http

import (
	"github.com/labstack/echo/v4/middleware"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func NewServer(handler CrossfitHandler) *echo.Echo {
	srv := echo.New()
	srv.HideBanner = true
	srv.Logger.SetLevel(log.DEBUG)
	srv.GET("/health", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	srv.Use(middleware.CORS())
	srv.Use(middleware.RequestID())
	srv.Use(middleware.Recover())

	RegisterHandlers(srv, handler)

	return srv
}
