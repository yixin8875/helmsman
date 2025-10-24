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
	tradesCachePrefixKey = "trades:"
	// TradesExpireTime expire time
	TradesExpireTime = 5 * time.Minute
)

var _ TradesCache = (*tradesCache)(nil)

// TradesCache cache interface
type TradesCache interface {
	Set(ctx context.Context, id uint64, data *model.Trades, duration time.Duration) error
	Get(ctx context.Context, id uint64) (*model.Trades, error)
	MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Trades, error)
	MultiSet(ctx context.Context, data []*model.Trades, duration time.Duration) error
	Del(ctx context.Context, id uint64) error
	SetPlaceholder(ctx context.Context, id uint64) error
	IsPlaceholderErr(err error) bool
}

// tradesCache define a cache struct
type tradesCache struct {
	cache cache.Cache
}

// NewTradesCache new a cache
func NewTradesCache(cacheType *database.CacheType) TradesCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.Trades{}
		})
		return &tradesCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.Trades{}
		})
		return &tradesCache{cache: c}
	}

	return nil // no cache
}

// GetTradesCacheKey cache key
func (c *tradesCache) GetTradesCacheKey(id uint64) string {
	return tradesCachePrefixKey + utils.Uint64ToStr(id)
}

// Set write to cache
func (c *tradesCache) Set(ctx context.Context, id uint64, data *model.Trades, duration time.Duration) error {
	if data == nil || id == 0 {
		return nil
	}
	cacheKey := c.GetTradesCacheKey(id)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *tradesCache) Get(ctx context.Context, id uint64) (*model.Trades, error) {
	var data *model.Trades
	cacheKey := c.GetTradesCacheKey(id)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *tradesCache) MultiSet(ctx context.Context, data []*model.Trades, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTradesCacheKey(v.ID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is id value
func (c *tradesCache) MultiGet(ctx context.Context, ids []uint64) (map[uint64]*model.Trades, error) {
	var keys []string
	for _, v := range ids {
		cacheKey := c.GetTradesCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.Trades)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[uint64]*model.Trades)
	for _, id := range ids {
		val, ok := itemMap[c.GetTradesCacheKey(id)]
		if ok {
			retMap[id] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *tradesCache) Del(ctx context.Context, id uint64) error {
	cacheKey := c.GetTradesCacheKey(id)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *tradesCache) SetPlaceholder(ctx context.Context, id uint64) error {
	cacheKey := c.GetTradesCacheKey(id)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *tradesCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
