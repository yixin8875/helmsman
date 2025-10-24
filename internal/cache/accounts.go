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
	accountsCachePrefixKey = "accounts:"
	// AccountsExpireTime expire time
	AccountsExpireTime = 5 * time.Minute
)

var _ AccountsCache = (*accountsCache)(nil)

// AccountsCache cache interface
type AccountsCache interface {
	Set(ctx context.Context, id uint64, data *model.Accounts, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Accounts, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Accounts, error)
	MultiSet(ctx context.Context, data []*model.Accounts, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// accountsCache define a cache struct
type accountsCache struct {
	cache cache.Cache
}

// NewAccountsCache new a cache
func NewAccountsCache(cacheType *database.CacheType) AccountsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Accounts{}
		})
		return &accountsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Accounts{}
		})
		return &accountsCache{cache: c}
	}

	return nil // no cache
}

// GetAccountsCacheKey cache key
func (c *accountsCache) GetAccountsCacheKey(id uint64) string {
	return accountsCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *accountsCache) Set(ctx context.Context, id uint64, data *model.Accounts, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetAccountsCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *accountsCache) Get(ctx context.Context, id uint64) (*model.Accounts, error) {
	var data *model.Accounts
	cacheKey := c.GetAccountsCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *accountsCache) MultiSet(ctx context.Context, data []*model.Accounts, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetAccountsCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *accountsCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Accounts, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetAccountsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Accounts)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Accounts)
	for _, id := range ids {
		val, ok := itemMap[c.GetAccountsCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *accountsCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetAccountsCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *accountsCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetAccountsCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *accountsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
