package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"time"

	"github.com/patrickmn/go-cache"
)

type BannerCacheImpl struct {
	cache *cache.Cache
}

func NewBannerCacheImpl(defaultExpiration, cleanupInterval time.Duration) *BannerCacheImpl {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &BannerCacheImpl{cache: c}
}

func (bc *BannerCacheImpl) GetBanner(ctx context.Context, tagID, featureID int) (string, bool, error) {
	key := fmt.Sprintf("%d_%d", tagID, featureID)
	data, found := bc.cache.Get(key)
	if !found {
		return "", false, fmt.Errorf("banner is not found in cache")
	}

	var banner models.BannerForUser
	err := json.Unmarshal(data.([]byte), &banner)
	if err != nil {
		return "", false, fmt.Errorf("error unmarshalling banner from cache: %w", err)
	}

	return banner.Content, banner.IsActive, nil
}

func (bc *BannerCacheImpl) SetBanner(ctx context.Context, tagID, featureID int, banner *models.BannerForUser) error {
	key := fmt.Sprintf("%d_%d", tagID, featureID)
	data, err := json.Marshal(banner)
	if err != nil {
		return fmt.Errorf("error marshalling banner: %w", err)
	}

	bc.cache.Set(key, data, cache.DefaultExpiration)
	return nil
}
