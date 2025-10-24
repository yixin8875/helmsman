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
	snapshotsCachePrefixKey = "snapshots:"
	// SnapshotsExpireTime expire time
	SnapshotsExpireTime = 5 * time.Minute
)

var _ SnapshotsCache = (*snapshotsCache)(nil)

// SnapshotsCache cache interface
type SnapshotsCache interface {
	Set(ctx context.Context, id uint64, data *model.Snapshots, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Snapshots, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Snapshots, error)
	MultiSet(ctx context.Context, data []*model.Snapshots, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// snapshotsCache define a cache struct
type snapshotsCache struct {
	cache cache.Cache
}

// NewSnapshotsCache new a cache
func NewSnapshotsCache(cacheType *database.CacheType) SnapshotsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Snapshots{}
		})
		return &snapshotsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Snapshots{}
		})
		return &snapshotsCache{cache: c}
	}

	return nil // no cache
}

// GetSnapshotsCacheKey cache key
func (c *snapshotsCache) GetSnapshotsCacheKey(id uint64) string {
	return snapshotsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *snapshotsCache) Set(ctx context.Context, id uint64, data *model.Snapshots, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetSnapshotsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *snapshotsCache) Get(ctx context.Context, id uint64) (*model.Snapshots, error) {
	var data *model.Snapshots
	cacheKey := c.GetSnapshotsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *snapshotsCache) MultiSet(ctx context.Context, data []*model.Snapshots, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetSnapshotsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *snapshotsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Snapshots, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetSnapshotsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Snapshots)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Snapshots)
	for _, id := range ids {
		val, ok := itemMap[c.GetSnapshotsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *snapshotsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetSnapshotsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *snapshotsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetSnapshotsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *snapshotsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
