package database_test

import (
	"books/internal/database"
	"books/utils"
	"testing"
)

func init() {
	utils.LoadEnv()
}

func TestDatabaseUsers(t *testing.T) {
	// connection for checknig actions, it will set up some initial data
	conn, err := utils.ConnectDB()
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
		testData := utils.TestScenario1Data()
		utils.SetupDB(conn, &testData)
		newUserID, err := db.AddUserIfNotExists("new_user", "new_user@mail.com", "new_user_image")
		if err != nil {
			t.Errorf("Error adding new user: %v", err)
			return
		}
		correctID := len(testData.Users) + 1
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
		if count != len(testData.Users)+1 {
			t.Errorf("Expected %d users, got %d", len(testData.Users)+1, count)
		}

	})

}
