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
	createTableBooks = `
		CREATE TABLE BOOKS (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(100) NOT NULL,
		editorial VARCHAR(100) NOT NULL,
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	insertUser = `INSERT INTO users (email, name, image_url) VALUES ($1, $2,$3) RETURNING id, email, name, image_url`
	getUser    = "SELECT id, name, email, image_url  FROM users WHERE email = $1"
	insertBook = `INSERT INTO books (title, author, editorial) VALUES ($1, $2, $3) RETURNING id`
)

func (s *service) clearDatabase() {
	s.db.Exec("DROP TABLE IF EXISTS users")
	s.db.Exec("DROP TABLE IF EXISTS books")
}

func (s *service) createTables() error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	_, err = s.db.ExecContext(ctx, createTableUsers)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = s.db.ExecContext(ctx, createTableBooks)
	if err != nil {
		tx.Rollback()
		return err
	}
	return tx.Commit()

}

func (s *service) insertUser(email, name, imageURL string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var user models.User
	err := s.db.QueryRowContext(ctx, insertUser, email, name, imageURL).Scan(&user.ID, &user.Email, &user.Name, &user.ImageURL)
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

func (s *service) insertBook(book models.Book) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	err := s.db.QueryRowContext(ctx, insertBook, book.Title, book.Author, book.Editorial).Scan(&book.ID)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}
