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
