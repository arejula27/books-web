package tests

import (
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
	_, err := setupDB()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	c.Set("user", users[0])
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

func TestBookPage(t *testing.T) {
	////////////////////////////////////////
	/////Setup database
	////////////////////////////////////////
	_, err := setupDB()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/book/1", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	c.SetParamNames("id")
	c.SetParamValues("1")
	c.Set("user", users[0])
	server := server.NewRouter(server.WithTestDB())

	////////////////////////////////////////
	// Assertions
	////////////////////////////////////////
	if err := server.BookPageHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
}
func TestGetBookPage(t *testing.T) {
	////////////////////////////////////////
	/////Setup database
	////////////////////////////////////////
	_, err := setupDB()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/getBook", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	c.Set("user", users[0])
	server := server.NewRouter(server.WithTestDB())

	////////////////////////////////////////
	// Assertions
	////////////////////////////////////////
	if err := server.GetBookHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
}
