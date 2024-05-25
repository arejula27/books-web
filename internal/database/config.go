package database

import (
	"context"
	"database/sql"
	"log"
	"time"
)

const (
	CreateTableUsers = `
		CREATE TABLE USERS (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		name VARCHAR(100) NOT NULL,
		image_url VARCHAR(255),
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	CreateTableBooks = `
		CREATE TABLE BOOKS (
		id SERIAL PRIMARY KEY,
		volume_id INT NOT NULL,
		registration_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		CONSTRAINT fk_volume
			FOREIGN KEY(volume_id)
				REFERENCES VOLUMES(id)
				ON DELETE CASCADE
				ON UPDATE CASCADE
	);`

	CreateTableVolumes = `
		CREATE TABLE VOLUMES (
		id SERIAL PRIMARY KEY,
		isbn VARCHAR(13) NOT NULL UNIQUE,
		title VARCHAR(255) NOT NULL,
		author VARCHAR(100) NOT NULL,
		editorial VARCHAR(100) NOT NULL
);	`

	CreateTableBooksTags = `
		CREATE TABLE BOOKS_TAGS (
		book_id INT NOT NULL,
		tag_id INT NOT NULL,
		CONSTRAINT pk_books_tags
			PRIMARY KEY (book_id, tag_id),
		CONSTRAINT fk_book
			FOREIGN KEY(book_id)
				REFERENCES BOOKS(id)
				ON DELETE CASCADE
				ON UPDATE CASCADE,
		CONSTRAINT fk_tag
			FOREIGN KEY(tag_id)
				REFERENCES TAGS(id)
				ON DELETE CASCADE	
				ON UPDATE CASCADE
	);`

	CreateTableTimelineRecords = `
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
	CreateTableReviews = `
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
	CreateTableTags = `
		CREATE TABLE TAGS (
		id SERIAL PRIMARY KEY,
		name VARCHAR(50) NOT NULL UNIQUE
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
	_, err := s.db.ExecContext(ctx, CreateTableUsers)
	if err != nil {
		s.logMsg("User table already exists")
	}
	// create tags create
	_, err = s.db.ExecContext(ctx, CreateTableTags)
	if err != nil {
		s.logMsg("Tags table already exists")
	}
	// create volumes table
	_, err = s.db.ExecContext(ctx, CreateTableVolumes)
	if err != nil {
		s.logMsg("Volumes table already exists")
	}
	// create books table
	_, err = s.db.ExecContext(ctx, CreateTableBooks)
	if err != nil {
		s.logMsg("Books table already exists")
	}
	_, err = s.db.ExecContext(ctx, CreateTableTimelineRecords)
	if err != nil {
		s.logMsg("Timeline records table already exists")
	}
	_, err = s.db.ExecContext(ctx, CreateTableReviews)
	if err != nil {
		s.logMsg("Reviews table already exists")
	}
	_, err = s.db.ExecContext(ctx, CreateTableBooksTags)
	if err != nil {
		s.logMsg("Books tags table already exists")
	}

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
		s.createTables()

		//populate default tags
		defaultTags := []string{"fiction", "terror", "classic", "romance", "sci-fiction", "mystery", "thriller", "biography", "fantasy", "non-fiction"}
		err := s.insertTags(defaultTags)
		if nil != err {
			return err
		}
	}

	return nil

}

// Health checks the health of the database, returns a map with a message or stops the application if the database is down
func (s *service) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.PingContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
