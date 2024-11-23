package server

import "github.com/labstack/echo/v4"

type Server struct {
	Port string
}

func NewServer(port string) *Server {
	return &Server{Port: port}
}

func (s *Server) Start() {
	e := echo.New()
	
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + s.Port))
}