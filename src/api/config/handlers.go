package config

import (
	"github.com/labstack/echo"
	"net/http"
)

var (
	NotFoundHandler = func(c echo.Context) error {

		return c.String(http.StatusNotFound, "Not Found")
	}

	ForbiddenHandler = func(c echo.Context) error {

		return c.String(http.StatusForbidden, "Forbidden")
	}

	MaintenanceHandler = func(c echo.Context) error {
		c.Response().Header().Set("Retry-After", "3600") // retry after 1 hourse
		return c.String(http.StatusServiceUnavailable, "Service Unavailable")
	}

	InternalErrorHandler = func(c echo.Context) error {
		return c.String(http.StatusInternalServerError, "Internal Server Error")
	}
)
