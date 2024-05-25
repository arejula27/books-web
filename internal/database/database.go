// Package database provides a service to connect and interact with the database
package database

import (
	"books/internal/models"
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

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
}

type service struct {
	db      *sql.DB
	connStr string
	reset   string
	appEnv  string
	verbose bool
}

var ()

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

func (s *service) logMsg(msg string) {
	if s.verbose {
		log.Println(msg)
	}
}

func (s *service) Connect() error {
	db, err := sql.Open("pgx", s.connStr)
	if err != nil {
		return err
	}
	s.db = db

	if s.reset == "true" && s.appEnv != "production" {
		s.logMsg("resetting database")
		s.clearDatabase()
	}
	s.createTables()
	return nil
}

// Health checks the health of the database, returns a map with a message or stops the application if the database is down
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

// TODO: change parameter to models.User
func (s *service) AddUserIfNotExists(email, name, imageURL string) (int, error) {
	//check if user exists
	user, err := s.getUserByEmail(email)

	if nil == err {
		//if user exists return no error
		return user.ID, nil
	} else if sql.ErrNoRows == err {
		//if user does not exist insert user
		userID, err := s.insertUser(email, name, imageURL)
		if nil != err {
			return 0, err
		}
		return userID, nil
	}
	//return error
	return 0, err
}
func (s *service) AddBook(book models.Book, userID int, review string) (int, error) {
	//insert book
	bookID, err := s.insertBook(book, userID, review)
	if nil != err {
		return 0, err
	}
	//return book
	return bookID, nil

}
func (s *service) GetBooksFromUser(userID int) ([]models.Book, error) {
	//get books from user
	books, err := s.selectBooksFromUser(userID)
	if nil != err {
		return nil, err
	}

	return books, nil
}

func (s *service) GetBookByID(bookID int) (models.Book, error) {
	book, err := s.selectBookByID(bookID)
	if nil != err {
		return models.Book{}, err
	}
	return book, nil
}

func (s *service) GetReviewsByBookID(bookID int) ([]models.Review, error) {
	reviews, err := s.selectReviewsByBookID(bookID)
	if nil != err {
		return nil, err
	}
	return reviews, nil
}
