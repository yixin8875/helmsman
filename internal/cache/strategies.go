package cache

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/go-dev-frame/sponge/pkg/cache"
	"github.com/go-dev-frame/sponge/pkg/encoding"
	"github.com/go-dev-frame/sponge/pkg/utils"

	"helmsman/internal/database"
	"helmsman/internal/model"
)

const (
	// cache prefix key, must end with a colon
	strategiesCachePrefixKey = "strategies:"
	// StrategiesExpireTime expire time
	StrategiesExpireTime = 5 * time.Minute
)

var _ StrategiesCache = (*strategiesCache)(nil)

// StrategiesCache cache interface
type StrategiesCache interface {
	Set(ctx context.Context, id uint64, data *model.Strategies, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Strategies, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Strategies, error)
	MultiSet(ctx context.Context, data []*model.Strategies, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// strategiesCache define a cache struct
type strategiesCache struct {
	cache cache.Cache
}

// NewStrategiesCache new a cache
func NewStrategiesCache(cacheType *database.CacheType) StrategiesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Strategies{}
		})
		return &strategiesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Strategies{}
		})
		return &strategiesCache{cache: c}
	}

	return nil // no cache
}

// GetStrategiesCacheKey cache key
func (c *strategiesCache) GetStrategiesCacheKey(id uint64) string {
	return strategiesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *strategiesCache) Set(ctx context.Context, id uint64, data *model.Strategies, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetStrategiesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *strategiesCache) Get(ctx context.Context, id uint64) (*model.Strategies, error) {
	var data *model.Strategies
	cacheKey := c.GetStrategiesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *strategiesCache) MultiSet(ctx context.Context, data []*model.Strategies, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetStrategiesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *strategiesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Strategies, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetStrategiesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Strategies)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Strategies)
	for _, id := range ids {
		val, ok := itemMap[c.GetStrategiesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *strategiesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetStrategiesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *strategiesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetStrategiesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *strategiesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
