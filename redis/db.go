package redis

import (
	"context"
	"time"

	gocache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"github.com/minipkg/db/redis/cache"
)

type IDB interface {
	DB() redis.Cmdable
	Close() error
	CacheOnce(cacheItem *cache.Item) error
	Cache(ctx context.Context, key string, value interface{}, funcToGetData func(*cache.Item) (interface{}, error), ttl time.Duration) error
}

type DB struct {
	client redis.UniversalClient
	cache  *gocache.Cache
	Exec   redis.Cmdable
}

var _ IDB = (*DB)(nil)
var _ cache.DB = (*DB)(nil)

type Config struct {
	Addrs    []string
	Login    string
	Password string
	DBName   int
}

// New creates a new DB connection
func New(conf Config) (*DB, error) {
	client := redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    conf.Addrs,
		Username: conf.Login,
		Password: conf.Password,
		DB:       conf.DBName,
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

func (d *DB) DB() redis.Cmdable {
	return d.Exec
}

//CacheOnce makes a cache
//Example:
//	CacheOnce(&gocache.Item{
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
	item := d.item2cacheItem(cacheItem)

	return d.cache.Once(item)
}

func (d *DB) Cache(ctx context.Context, key string, value interface{}, funcToGetData func(*cache.Item) (interface{}, error), ttl time.Duration) error {
	item := d.item2cacheItem(&cache.Item{
		Ctx:            ctx,
		Key:            key,
		Value:          value,
		TTL:            ttl,
		Do:             funcToGetData,
		SetXX:          false,
		SetNX:          false,
		SkipLocalCache: true,
	})
	return d.cache.Once(item)
}

func (d *DB) item2cacheItem(item *cache.Item) *gocache.Item {
	i := &gocache.Item{
		Ctx:            item.Ctx,
		Key:            item.Key,
		Value:          item.Value,
		TTL:            item.TTL,
		SetXX:          item.SetXX,
		SetNX:          item.SetNX,
		SkipLocalCache: item.SkipLocalCache,
	}
	i.Do = func(i *gocache.Item) (interface{}, error) {
		return item.Do(item)
	}
	return i
}

func (d *DB) Close() error {

	if d.client != nil {
		if err := d.client.Close(); err != nil {
			return err
		}
	}
	return nil
}
