package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"plex-query/controller"
)

type Server struct {
	echo *echo.Echo
	ctrl *controller.Controller
	port string
}

func New(p string, c *controller.Controller) *Server {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	return &Server{
		echo: e,
		ctrl: c,
		port: p,
	}
}

func (s *Server) Run() error {
	s.echo.GET("/media", s.ctrl.GetMediaByTag)
	return s.echo.Start(":" + s.port)
}
