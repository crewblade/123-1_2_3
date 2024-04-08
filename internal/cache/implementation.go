package cache

import (
	"context"
	"errors"
	"fmt"
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

func (bc *BannerCacheImpl) GetBannerContent(ctx context.Context, tagID, featureID int) (string, bool, error) {
	content, ok := bc.cache.Get(fmt.Sprintf("%d_%d", tagID, featureID))
	if !ok {
		return "", false, errors.New("content is not found in cache")
	}
	return content.(string), true, nil
}

func (bc *BannerCacheImpl) SetBannerContent(ctx context.Context, tagID, featureID int, content string, isActive bool) error {
	bc.cache.Set(fmt.Sprintf("%d_%d", tagID, featureID), content, cache.DefaultExpiration)
	return nil
}
