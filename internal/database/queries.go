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
	createTableTimelineRecords = `
	CREATE TABLE TIMELINE_RECORDS (
		id SERIAL PRIMARY KEY,
		book_id INT NOT NULL, 
		user_id INT NOT NULL,
		received BOOLEAN DEFAULT FALSE,
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
	createTableReviews = `
		CREATE TABLE REVIEWS (
		timeline_record_id INT PRIMARY KEY,
		comment TEXT NOT NULL,
		CONSTRAINT fk_timeline_record
			FOREIGN KEY(timeline_record_id)
			REFERENCES TIMELINE_RECORDS(id)
			ON DELETE CASCADE
			ON UPDATE CASCADE
	);
	`
)

func (s *service) clearDatabase() {

	_, err := s.db.Exec("DROP TABLE IF EXISTS reviews")
	if err != nil {
		log.Println(err)
	}
	_, err = s.db.Exec("DROP TABLE IF EXISTS timeline_records")
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
	_, err = s.db.ExecContext(ctx, createTableTimelineRecords)
	if err != nil {
		log.Println("Timeline_records table already exists")
	}
	_, err = s.db.ExecContext(ctx, createTableReviews)
	if err != nil {
		log.Println("Reviews table already exists")
	}

}

func (s *service) insertUser(email, name, imageURL string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `INSERT INTO users (email, name, image_url) VALUES ($1, $2,$3) RETURNING id`
	var userID int
	err := s.db.QueryRowContext(ctx, query, email, name, imageURL).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (s *service) getUserByEmail(email string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := "SELECT id, name, email, image_url  FROM users WHERE email = $1"
	var user models.User
	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Name, &user.Email, &user.ImageURL)

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
	insertBook := `INSERT INTO books (title, author, editorial) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRowContext(ctx, insertBook, book.Title, book.Author, book.Editorial).Scan(&bookID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}
	var recordID int
	insertTimelineRecord := `INSERT INTO timeline_records (book_id, user_id,received) VALUES ($1, $2, $3 ) RETURNING id`
	err = tx.QueryRowContext(ctx, insertTimelineRecord, bookID, userID, true).Scan(&recordID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	insertReview := `INSERT INTO reviews (timeline_record_id, comment) VALUES ($1, $2) RETURNING timeline_record_id`
	_, err = tx.ExecContext(ctx, insertReview, recordID, review)
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
func (s *service) insertTimelineRecord(userID int, bookID int) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	var reviewID int
	query := `INSERT INTO timeline_records (book_id, user_id) VALUES ($1, $2 ) RETURNING id`
	err := s.db.QueryRowContext(ctx, query, bookID, userID).Scan(&reviewID)
	if err != nil {
		return 0, err
	}
	return reviewID, nil
}
func (s *service) selectBooksFromUser(userID int) ([]models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `SELECT b.id, b.title, b.author, b.editorial FROM books 
		b INNER JOIN timeline_records tr ON b.id = tr.book_id WHERE tr.user_id = $1
		ORDER BY tr.creation_date DESC`
	rows, err := s.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []models.Book
	for rows.Next() {
		var book models.Book
		err = rows.Scan(&book.ID, &book.Title, &book.Author, &book.Editorial)
		if err != nil {
			return nil, err
		}
		books = append(books, book)
	}
	return books, nil
}
