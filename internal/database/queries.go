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
		volume_id INT NOT NULL,
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_volume
			FOREIGN KEY(volume_id)
				REFERENCES VOLUMES(id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	);
	`
	createTableVolumes = `
		CREATE TABLE VOLUMES (
		id SERIAL PRIMARY KEY,
		isbn VARCHAR(13) NOT NULL UNIQUE,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(100) NOT NULL,
		editorial VARCHAR(100) NOT NULL
);	`
	createTableTimelineRecords = `
	CREATE TABLE TIMELINE_RECORDS (
		id SERIAL PRIMARY KEY,
		book_id INT NOT NULL, 
		user_id INT NOT NULL,
		received BOOLEAN DEFAULT FALSE,
		willing_to_send BOOLEAN DEFAULT FALSE,
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

// clearDatabase drops all tables
func (s *service) clearDatabase() {
	_, err := s.db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		log.Fatalf("Error dropping database: %v", err)
	}
}

func (s *service) createTables() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// create users table
	_, err := s.db.ExecContext(ctx, createTableUsers)
	if err != nil {
		s.logMsg("User table already exists")
	}
	// create volumes table
	_, err = s.db.ExecContext(ctx, createTableVolumes)
	if err != nil {
		s.logMsg("Volumes table already exists")
	}
	// create books table
	_, err = s.db.ExecContext(ctx, createTableBooks)
	if err != nil {
		s.logMsg("Books table already exists")
	}
	_, err = s.db.ExecContext(ctx, createTableTimelineRecords)
	if err != nil {
		s.logMsg("Timeline records table already exists")
	}
	_, err = s.db.ExecContext(ctx, createTableReviews)
	if err != nil {
		s.logMsg("Reviews table already exists")
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
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	var VolumeID int
	// insert volume
	// if the volume already exists, update the id for returning it
	//When updating the query will detect a change on a row and will return it, if not
	//the query will not trigger any change and will return nothing
	insertVolume := `INSERT INTO volumes (isbn, title, author, editorial) VALUES ($1, $2, $3, $4)
	ON CONFLICT (isbn) DO UPDATE SET id = volumes.id RETURNING id;`
	err = tx.QueryRowContext(ctx, insertVolume, book.ISBN, book.Title, book.Author, book.Editorial).Scan(&VolumeID)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	var bookID int
	insertBook := `INSERT INTO books (volume_id) VALUES ($1) RETURNING id`
	err = tx.QueryRowContext(ctx, insertBook, VolumeID).Scan(&bookID)
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
	query := `SELECT b.id, v.title, v.author, v.editorial FROM volumes v INNER JOIN books 
		b ON v.id = b.volume_id INNER JOIN timeline_records tr ON b.id = tr.book_id WHERE tr.user_id = $1
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
func (s *service) selectBookByID(bookID int) (models.Book, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `SELECT b.id, v.title, v.author, v.editorial v,isbn FROM volumes v INNER JOIN books b ON v.id = b.volume_id WHERE b.id = $1`
	var book models.Book
	err := s.db.QueryRowContext(ctx, query, bookID).Scan(&book.ID, &book.Title, &book.Author, &book.Editorial, &book.ISBN)
	if err != nil {
		return models.Book{}, err
	}
	return book, nil
}

func (s *service) selectReviewsByBookID(bookID int) ([]models.Review, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `SELECT r.comment,tr.creation_date, u.name FROM reviews r LEFT JOIN timeline_records tr ON r.timeline_record_id = tr.id
		INNER JOIN users u ON tr.user_id = u.id WHERE tr.book_id = $1 
		ORDER BY tr.creation_date DESC`
	rows, err := s.db.QueryContext(ctx, query, bookID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		review.User = models.User{}
		err = rows.Scan(&review.Comment, &review.Date, &review.User.Name)
		if err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	return reviews, nil
}
