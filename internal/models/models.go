package models

import "time"

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

type Review struct {
	ID      int
	Comment string
	User    User
	Date    time.Time
}
