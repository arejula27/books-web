package tests

import (
	"books/internal/database"
	"books/internal/models"
	"log"
)

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
		},
		{
			Title:     "test_title2",
			Author:    "test_autor2",
			Editorial: "test_editorial2",
		},
		{
			Title:     "test_title3",
			Author:    "test_autor3",
			Editorial: "test_editorial3",
		},
	}
)

func setupDB() (*database.Service, error) {

	db := database.New(database.TestConfig(), database.ResetDatabase())
	err := db.Connect()
	if err != nil {
		return nil, err
	}
	//add two test users
	for i, user := range users {
		users[i].ID, err = db.AddUserIfNotExists(user.Name, user.Email, "test_image")
		if err != nil {
			return nil, err
		}
	}
	log.Println(users[0].ID, users[1].ID)

	//add books

	for i, book := range books {

		book.ID, err = db.AddBook(book, users[i%2].ID, "")
		if err != nil {
			return nil, err
		}
	}
	return &db, nil
}
