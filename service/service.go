package service

import (
	"net/http"

	"github.com/darron/gips/config"
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HTTPService struct {
	conf *config.App
}

var (
	APIV1Path        = "/api/v1"
	IPsPath          = "/search/:ip"
	IPsPathFull      = APIV1Path + IPsPath
	ProjectPath      = "/project/:project"
	ProjectPathFull  = APIV1Path + ProjectPath
	ProjectsPath     = "/projects"
	ProjectsPathFull = APIV1Path + ProjectsPath
)

func Get(conf *config.App) (*echo.Echo, error) {
	s := HTTPService{}

	s.conf = conf

	// Echo instance
	e := echo.New()

	// Enable Prometheus
	p := prometheus.NewPrometheus("echo", nil)
	p.Use(e)

	// Let's allow some static files
	e.Static("/", "public")

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Skipper: func(c echo.Context) bool {
			if c.Request().URL.Path == "/metrics" || c.Request().URL.Path == "/version" {
				return true
			}
			return false
		},
	}))

	// Routes
	e.GET("/", s.Root)
	e.GET(IPsPathFull, s.IPSearch)
	e.GET(ProjectPathFull, s.Project)
	e.GET(ProjectsPathFull, s.Projects)

	// Infra
	e.GET("/version", s.Version)

	return e, nil
}

func (s HTTPService) Root(c echo.Context) error {
	return c.JSON(http.StatusOK, "hello")
}
