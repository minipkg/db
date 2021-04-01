package cache

import (
	"context"
	"time"
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
