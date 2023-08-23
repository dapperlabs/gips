package service

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s HTTPService) Project(c echo.Context) error {
	project := c.Param("project")
	if project == "" {
		return c.JSON(http.StatusNotFound, "project must not be blank")
	}
	outputProject, err := s.conf.ProjectStore.Find(project)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSONPretty(http.StatusOK, *outputProject, "  ")
}
