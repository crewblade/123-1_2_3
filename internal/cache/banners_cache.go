package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"time"
)

func (bc *BannerCacheImpl) GetBanners(ctx context.Context, featureID, tagID, limit, offset int) ([]models.Banner, bool, error) {
	key := fmt.Sprintf("%d_%d_%d_%d", featureID, tagID, limit, offset)
	data, found := bc.cache.Get(key)
	if !found {
		return nil, false, errors.New("banners not found in cache")
	}

	banners := make([]models.Banner, 0)
	err := json.Unmarshal(data.([]byte), &banners)
	if err != nil {
		return nil, false, fmt.Errorf("error unmarshalling banners from cache: %w", err)
	}

	return banners, true, nil
}

func (bc *BannerCacheImpl) SetBanners(ctx context.Context, featureID, tagID, limit, offset int, banners []models.Banner, expiration time.Duration) error {
	key := fmt.Sprintf("%d_%d_%d_%d", featureID, tagID, limit, offset)
	data, err := json.Marshal(banners)
	if err != nil {
		return fmt.Errorf("error marshalling banners for cache: %w", err)
	}

	bc.cache.Set(key, data, expiration)
	return nil
}
