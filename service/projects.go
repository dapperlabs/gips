package service

import (
	"net/http"

	"github.com/darron/gips/core"
	"github.com/labstack/echo/v4"
)

func (s HTTPService) Projects(c echo.Context) error {
	outputProjects, err := s.conf.ProjectStore.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Convert so we can display the JSON nicely.
	var drefProjects []core.Project
	for _, p := range outputProjects {
		drefProjects = append(drefProjects, *p)
	}
	return c.JSONPretty(http.StatusOK, drefProjects, "  ")
}
