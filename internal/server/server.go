package server

import (
	"github.com/gokuls-codes/on-the-go/internal/services/docker"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port string
}

func NewServer(port string) *Server {
	return &Server{port: port}
}

func (s *Server) Start() {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())
	
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Static("/static", "static")

	e.GET("/", func(c echo.Context) error {
		return utils.Render(c, templates.Base())
	})

	e.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	dashboardGroup := e.Group("/dashboard")
	dashboardGroup.GET("", func(c echo.Context) error {
		return utils.Render(c, pages.Dashboard())
	})

	dockerGroup := dashboardGroup.Group("/docker")

	dockerHandler := docker.Handler{}
	dockerHandler.RegisterRoutes(dockerGroup)

	e.Logger.Fatal(e.Start(":" + s.port))
}