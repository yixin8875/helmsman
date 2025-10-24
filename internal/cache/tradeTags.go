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
	tradeTagsCachePrefixKey = "tradeTags:"
	// TradeTagsExpireTime expire time
	TradeTagsExpireTime = 5 * time.Minute
)

var _ TradeTagsCache = (*tradeTagsCache)(nil)

// TradeTagsCache cache interface
type TradeTagsCache interface {
	Set(ctx context.Context, tradeID int, data *model.TradeTags, duration time.Duration) error
	Get(ctx context.Context, tradeID int) (*model.TradeTags, error)
	MultiGet(ctx context.Context, tradeIDs []int) (map[int]*model.TradeTags, error)
	MultiSet(ctx context.Context, data []*model.TradeTags, duration time.Duration) error
	Del(ctx context.Context, tradeID int) error
	SetPlaceholder(ctx context.Context, tradeID int) error
	IsPlaceholderErr(err error) bool
}

// tradeTagsCache define a cache struct
type tradeTagsCache struct {
	cache cache.Cache
}

// NewTradeTagsCache new a cache
func NewTradeTagsCache(cacheType *database.CacheType) TradeTagsCache {
	jsonEncoding := encoding.JSONEncoding{}
	cachePrefix := ""

	cType := strings.ToLower(cacheType.CType)
	switch cType {
	case "redis":
		c := cache.NewRedisCache(cacheType.Rdb, cachePrefix, jsonEncoding, func() interface{} {
			return &model.TradeTags{}
		})
		return &tradeTagsCache{cache: c}
	case "memory":
		c := cache.NewMemoryCache(cachePrefix, jsonEncoding, func() interface{} {
			return &model.TradeTags{}
		})
		return &tradeTagsCache{cache: c}
	}

	return nil // no cache
}

// GetTradeTagsCacheKey cache key
func (c *tradeTagsCache) GetTradeTagsCacheKey(tradeID int) string {
	return tradeTagsCachePrefixKey + utils.IntToStr(tradeID)
}

// Set write to cache
func (c *tradeTagsCache) Set(ctx context.Context, tradeID int, data *model.TradeTags, duration time.Duration) error {
	if data == nil {
		return nil
	}
	cacheKey := c.GetTradeTagsCacheKey(tradeID)
	err := c.cache.Set(ctx, cacheKey, data, duration)
	if err != nil {
		return err
	}
	return nil
}

// Get cache value
func (c *tradeTagsCache) Get(ctx context.Context, tradeID int) (*model.TradeTags, error) {
	var data *model.TradeTags
	cacheKey := c.GetTradeTagsCacheKey(tradeID)
	err := c.cache.Get(ctx, cacheKey, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// MultiSet multiple set cache
func (c *tradeTagsCache) MultiSet(ctx context.Context, data []*model.TradeTags, duration time.Duration) error {
	valMap := make(map[string]interface{})
	for _, v := range data {
		cacheKey := c.GetTradeTagsCacheKey(v.TradeID)
		valMap[cacheKey] = v
	}

	err := c.cache.MultiSet(ctx, valMap, duration)
	if err != nil {
		return err
	}

	return nil
}

// MultiGet multiple get cache, return key in map is tradeID value
func (c *tradeTagsCache) MultiGet(ctx context.Context, tradeIDs []int) (map[int]*model.TradeTags, error) {
	var keys []string
	for _, v := range tradeIDs {
		cacheKey := c.GetTradeTagsCacheKey(v)
		keys = append(keys, cacheKey)
	}

	itemMap := make(map[string]*model.TradeTags)
	err := c.cache.MultiGet(ctx, keys, itemMap)
	if err != nil {
		return nil, err
	}

	retMap := make(map[int]*model.TradeTags)
	for _, tradeID := range tradeIDs {
		val, ok := itemMap[c.GetTradeTagsCacheKey(tradeID)]
		if ok {
			retMap[tradeID] = val
		}
	}

	return retMap, nil
}

// Del delete cache
func (c *tradeTagsCache) Del(ctx context.Context, tradeID int) error {
	cacheKey := c.GetTradeTagsCacheKey(tradeID)
	err := c.cache.Del(ctx, cacheKey)
	if err != nil {
		return err
	}
	return nil
}

// SetPlaceholder set placeholder value to cache
func (c *tradeTagsCache) SetPlaceholder(ctx context.Context, tradeID int) error {
	cacheKey := c.GetTradeTagsCacheKey(tradeID)
	return c.cache.SetCacheWithNotFound(ctx, cacheKey)
}

// IsPlaceholderErr check if cache is placeholder error
func (c *tradeTagsCache) IsPlaceholderErr(err error) bool {
	return errors.Is(err, cache.ErrPlaceholder)
}
