package server

import (
	"encoding/json"
	"github.com/labstack/echo/v4"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

// AuthorizationMiddleware is a middleware to check if the user is authorized
func AuthorizationMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		//check the session is valid
		userJSON, err := gothic.GetFromSession("user", c.Request())
		if err != nil {
			return c.JSON(401, map[string]string{"message": "Unauthorized"})
		}
		var user goth.User
		err = json.Unmarshal([]byte(userJSON), &user)

		if err != nil {
			return c.JSON(401, map[string]string{"message": "Unauthorized"})
		}
		c.Set("user", user)
		return next(c)
	}

}
