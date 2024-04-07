package postgres

import "context"

func (s *Storage) IsAdmin(ctx context.Context, token string) (bool, error) {
	return true, nil
}
