package postgres

import (
	"context"
	"github.com/crewblade/banner-management-service/internal/domain/models"
)

func (s *Storage) GetBanners(ctx context.Context) ([]models.Banner, error) {
	return nil, nil
}
func (s *Storage) SaveBanner(
	ctx context.Context,
	tagIDs []int,
	featureID int,
	content map[string]string,
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
	content map[string]string,
	isActive bool,
) error {
	return nil
}
func (s *Storage) GetUserBanner(
	ctx context.Context,
	tagID int,
	featureID int,
	isAdmin bool,
) (map[string]string, error) {
	return nil, nil
}
