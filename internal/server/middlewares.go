package server

import (
	"books/internal/models"
	"net/http"
	"strings"

	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// AuthorizationMiddleware is a middleware to check if the user is authorized
func (s *Router) authorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		//check the session is valid
		userJSON, err := gothic.GetFromSession("user", c.Request())
		if err != nil {
			return c.JSON(401, map[string]string{"message": "Unauthorized"})
		}
		//parse the user from the session
		var gothUser goth.User
		err = json.Unmarshal([]byte(userJSON), &gothUser)
		if err != nil {
			return c.JSON(401, map[string]string{"message": "Unauthorized"})
		}
		user := models.User{
			Email:    gothUser.Email,
			Name:     gothUser.Name,
			ImageURL: gothUser.AvatarURL,
		}
		if user.Name == "" {
			user.Name = strings.Split(user.Email, "@")[0]
		}
		//if is new user, add it to the database
		user.ID, err = s.db.AddUserIfNotExists(user.Email, user.Name, user.ImageURL)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		//set the user in the context
		c.Set("user", user)

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
