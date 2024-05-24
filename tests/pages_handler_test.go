package tests

import (
	"books/internal/database"
	"books/internal/models"
	"books/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

func TestHomePage(t *testing.T) {
	////////////////////////////////////////
	/////Setup database
	////////////////////////////////////////
	db := database.New(database.TestConfig(), database.ResetDatabase())
	err := db.Connect()
	if err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	//add two test users
	testUser1 := models.User{
		Name:  "test_name",
		Email: "test_mail@mail.test",
	}
	testUser1.ID, err = db.AddUserIfNotExists(testUser1.Name, testUser1.Email, "test_image")
	testUser2 := models.User{
		Name:  "test_name2",
		Email: "test_mail2@mail.test",
	}
	testUser2.ID, err = db.AddUserIfNotExists(testUser2.Name, testUser2.Email, "test_image2")

	if err != nil {
		t.Error("Error adding test user 2")
		return
	}

	//add books
	books := []models.Book{
		{
			Title:     "test_title",
			Author:    "test_autor",
			Editorial: "test_editorial",
		},
		{
			Title:     "test_title2",
			Author:    "test_autor2",
			Editorial: "test_editorial2",
		},
		{
			Title:     "test_title3",
			Author:    "test_autor3",
			Editorial: "test_editorial3",
		},
	}

	for i, book := range books {
		if i%2 == 0 {
			book.ID, err = db.AddBook(book, testUser2.ID, "")

		} else {
			book.ID, err = db.AddBook(book, testUser1.ID, "")
		}
		if err != nil {
			t.Error("Error adding test book")
			return
		}
	}
	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	c.Set("user", testUser2)
	server := server.NewRouter(server.WithTestDB())

	////////////////////////////////////////
	// Assertions
	////////////////////////////////////////
	if err := server.HomePageHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	booksHTML := doc.Find(`[data-testid="BooksList"] a`)
	if booksHTML.Length() != 2 {
		t.Errorf("handler() wrong number of books = %v", booksHTML.Length())
	}
	booksHTML.Each(func(i int, s *goquery.Selection) {
		if i == 0 {
			if s.Text() != books[2].Title {
				t.Errorf("handler() wrong title = %v", s.Text())
			}
		} else if i == 1 {
			if s.Text() != books[0].Title {
				t.Errorf("handler() wrong title = %v", s.Text())
			}
		}
	})
}
