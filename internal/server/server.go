package server

import (
	"database/sql"

	"github.com/gokuls-codes/on-the-go/internal/services/docker"
	"github.com/gokuls-codes/on-the-go/internal/services/git"
	"github.com/gokuls-codes/on-the-go/internal/utils"
	"github.com/gokuls-codes/on-the-go/internal/web/templates"
	"github.com/gokuls-codes/on-the-go/internal/web/templates/pages"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	port string
	db *sql.DB
}

func NewServer(port string, db *sql.DB) *Server {
	return &Server{port: port, db: db}
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
	
	dockerHandler := docker.Handler{}
	dockerHandler.RegisterRoutes(dashboardGroup)

	gitGroup := e.Group("/git")
	gitHandler := git.Handler{}
	gitHandler.RegisterRoutes(gitGroup)

	e.Logger.Fatal(e.Start(":" + s.port))
}