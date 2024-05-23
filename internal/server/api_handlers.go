package server

import (
	"books/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
func (s *Server) addBookHandler(c echo.Context) error {
	//get the book from the form

	var book models.Book
	book.Title = c.FormValue("title")
	book.Author = c.FormValue("author")
	book.Editorial = c.FormValue("editorial")

	book, err := s.db.AddBook(book)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/home")
}
