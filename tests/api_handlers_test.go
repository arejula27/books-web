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

	conn, err := connectDB()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}

	type testCases struct {
		testName        string
		forms           []map[string][]string
		expectedCode    int
		expedtedBooks   int
		expectedVolumes int
	}

	tests := []testCases{
		testCases{
			testName: "Add new book and new volume",
			forms: []map[string][]string{
				map[string][]string{
					"title":     {"test_title"},
					"author":    {"test_author"},
					"editorial": {"test_editorial"},
					"isbn":      {"new_test_isbn"},
					"review":    {"test_review"},
					tags[0]:     {"true"},
				},
			},
			expectedCode:    http.StatusFound,
			expedtedBooks:   len(books) + 1,
			expectedVolumes: len(books) + 1,
		},
		testCases{
			testName: "Add book of an existing volume",
			forms: []map[string][]string{
				map[string][]string{
					"title":     {"test_title"},
					"author":    {"test_author"},
					"editorial": {"test_editorial"},
					"isbn":      {"test_isbn"},
					"review":    {"test_review"},
					tags[0]:     {"true"},
				},
			},
			expectedCode:    http.StatusFound,
			expedtedBooks:   len(books) + 1,
			expectedVolumes: len(books),
		},
	}

	////////////////////////////////////////
	// Setup server
	////////////////////////////////////////
	e := echo.New()

	for _, testCase := range tests {

		t.Run(testCase.testName, func(t *testing.T) {
			for _, form := range testCase.forms {
				setupDB(conn)
				req := httptest.NewRequest(http.MethodPost, "/addBook", nil)
				req.Form = form
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
				if resp.Code != testCase.expectedCode {
					t.Errorf("handler() wrong status code = %v", resp.Code)
					return
				}

				//check if the book was added
				var count int
				err = conn.QueryRow("SELECT COUNT(*) FROM books").Scan(&count)
				if err != nil {
					t.Errorf("handler() error = %v", err)
					return
				}
				if count != testCase.expedtedBooks {
					t.Error("Book was not added")
					return
				}
				//check if a new volume was added
				err = conn.QueryRow("SELECT COUNT(*) FROM volumes").Scan(&count)
				if err != nil {
					t.Errorf("handler() error = %v", err)
					return
				}
				if count != testCase.expectedVolumes {
					t.Error("Volume was not added")
					return
				}
			}
		})
	}

}
