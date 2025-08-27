package server

import (
	"github.com/gokuls-codes/on-the-go/internal/services/docker"
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

	e.GET("/hello", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	dockerGroup := e.Group("/docker")

	dockerHandler := docker.Handler{}
	dockerHandler.RegisterRoutes(dockerGroup)

	e.Logger.Fatal(e.Start(":" + s.port))
}