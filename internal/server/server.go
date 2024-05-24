// Package server provides a server to handle the HTTP requests
package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/sessions"
	_ "github.com/joho/godotenv/autoload" // load .env file
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"

	"books/internal/database"
)

const (
	key    = "secret"
	MaxAge = 86400 * 30 // MaxAge is 30 days
	IsProd = false      // IsProd this to true when deploying to production
)

// Router is a struct that holds the server configuration
type Router struct {
	port int
	db   database.Service
}
type optFunc func(*opts)
type opts struct {
	db   database.Service
	port int
}

func defaultOpts() *opts {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 3000
	}
	return &opts{
		port: port,
		db:   database.New(),
	}

}
func WithTestDB() optFunc {
	return func(o *opts) {
		o.db = database.New(database.TestConfig())
	}
}

func NewRouter(opts ...optFunc) *Router {
	opt := defaultOpts()
	for _, fn := range opts {
		fn(opt)
	}
	err := opt.db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	return &Router{
		port: opt.port,
		db:   opt.db,
	}
}

// New creates a new server instance
func New() *http.Server {

	NewServer := NewRouter()

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
