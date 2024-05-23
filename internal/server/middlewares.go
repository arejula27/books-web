package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func (s *Server) createUserMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(User)
		err := s.db.AddUserIfNotExists(user.Email, user.Name)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		return next(c)
	}
}

func redirectIfLogged(next echo.HandlerFunc) echo.HandlerFunc {
	//middleware to check if the user is already logged in and redirect to home page
	return func(c echo.Context) error {
		_, err := gothic.GetFromSession("user", c.Request())
		if err == nil {
			return c.Redirect(http.StatusFound, "/home")
		}
		//if the user is not logged in, continue to the index page
		return next(c)
	}
}
