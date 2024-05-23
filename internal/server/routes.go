// Package server provides a server to handle the HTTP requests
package server

import (
	"net/http"

	"books/internal/models"
	"books/web"
	"books/web/pages"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes registers the routes for the server
func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//static file server
	jsFileServer := http.FileServer(http.FS(web.JS))
	e.GET("/js/*", echo.WrapHandler(jsFileServer))
	cssFileServer := http.FileServer(http.FS(web.CSS))
	e.GET("/css/*", echo.WrapHandler(cssFileServer))

	// Public pages routes
	e.GET("/", echo.WrapHandler(templ.Handler(pages.IndexPage())), redirectIfLogged)

	//TODO: Handle error if db is not connected
	//e.GET("/health", s.healthHandler)

	//Private pages routes
	pr := e.Group("", s.authorizationMiddleware)
	pr.GET("/home", homePageHandler)

	// API routes
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	// Auth routes
	e.GET("/auth/:provider/callback", googleCallback)
	e.GET("/logout/:provider", logout)
	e.GET("/auth/:provider", login)

	return e
}

func homePageHandler(c echo.Context) error {
	user := c.Get("user").(models.User)
	component := pages.HomePage(user)
	return component.Render(c.Request().Context(), c.Response())
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
