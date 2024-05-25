package tests

import (
	"books/internal/server"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
)

func TestHealthHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	router := server.NewRouter(server.WithTestDB())
	// Assertions
	if err := router.HealthHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}

}

func TestAddBookHandler(t *testing.T) {
	////////////////////////////////////////
	// Setup database
	////////////////////////////////////////

	db, err := setupDB()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/addBook", nil)
	req.Form = map[string][]string{
		"title":     {"test_title"},
		"author":    {"test_author"},
		"editorial": {"test_editorial"},
		"isbn":      {"test_isbn"},
		"review":    {"test_review"},
	}
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	c.Set("user", users[0])

	//set form values
	router := server.NewRouter(server.WithTestDB())

	////////////////////////////////////////
	// Assertions
	////////////////////////////////////////
	if err := router.AddBookHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusFound {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}

	//check if the book was added
	books, err := (*db).GetBookByID(1)
	if err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if books.ID != 1 {
		t.Error("Book was not added")
		return
	}
}
