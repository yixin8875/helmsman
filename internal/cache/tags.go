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
	tagsCachePrefixKey = "tags:"
	// TagsExpireTime expire time
	TagsExpireTime = 5 * time.Minute
)

var _ TagsCache = (*tagsCache)(nil)

// TagsCache cache interface
type TagsCache interface {
	Set(ctx context.Context, id uint64, data *model.Tags, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Tags, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Tags, error)
	MultiSet(ctx context.Context, data []*model.Tags, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// tagsCache define a cache struct
type tagsCache struct {
	cache cache.Cache
}

// NewTagsCache new a cache
func NewTagsCache(cacheType *database.CacheType) TagsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Tags{}
		})
		return &tagsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Tags{}
		})
		return &tagsCache{cache: c}
	}

	return nil // no cache
}

// GetTagsCacheKey cache key
func (c *tagsCache) GetTagsCacheKey(id uint64) string {
	return tagsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *tagsCache) Set(ctx context.Context, id uint64, data *model.Tags, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTagsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *tagsCache) Get(ctx context.Context, id uint64) (*model.Tags, error) {
	var data *model.Tags
	cacheKey := c.GetTagsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *tagsCache) MultiSet(ctx context.Context, data []*model.Tags, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTagsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *tagsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Tags, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTagsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Tags)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Tags)
	for _, id := range ids {
		val, ok := itemMap[c.GetTagsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *tagsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTagsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *tagsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetTagsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *tagsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
