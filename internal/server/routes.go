// Package server provides a server to handle the HTTP requests
package server

import (
	"net/http"

	"books/web"
	"books/web/pages"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RegisterRoutes registers the routes for the server
func (s *Router) RegisterRoutes() http.Handler {
	e := echo.New()

	// Middlewares
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//static file server
	jsFileServer := http.FileServer(http.FS(web.JS))
	e.GET("/js/*", echo.WrapHandler(jsFileServer))
	cssFileServer := http.FileServer(http.FS(web.CSS))
	e.GET("/css/*", echo.WrapHandler(cssFileServer))
	fontsFileServer := http.FileServer(http.FS(web.Fonts))
	e.GET("/fonts/*", echo.WrapHandler(fontsFileServer))
	// Public pages routes
	e.GET("/", echo.WrapHandler(templ.Handler(pages.IndexPage())), redirectIfLogged)

	//TODO: Handle error if db is not connected
	//e.GET("/health", s.healthHandler)

	//Private pages routes
	pr := e.Group("", s.authorizationMiddleware)
	pr.GET("/home", s.HomePageHandler)
	pr.GET("/addBook", s.AddBookPageHandler)
	pr.GET("/book/:id", s.BookPageHandler)
	pr.GET("/getBook", s.GetBookHandler)

	// API routes
	pr.POST("/addBook", s.AddBookHandler)

	// Auth routes
	e.GET("/auth/:provider/callback", googleCallback)
	e.GET("/logout/:provider", logout)
	e.GET("/auth/:provider", login)

	return e
}
