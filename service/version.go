package service

import (
	"net/http"

	"github.com/darron/gips/config"
	"github.com/labstack/echo/v4"
)

func (s HTTPService) Version(c echo.Context) error {
	// Get all records.
	versionInfo := config.GetVersionInfo()
	return c.JSON(http.StatusOK, versionInfo)
}
