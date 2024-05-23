package database

import (
	"books/internal/models"
	"context"
	"time"
)

const (
	createTableUsers = `
		CREATE TABLE USERS (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	insertUser = `INSERT INTO USERS (email, name) VALUES ($1, $2)`
	getUser    = "SELECT id, name, email  FROM users WHERE email = $1"
)

func (s *service) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, createTableUsers)

	return err
}

func (s *service) insertUser(email, name string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, name) VALUES ($1, $2) RETURNING id, email, name`
	var user models.User
	err := s.db.QueryRowContext(ctx, query, email, name).Scan(&user.ID, &user.Email, &user.Name)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *service) getUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var user models.User
	err := s.db.QueryRowContext(ctx, getUser, email).Scan(&user.ID, &user.Name, &user.Email)

	return user, err
}
