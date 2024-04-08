package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/errs"
	"github.com/lib/pq"
)

func (s *Storage) GetBanners(ctx context.Context) ([]models.Banner, error) {
	return nil, nil
}
func (s *Storage) SaveBanner(
	ctx context.Context,
	tagIDs []int,
	featureID int,
	content string,
	isActive bool,
) (int, error) {
	const op = "repo.postgres.SaveBanner"
	const errValue = -1
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return errValue, fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()
	stmt, err := s.db.PrepareContext(ctx, "INSERT INTO banners(tag_ids, feature_id, content, is_active) VALUES ($1, $2, $3, $4) RETURNING id")
	if err != nil {
		return errValue, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	var bannerID int
	err = stmt.QueryRowContext(ctx, pq.Array(tagIDs), featureID, content, isActive).Scan(&bannerID)
	if err != nil {
		return errValue, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return errValue, fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return bannerID, nil
}

func (s *Storage) DeleteBanner(ctx context.Context, bannerID int) error {
	const op = "repo.postgres.DeleteBanner"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: begin transaction: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "DELETE FROM banners WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: prepare context: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, bannerID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) UpdateBanner(
	ctx context.Context,
	bannerID int,
	tagIDs []int,
	featureID int,
	content string,
	isActive bool,
) error {
	const op = "repo.postgres.UpdateBanner"

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("%s: start transaction: %w", op, err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE banners SET tag_ids = $1, feature_id = $2, content = $3, is_active = $4, updated_at = NOW() WHERE id = $5")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, pq.Array(tagIDs), featureID, content, isActive, bannerID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return nil
}
func (s *Storage) GetUserBanner(
	ctx context.Context,
	tagID int,
	featureID int,
) (string, bool, error) {

	const op = "repo.postgres.GetUserBanner"

	stmt, err := s.db.PrepareContext(ctx, "SELECT content, is_active FROM banners WHERE feature_id = $1 AND $2 = ANY(tag_ids)")
	if err != nil {
		return "", false, fmt.Errorf("%s: prepare context %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, featureID, tagID)
	var content string
	var isActive bool
	err = row.Scan(&content, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, fmt.Errorf("%s: row scan %f", op, errs.ErrBannerNotFound)
		}

		return "", false, fmt.Errorf("%s: row scan %w", op, err)
	}

	return content, isActive, nil

}
