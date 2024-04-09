package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

type BannerCacheImpl struct {
	cache *cache.Cache
}

func NewBannerCacheImpl(defaultExpiration, cleanupInterval time.Duration) *BannerCacheImpl {
	c := cache.New(defaultExpiration, cleanupInterval)
	return &BannerCacheImpl{cache: c}
}
