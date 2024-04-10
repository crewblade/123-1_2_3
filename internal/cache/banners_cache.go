package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/utils"
	"github.com/patrickmn/go-cache"
	"strings"
)

func (bc *BannerCacheImpl) GetBanners(ctx context.Context, featureID, tagID, limit, offset *int) ([]models.Banner, error) {
	key := buildBannersCacheKey(featureID, tagID, limit, offset)
	data, found := bc.cache.Get(key)
	if !found {
		return nil, fmt.Errorf("banners not found in cache")
	}

	banners := make([]models.Banner, 0)
	err := json.Unmarshal(data.([]byte), &banners)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling banners from cache: %w", err)
	}

	return banners, nil
}

func (bc *BannerCacheImpl) SetBanners(ctx context.Context, featureID, tagID, limit, offset *int, banners []models.Banner) error {
	key := buildBannersCacheKey(featureID, tagID, limit, offset)
	data, err := json.Marshal(banners)
	if err != nil {
		return fmt.Errorf("error marshalling banners for cache: %w", err)
	}

	bc.cache.Set(key, data, cache.DefaultExpiration)
	return nil
}

func buildBannersCacheKey(featureID, tagID, limit, offset *int) string {

	keyParts := []string{
		utils.IntPointertoaOrDefault(featureID, "nil_feature"),
		utils.IntPointertoaOrDefault(tagID, "nil_tag"),
		utils.IntPointertoaOrDefault(limit, "nil_limit"),
		utils.IntPointertoaOrDefault(offset, "nil_offset"),
	}

	return strings.Join(keyParts, "_")
}
