package redis

import (
	"context"

	gocache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/minipkg/db/redis/cache"
)

type IDB interface {
	DB() redis.Cmdable
	Close() error
	CacheSet(cacheItem *cache.Item) error
	CacheOnce(cacheItem *cache.Item) error
	CacheGet(ctx context.Context, key string, value interface{}) error
	CacheStats() *cache.Stats
}

type DB struct {
	client redis.UniversalClient
	cache  *gocache.Cache
	Exec   redis.Cmdable
}

var _ IDB = (*DB)(nil)

type Config struct {
	Addrs    []string
	Login    string
	Password string
	DBNum    int
}

// New creates a new DB connection
func New(conf Config) (*DB, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Addrs,
		Username: conf.Login,
		Password: conf.Password,
		DB:       conf.DBNum,
	})
	// @todo: try before timeout
	err := client.Ping(context.TODO()).Err()

	if err != nil {
		return nil, err
	}

	dbobj := &DB{
		client: client,
		Exec:   client,
		cache: gocache.New(&gocache.Options{
			Redis: client,
		}),
	}
	return dbobj, nil
}

// Close closes a client
func (d *DB) Close() error {

	if d.client != nil {
		if err := d.client.Close(); err != nil {
			return err
		}
	}
	return nil
}

// DB returns an object for execution commands
func (d *DB) DB() redis.Cmdable {
	return d.Exec
}

// CacheOnceItem makes a cache
// Once gets the item.Object for the given item.Key from the cache or executes, caches, and returns the results of the given item.Func, making sure that only one execution is in-flight for a given item.Key at a time. If a duplicate comes in, the duplicate caller waits for the original to complete and receives the same results.
// Example:
//	CacheOnceItem(&cache.Item{
//		Key:   "mykey",
//		Value: obj, // destination
//		Do: func(*cache.Item) (interface{}, error) {
//			return &Object{
//				Str: "mystring",
//				Num: 42,
//			}, nil
//		},
//	})
func (d *DB) CacheOnce(cacheItem *cache.Item) error {
	return d.cache.Once(cache.Item2CacheItem(cacheItem))
}

// CacheSet sets a cache
// Example:
//    ctx := context.Background()
//    key := "mykey"
//    obj := &Object{
//        Str: "mystring",
//        Num: 42,
//    }
//
//    if err := mycache.Set(&cache.Item{
//        Ctx:   ctx,
//        Key:   key,
//        Value: obj,
//        TTL:   time.Hour,
//    }); err != nil {
//        panic(err)
//    }
func (d *DB) CacheSet(cacheItem *cache.Item) error {
	return d.cache.Set(cache.Item2CacheItem(cacheItem))
}

// CacheGet gets a cached value
// Example:
//    var wanted Object
//    if err := mycache.Get(ctx, key, &wanted); err == nil {
//        fmt.Println(wanted)
//    }
func (d *DB) CacheGet(ctx context.Context, key string, value interface{}) error {
	return d.cache.Get(ctx, key, value)
}

func (d *DB) CacheStats() *cache.Stats {
	return cache.CacheStats2Stats(d.cache.Stats())
}
