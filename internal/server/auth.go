package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/markbates/goth/gothic"
)

func googleCallback(c echo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		log.Println(err)
		return c.String(http.StatusInternalServerError, "Error")
	}
	//save the user to the session
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	err = gothic.StoreInSession("user", string(userJSON), c.Request(), c.Response())
	if err != nil {
		return err
	}
	return c.Redirect(http.StatusTemporaryRedirect, "/home")
}
func login(c echo.Context) error {

	_, err := gothic.GetFromSession("user", c.Request())
	//if the user is already logged in, redirect to the user page
	if err == nil {
		return c.Redirect(http.StatusTemporaryRedirect, "/home")
	}

	provider := c.Param("provider")
	res := c.Response()
	req := c.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.BeginAuthHandler(res, req)
	return nil
}
func logout(c echo.Context) error {
	provider := c.Param("provider")
	res := c.Response()
	req := c.Request().WithContext(context.WithValue(context.Background(), "provider", provider))
	gothic.Logout(res, req)
	return c.Redirect(http.StatusTemporaryRedirect, "/")
}
