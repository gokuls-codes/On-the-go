package server

import (
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type Server struct {
	port string
	conn *pgx.Conn
}

func NewServer(port string, conn *pgx.Conn) *Server {
	return &Server{port: port, conn: conn}
}

func (s *Server) Start() {
	e := echo.New()
	
	e.GET("/", func(c echo.Context) error {
		return c.String(200, "Hello, World!")
	})

	e.Logger.Fatal(e.Start(":" + s.port))
}