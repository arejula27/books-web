package database

import (
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
	getUser    = "SELECT name  FROM users WHERE email = $1"
)

func (s *service) createTable() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, createTableUsers)

	return err
}

func (s *service) insertUser(email, name string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err := s.db.ExecContext(ctx, insertUser, email, name)

	return err
}

func (s *service) getUserByEmail(email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var name string
	err := s.db.QueryRowContext(ctx, getUser, email).Scan(&name)

	return name, err
}
