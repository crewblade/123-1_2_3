package postgres

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(connectionString string) (*Storage, error) {
	const op = "repo.postgres.New"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	return &Storage{db: db}, nil
}

func (s *Storage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

func (s *Storage) ClearData(ctx context.Context) error {
	const op = "repo.postgres.ClearData"

	_, err := s.db.ExecContext(ctx, "DELETE FROM banners")
	if err != nil {
		return fmt.Errorf("%s: executing delete banners query '%s': %w", op, err)
	}

	_, err = s.db.ExecContext(ctx, "DELETE FROM users")
	if err != nil {
		return fmt.Errorf("%s: executing delete users query '%s': %w", op, err)
	}
	return nil

}
