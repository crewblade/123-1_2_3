package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

type Storage struct {
	db *pgxpool.Pool
}

func New(connectionString string, minConnections int, maxConnections int, maxConnIdleTime uint64) (*Storage, error) {
	const op = "repo.postgres.New"
	cfx, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	cfx.MaxConns = int32(maxConnections)
	cfx.MinConns = int32(minConnections)
	cfx.MaxConnIdleTime = time.Duration(maxConnIdleTime) * time.Millisecond

	db, err := pgxpool.NewWithConfig(context.Background(), cfx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if err = db.Ping(context.Background()); err != nil {
		db.Close()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

//func (s *Storage) Close() error {
//	if s.db != nil {
//		return s.db.Close()
//	}
//	return nil
//}

func (s *Storage) ClearData(ctx context.Context) error {
	const op = "repo.postgres.ClearData"

	_, err := s.db.Exec(ctx, "DELETE FROM banners")
	if err != nil {
		return fmt.Errorf("%s: executing delete banners query : %w", op, err)
	}

	_, err = s.db.Exec(ctx, "DELETE FROM users")
	if err != nil {
		return fmt.Errorf("%s: executing delete users query: %w", op, err)
	}
	return nil

}

func (s *Storage) CountRows(ctx context.Context) (int, error) {
	const op = "repo.postgres.CountRows"

	var count int
	err := s.db.QueryRow(ctx, "SELECT COUNT(*) FROM banners").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("%s: execute query: %w", op, err)
	}

	return count, nil
}

func (s *Storage) PrepareForTest(ctx context.Context) error {
	const op = "repo.postgres.PrepareForTest"

	queryBanners := `CREATE TABLE IF NOT EXISTS banners (
                                       id SERIAL PRIMARY KEY,
                                       content JSONB NOT NULL,
                                       feature_id INT NOT NULL,
                                       tag_ids INT[] NOT NULL,
                                       is_active BOOLEAN NOT NULL DEFAULT true,
                                       created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                       deleted BOOLEAN NOT NULL DEFAULT false
);`
	_, err := s.db.Exec(ctx, queryBanners)
	if err != nil {
		return fmt.Errorf("%s: executing delete banners query : %w", op, err)
	}

	queryUsers := `CREATE TABLE IF NOT EXISTS users (
                                     token TEXT NOT NULL,
                                     is_admin BOOLEAN NOT NULL DEFAULT false
);`

	_, err = s.db.Exec(ctx, queryUsers)
	if err != nil {
		return fmt.Errorf("%s: executing delete users query: %w", op, err)
	}
	return nil
}
