package postgres

import (
	"context"
	"fmt"
)

func (s *Storage) IsAdmin(ctx context.Context, token string) (bool, error) {
	const op = "repo.postgres.IsAdmin"
	stmt, err := s.db.Prepare("SELECT is_admin FROM users WHERE token = $1")
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, token)
	var isAdmin bool

	err = row.Scan(&isAdmin)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}

func (s *Storage) AddUser(ctx context.Context, token string, isAdmin bool) error {
	const op = "repo.postgres.AddUser"

	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO users (token, is_admin) VALUES ($1, $2)")
	if err != nil {
		return fmt.Errorf("%s: preparing statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, token, isAdmin)
	if err != nil {
		return fmt.Errorf("%s: executing statement: %w", op, err)
	}

	return nil
}
