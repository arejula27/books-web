package server

import (
	"books/internal/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Router) HealthHandler(c echo.Context) error {
	err := s.db.Health()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, nil)
}
func (s *Router) AddBookHandler(c echo.Context) error {
	//get the book from the form

	var book models.Book
	var review string
	book.Tags = []string{}
	userID := c.Get("user").(models.User).ID
	params, err := c.FormParams()
	if err != nil {
		return err
	}
	for param, value := range params {
		if param == "title" {
			book.Title = value[0]
		} else if param == "author" {
			book.Author = value[0]
		} else if param == "editorial" {
			book.Editorial = value[0]
		} else if param == "isbn" {
			book.ISBN = value[0]
		} else if param == "review" {
			review = value[0]
		} else {
			book.Tags = append(book.Tags, param)

		}

	}

	_, err = s.db.AddBook(book, userID, review)
	if err != nil {
		return err
	}

	return c.Redirect(http.StatusFound, "/home")
}
