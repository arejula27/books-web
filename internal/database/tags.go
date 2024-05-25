package database

import (
	"context"
	"time"
)

func (s *service) insertTags(tags []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `INSERT INTO tags (name) VALUES ($1)`
	for _, tag := range tags {
		_, err := s.db.ExecContext(ctx, query, tag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *service) GetAllTags() ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	query := `SELECT name FROM tags`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tags []string
	for rows.Next() {
		var tag string
		err = rows.Scan(&tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
