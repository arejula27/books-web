package database

import (
	"books/internal/models"
	"context"
	"log"
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
	createTableReviews = `
	CREATE TABLE REVIEWS (
		id SERIAL PRIMARY KEY,
		book_id INT NOT NULL, 
		user_id INT NOT NULL,
		comment TEXT,
		creation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_book
			FOREIGN KEY(book_id) 
				REFERENCES BOOKS(id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
		CONSTRAINT fk_user
			FOREIGN KEY(user_id)
				REFERENCES USERS(id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	);
	`
	insertUser   = `INSERT INTO users (email, name, image_url) VALUES ($1, $2,$3) RETURNING id`
	getUser      = "SELECT id, name, email, image_url  FROM users WHERE email = $1"
	insertBook   = `INSERT INTO books (title, author, editorial) VALUES ($1, $2, $3) RETURNING id`
	insertReview = `INSERT INTO reviews (book_id, user_id, comment) VALUES ($1, $2, $3) RETURNING id`
)

func (s *service) clearDatabase() {
	_, err := s.db.Exec("DROP TABLE IF EXISTS reviews")
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec("DROP TABLE IF EXISTS users")
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec("DROP TABLE IF EXISTS books")

	if err != nil {
		log.Println(err)
	}
}

func (s *service) createTables() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// create users table
	_, err := s.db.ExecContext(ctx, createTableUsers)
	if err != nil {
		log.Println("User table already exists")
	}
	// create books table
	_, err = s.db.ExecContext(ctx, createTableBooks)
	if err != nil {
		log.Println("Books table already exists")
	}
	_, err = s.db.ExecContext(ctx, createTableReviews)
	if err != nil {
		log.Println("Reviews table already exists")
	}
}

func (s *service) insertUser(email, name, imageURL string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var userID int
	err := s.db.QueryRowContext(ctx, insertUser, email, name, imageURL).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *service) getUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var user models.User
	err := s.db.QueryRowContext(ctx, getUser, email).Scan(&user.ID, &user.Name, &user.Email, &user.ImageURL)

	return user, err
}

func (s *service) insertBook(book models.Book, userID int, review string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var bookID int
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	err = tx.QueryRowContext(ctx, insertBook, book.Title, book.Author, book.Editorial).Scan(&bookID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}
	_, err = tx.ExecContext(ctx, insertReview, bookID, userID, review)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return bookID, nil
}
func (s *service) insertReview(review string, userID int, bookID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var reviewID int
	err := s.db.QueryRowContext(ctx, insertReview, bookID, userID, review).Scan(&reviewID)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}
