// Package server provides a server to handle the HTTP requests
package server

import (
	"net/http"

	"books/web"
	"books/web/pages"

	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/markbates/goth/gothic"
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
	e.GET("/", echo.WrapHandler(templ.Handler(pages.IndexPage())), func(next echo.HandlerFunc) echo.HandlerFunc {
		//middleware to check if the user is already logged in and redirect to home page
		return func(c echo.Context) error {
			_, err := gothic.GetFromSession("user", c.Request())
			if err == nil {
				return c.Redirect(http.StatusFound, "/home")
			}
			//if the user is not logged in, continue to the index page
			return next(c)
		}
	})

	//TODO: Handle error if db is not connected
	//e.GET("/health", s.healthHandler)

	//Private pages routes
	pr := e.Group("", AuthorizationMiddleware)
	pr.GET("/home", echo.WrapHandler(templ.Handler(pages.HomePage())))

	// API routes
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	// Auth routes
	e.GET("/auth/:provider/callback", googleCallback)
	e.GET("/logout/:provider", logout)
	e.GET("/auth/:provider", login)

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
