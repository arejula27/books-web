// Package server provides a server to handle the HTTP requests
package server

import (
	"net/http"

	"books/web"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes registers the routes for the server
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	jsFileServer := http.FileServer(http.FS(web.Files))
	e.GET("/js/*", echo.WrapHandler(jsFileServer))
	cssFileServer := http.FileServer(http.FS(web.CSS))
	e.GET("/css/*", echo.WrapHandler(cssFileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	e.GET("/", s.HelloWorldHandler)

	e.GET("/health", s.healthHandler)

	return e
}

// HelloWorldHandler is a handler that returns a simple message, used for testing purposes
func (s *Server) HelloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
