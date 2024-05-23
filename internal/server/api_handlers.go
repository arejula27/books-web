package server

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}
func (s *Server) addBookHandler(c echo.Context) error {
	return c.Redirect(http.StatusFound, "/home")
}
