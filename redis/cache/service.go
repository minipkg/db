package cache

import (
	"context"
	"fmt"
	"time"
)

const (
	keyPrefixForCache = "cache_"
)

type Service interface {
	Cache(ctx context.Context, key string, value interface{}, funcToGetData func(*Item) (interface{}, error)) error
	CacheOnce(cacheItem *Item) error
}

var _ Service = (*service)(nil)

type service struct {
	db  DB
	ttl time.Duration
}

type DB interface {
	CacheOnce(cacheItem *Item) error
	Cache(ctx context.Context, key string, value interface{}, funcToGetData func(*Item) (interface{}, error), ttl time.Duration) error
}

// NewService creates a new authentication service.
func NewService(db DB, ttlInHour uint) *service {
	return &service{
		db:  db,
		ttl: time.Duration(int64(ttlInHour)) * time.Hour,
	}
}

func (s *service) CacheOnce(cacheItem *Item) error {
	(*cacheItem).Key = s.key((*cacheItem).Key)
	(*cacheItem).TTL = s.ttl
	return s.db.CacheOnce(cacheItem)
}

func (s *service) Cache(ctx context.Context, key string, value interface{}, funcToGetData func(*Item) (interface{}, error)) error {
	return s.db.Cache(ctx, s.key(key), value, funcToGetData, s.ttl)
}

func (s *service) key(key string) string {
	return fmt.Sprintf("%s%s", keyPrefixForCache, key)
}
