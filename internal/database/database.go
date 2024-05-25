// Package database provides a service to connect and interact with the database
package database

import (
	"books/internal/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // initialize pgx driver
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload" // load .env file
)

// Service is an interface for the database service
type Service interface {
	Connect() error
	Health() map[string]string
	AddUserIfNotExists(email, name, imageURL string) (int, error)
	AddBook(book models.Book, userID int, review string) (int, error)
	GetBooksFromUser(userID int) ([]models.Book, error)
	GetBookByID(bookID int) (models.Book, error)
	GetReviewsByBookID(bookID int) ([]models.Review, error)
	GetAllTags() ([]string, error)
}

type service struct {
	db      *sql.DB
	connStr string
	reset   string
	appEnv  string
	verbose bool
}

type OptFunc func(*opts)

type opts struct {
	database string
	password string
	username string
	port     string
	host     string
	reset    string
	appEnv   string
	verbose  bool
}

func defaultOpts() *opts {

	return &opts{
		database: os.Getenv("DB_DATABASE"),
		password: os.Getenv("DB_PASSWORD"),
		username: os.Getenv("DB_USERNAME"),
		verbose:  os.Getenv("VERBOSE") == "true",
		port:     os.Getenv("DB_PORT"),
		host:     os.Getenv("DB_HOST"),
		reset:    os.Getenv("DB_RESET"),
		appEnv:   os.Getenv("APP_ENV")}
}
func TestConfig() OptFunc {
	return func(o *opts) {
		//print the current directory
		err := godotenv.Load("../.env.test")
		if err != nil {
			log.Println("Error loading .env file")

		}

		o.database = "books_test"
		if os.Getenv("DB_DATABASE") != "" {
			o.database = os.Getenv("DB_DATABASE")
		}
		o.appEnv = "test"
		o.username = os.Getenv("DB_USERNAME")
		o.password = os.Getenv("DB_PASSWORD")
		o.host = os.Getenv("DB_HOST")
		o.port = os.Getenv("DB_PORT")
		o.verbose = os.Getenv("VERBOSE") == "true"

	}
}
func ResetDatabase() OptFunc {
	return func(o *opts) {
		o.reset = "true"
	}
}

// New creates a new database service
func New(opts ...OptFunc) Service {
	opt := defaultOpts()
	for _, fn := range opts {
		fn(opt)
	}
	s := &service{
		connStr: fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", opt.username, opt.password, opt.host, opt.port, opt.database),
		reset:   opt.reset,
		appEnv:  opt.appEnv,
		verbose: opt.verbose,
	}

	return s
}
