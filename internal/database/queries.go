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
		image_url VARCHAR(255),
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	insertUser = `INSERT INTO USERS (email, name) VALUES ($1, $2)`
	getUser    = "SELECT id, name, email,image_url  FROM users WHERE email = $1"
)

func (s *service) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, createTableUsers)

	return err
}

func (s *service) insertUser(email, name, imageURL string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	query := `INSERT INTO users (email, name, image_url) VALUES ($1, $2,$3) RETURNING id, email, name, image_url`
	var user models.User
	err := s.db.QueryRowContext(ctx, query, email, name, imageURL).Scan(&user.ID, &user.Email, &user.Name, &user.ImageURL)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (s *service) getUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var user models.User
	err := s.db.QueryRowContext(ctx, getUser, email).Scan(&user.ID, &user.Name, &user.Email, &user.ImageURL)

	return user, err
}
