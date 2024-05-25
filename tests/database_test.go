package tests

import (
	"books/internal/database"
	"testing"
)

func TestDatabaseUsers(t *testing.T) {
	// connection for checknig actions, it will set up some initial data
	conn, err := connectDB()
	if err != nil {
		t.Errorf("Error setting up database: %v", err)
		return
	}
	db := database.New(database.TestConfig())
	err = db.Connect()
	if err != nil {
		t.Errorf("Error connecting to the database: %v", err)
		return
	}

	// AddUserIfNotExists
	t.Run("Add new user", func(t *testing.T) {
		setupDB(conn)
		newUserID, err := db.AddUserIfNotExists("new_user", "new_user@mail.com", "new_user_image")
		if err != nil {
			t.Errorf("Error adding new user: %v", err)
			return
		}
		correctID := len(users) + 1
		if newUserID != correctID {
			t.Errorf("Expected user ID to be %d, got %d", correctID, newUserID)
		}
		//check if the new user has been added to the database
		var count int
		err = conn.QueryRow("SELECT COUNT(*) FROM users").Scan(&count)
		if err != nil {
			t.Errorf("Error counting users: %v", err)
			return
		}
		if count != len(users)+1 {
			t.Errorf("Expected %d users, got %d", len(users)+1, count)
		}

	})

}
