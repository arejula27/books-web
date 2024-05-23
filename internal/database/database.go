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
	_ "github.com/joho/godotenv/autoload"
)

// Service is an interface for the database service
type Service interface {
	Health() map[string]string
	AddUserIfNotExists(email, name string) (models.User, error)
}

type service struct {
	db *sql.DB
}

var (
	database = os.Getenv("DB_DATABASE")
	password = os.Getenv("DB_PASSWORD")
	username = os.Getenv("DB_USERNAME")
	port     = os.Getenv("DB_PORT")
	host     = os.Getenv("DB_HOST")
)

// New creates a new database service
func New() Service {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatal(err)
	}
	s := &service{db: db}
	s.createTable()
	return s
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

func (s *service) AddUserIfNotExists(email, name string) (models.User, error) {
	//check if user exists
	user, err := s.getUserByEmail(email)
	if nil != err {
		//if user does not exist insert user
		user, err = s.insertUser(email, name)
		if nil != err {
			return models.User{}, err
		}
	}
	//if user exists return no error
	return user, nil
}
