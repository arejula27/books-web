package server

import (
	"books/internal/models"
	"books/web/pages"
	"strconv"

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

func (s *Router) BookPageHandler(c echo.Context) error {
	id := c.Param("id")
	bookID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	//Retrieve book from database
	book, err := s.db.GetBookByID(bookID)
	if err != nil {
		return err
	}
	//Retrieve reviews from database
	reviews, err := s.db.GetReviewsByBookID(bookID)
	props := pages.TimelinePageProps{
		Book:     book,
		Timeline: reviews,
	}

	component := pages.TimelinePage(props)
	return component.Render(c.Request().Context(), c.Response())
}
