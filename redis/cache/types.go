package cache

import (
	"context"
	"time"

	"github.com/go-redis/cache/v9"
)

type Item struct {
	Ctx context.Context

	Key   string
	Value interface{}

	// TTL is the cache expiration time.
	// Default TTL is 1 hour.
	TTL time.Duration

	// Do returns value to be cached.
	Do func(*Item) (interface{}, error)

	// SetXX only sets the key if it already exists.
	SetXX bool

	// SetNX only sets the key if it does not already exist.
	SetNX bool

	// SkipLocalCache skips local cache as if it is not set.
	SkipLocalCache bool
}

func Item2CacheItem(item *Item) *cache.Item {
	i := &cache.Item{
		Ctx:            item.Ctx,
		Key:            item.Key,
		Value:          item.Value,
		TTL:            item.TTL,
		SetXX:          item.SetXX,
		SetNX:          item.SetNX,
		SkipLocalCache: item.SkipLocalCache,
	}
	i.Do = func(i *cache.Item) (interface{}, error) {
		return item.Do(item)
	}
	return i
}

type Stats struct {
	Hits   uint64
	Misses uint64
}

func CacheStats2Stats(cacheStats *cache.Stats) *Stats {
	return &Stats{
		Hits:   cacheStats.Hits,
		Misses: cacheStats.Misses,
	}
}
