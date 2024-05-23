package server

import (
	"books/internal/models"
	"books/web/pages"

	"github.com/labstack/echo/v4"
)

func homePageHandler(c echo.Context) error {
	user := c.Get("user").(models.User)
	component := pages.HomePage(user)
	return component.Render(c.Request().Context(), c.Response())
}
