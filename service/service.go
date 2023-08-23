package service

import (
	"net/http"
	"time"

	cache "github.com/SporkHubr/echo-http-cache"
	"github.com/SporkHubr/echo-http-cache/adapter/memory"
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

	// Middleware Cache settings
	cacheCapacity   = 10000
	cacheRefreshKey = "opn" // ?$cacheRefreshKey=true to a page to force a refresh
	cacheTTL        = 32 * time.Hour
	nonCachedPaths  = []string{"/api"}
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
			return c.Request().URL.Path == "/metrics"
		},
	}))

	// Let's setup the in memory cache as middleware.
	if s.conf.MiddlewareHTMLCache {
		memcached, err := memory.NewAdapter(
			memory.AdapterWithAlgorithm(memory.LRU),
			memory.AdapterWithCapacity(cacheCapacity),
		)
		if err != nil {
			return e, err
		}
		cacheClient, err := cache.NewClient(
			cache.ClientWithAdapter(memcached),
			cache.ClientWithTTL(cacheTTL),
			cache.ClientWithRefreshKey(cacheRefreshKey),
			cache.ClientWithRestrictedPaths(nonCachedPaths),
		)
		if err != nil {
			return e, err
		}
		e.Use(cacheClient.Middleware())
	}

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
