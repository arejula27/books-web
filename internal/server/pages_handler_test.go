package server_test

import (
	"books/internal/server"
	"books/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/labstack/echo/v4"
)

func init() {
	utils.LoadEnv()
}

func TestHomePage(t *testing.T) {
	////////////////////////////////////////
	/////Setup database
	////////////////////////////////////////
	conn, err := utils.ConnectDB()
	testData := utils.TestScenario1Data()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	err = utils.SetupDB(conn, &testData)
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
	c.Set("user", testData.Users[0])
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
	if booksHTML.Length() != len(testData.Books) {
		t.Errorf("handler() wrong number of books = %v", booksHTML.Length())
	}
	booksHTML.Each(func(i int, s *goquery.Selection) {

		if s.Text() != testData.Books[len(testData.Books)-i-1].Title {
			t.Errorf("handler() wrong title = %v", s.Text())
		}
	})
}

func TestAddBookPage(t *testing.T) {
	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	conn, err := utils.ConnectDB()
	testData := utils.TestScenario1Data()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	err = utils.SetupDB(conn, &testData)
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/addBook", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	server := server.NewRouter()

	////////////////////////////////////////
	// Assertions
	////////////////////////////////////////
	if err := server.AddBookPageHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
}
func TestBookPage(t *testing.T) {
	////////////////////////////////////////
	/////Setup database
	////////////////////////////////////////
	conn, err := utils.ConnectDB()
	testData := utils.TestScenario1Data()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	err = utils.SetupDB(conn, &testData)
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
	c.Set("user", testData.Users[0])
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
	conn, err := utils.ConnectDB()
	testData := utils.TestScenario1Data()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	err = utils.SetupDB(conn, &testData)
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
	c.Set("user", testData.Users[0])
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
