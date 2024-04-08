package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/errs"
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
	return 1, nil
}
func (s *Storage) DeleteBanner(ctx context.Context, bannerID int) error {
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
	return nil
}
func (s *Storage) GetUserBanner(
	ctx context.Context,
	tagID int,
	featureID int,
	isAdmin bool,
) (string, bool, error) {

	const op = "repo.postgres.GetUserBanner"

	stmt, err := s.db.Prepare("SELECT content, is_active FROM banners WHERE feature_id = $1 AND $2 = ANY(tag_ids)")
	if err != nil {
		return "", false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, featureID, tagID)
	var content string
	var isActive bool
	err = row.Scan(&content, &isActive)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", false, fmt.Errorf("%s: %f", op, errs.ErrBannerNotFound)
		}

		return "", false, fmt.Errorf("%s: %w", op, err)
	}

	return content, isAdmin, nil

}
