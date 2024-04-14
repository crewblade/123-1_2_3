package postgres

import (
	"context"
	"fmt"
)

func (s *Storage) IsAdmin(ctx context.Context, token string) (bool, error) {
	const op = "repo.postgres.IsAdmin"

	var isAdmin bool
	err := s.db.QueryRow(ctx, "SELECT is_admin FROM users WHERE token = $1", token).Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) AddUser(ctx context.Context, token string, isAdmin bool) error {
	const op = "repo.postgres.AddUser"

	_, err := s.db.Exec(ctx, "INSERT INTO users (token, is_admin) VALUES ($1, $2)", token, isAdmin)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
