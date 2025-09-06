package server

import (
	"github.com/gokuls-codes/on-the-go/internal/db"
	"github.com/gokuls-codes/on-the-go/internal/services/docker"
	"github.com/gokuls-codes/on-the-go/internal/services/git"
	"github.com/gokuls-codes/on-the-go/internal/services/system"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port  string
	store *db.Store
}

func NewServer(port string, store *db.Store) *Server {
	return &Server{port: port, store: store}
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
		return utils.Render(c, pages.DashboardPage())
	})

	dockerHandler := docker.Handler{
		Store: s.store,
	}
	dockerHandler.RegisterRoutes(dashboardGroup)

	gitGroup := e.Group("/git")
	gitHandler := git.Handler{
		Store: s.store,
	}
	gitHandler.RegisterRoutes(gitGroup)

	systemGroup := dashboardGroup.Group("/system")
	systemHandler := system.Handler{}
	systemHandler.RegisterRoutes(systemGroup)

	e.Logger.Fatal(e.Start(":" + s.port))
}
