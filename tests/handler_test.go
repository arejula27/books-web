package tests

import (
	"books/internal/models"
	"books/internal/server"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)
	c.Set("user", models.User{})
	// Assertions
	if err := server.HomePageHandler(c); err != nil {
		t.Errorf("handler() error = %v", err)
		return
	}
	if resp.Code != http.StatusOK {
		t.Errorf("handler() wrong status code = %v", resp.Code)
		return
	}
}
