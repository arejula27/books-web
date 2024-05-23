// Package server provides a server to handle the HTTP requests
package server

import (
	"fmt"
	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload" // load .env file
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"net/http"
	"os"
	"strconv"
	"time"

	"books/internal/database"
)

const (
	key    = "secret"
	MaxAge = 86400 * 30 // MaxAge is 30 days
	IsProd = false      // IsProd this to true when deploying to production
)

// Server is a struct that holds the server configuration
type Server struct {
	port int

	db database.Service
}

// NewServer creates a new server instance
func NewServer() *http.Server {

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,

		db: database.New(),
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	store := sessions.NewCookieStore([]byte(key))

	store.MaxAge(MaxAge)
	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = IsProd

	gothic.Store = store

	goth.UseProviders(google.New(googleClientID, googleClientSecret, "http://localhost:3000/auth/google/callback"))

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
