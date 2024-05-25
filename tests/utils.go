package tests

import (
	"books/internal/database"
	"books/internal/models"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

////////////////////////////////////////////////////////////
// The user test_name has 3 books, each with a different tag
////////////////////////////////////////////////////////////

var (
	users = []models.User{
		{
			Name:  "test_name",
			Email: "test_mail@mail.test",
		},
		{
			Name:  "test_name2",
			Email: "test_mail2@mail.test",
		},
	}
	books = []models.Book{
		{
			Title:     "test_title",
			Author:    "test_autor",
			Editorial: "test_editorial",
			ISBN:      "test_isbn",
		},
		{
			Title:     "test_title2",
			Author:    "test_autor2",
			Editorial: "test_editorial2",
			ISBN:      "test_isbn2",
		},
		{
			Title:     "test_title3",
			Author:    "test_autor3",
			Editorial: "test_editorial3",
			ISBN:      "test_isbn3",
		},
	}
	tags = []string{
		"tag1",
		"tag2",
		"tag3",
	}
)

func connectDB() (*sql.DB, error) {
	err := godotenv.Load("../.env.test")
	if err != nil {
		log.Println("Error loading .env file")
		return nil, err
	}

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_DATABASE")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, databaseName)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	return db, nil
}

func setupDB(db *sql.DB) error {
	//clear database
	_, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		log.Fatalf("Error dropping database: %v", err)
	}

	//create tables
	tables := []string{
		database.CreateTableUsers,
		database.CreateTableVolumes,
		database.CreateTableTags,
		database.CreateTableBooks,
		database.CreateTableBooksTags,
		database.CreateTableTimelineRecords,
		database.CreateTableReviews,
	}

	for _, table := range tables {
		_, err = db.Exec(table)
		if err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
	}

	//insert some users
	for i, user := range users {
		var id int
		err = db.QueryRow("INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id", user.Name, user.Email).Scan(&id)
		if err != nil {

			log.Fatalf("Error inserting user: %v", err)
		}
		users[i].ID = id
	}

	//insert tags
	for _, tag := range tags {
		var id int
		err = db.QueryRow("INSERT INTO tags (name) VALUES ($1) RETURNING id", tag).Scan(&id)
		if err != nil {
			log.Fatalf("Error inserting tag: %v", err)
		}
	}

	//insert books
	for i, book := range books {
		var volumeID int
		err = db.QueryRow("INSERT INTO volumes (title, author, editorial, isbn) VALUES ($1, $2, $3, $4) RETURNING id", book.Title, book.Author, book.Editorial, book.ISBN).Scan(&volumeID)
		if err != nil {
			log.Fatalf("Error inserting volume: %v", err)
		}
		var bookID int
		err = db.QueryRow("INSERT INTO books (volume_id) VALUES ($1) RETURNING id", volumeID).Scan(&bookID)
		if err != nil {
			log.Fatalf("Error inserting book: %v", err)
		}
		//insert tags
		tagID := (i % len(tags)) + 1
		_, err = db.Exec("INSERT INTO books_tags (book_id, tag_id) VALUES ($1, $2)", bookID, tagID)
		if err != nil {
			log.Fatalf("Error inserting book tag: %v", err)
		}

		//insert timeline records
		var recordID int
		err = db.QueryRow("INSERT INTO timeline_records (user_id, book_id, received) VALUES ($1, $2,true) RETURNING id", 1, bookID).Scan(&recordID)
		if err != nil {
			log.Fatalf("Error inserting timeline record: %v", err)
		}
		//insert reviews
		_, err = db.Exec("INSERT INTO reviews (timeline_record_id, comment) VALUES ($1, $2)", recordID, "test_review")

	}

	return nil
}
