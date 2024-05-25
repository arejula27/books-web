package database

import (
	"books/internal/models"
	"context"
	"time"
)

func (s *service) AddUserIfNotExists(email, name, imageURL string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `INSERT INTO users (email, name, image_url) VALUES ($1, $2,$3) ON CONFLICT (email) DO UPDATE SET id = users.id RETURNING id`
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
