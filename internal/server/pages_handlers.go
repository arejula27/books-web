package server

import (
	"books/internal/models"
	"books/web/pages"

	"github.com/labstack/echo/v4"
)

func (s *Router) HomePageHandler(c echo.Context) error {
	user := c.Get("user").(models.User)
	books, err := s.db.GetBooksFromUser(user.ID)
	if err != nil {
		return err
	}
	component := pages.HomePage(user, books)
	return component.Render(c.Request().Context(), c.Response())
}
