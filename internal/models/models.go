package models

// User is a struct that represents a user
type User struct {
	ID       int
	Name     string
	Email    string
	ImageURL string
}

type Book struct {
	ID        int
	Title     string
	Author    string
	Editorial string
}
