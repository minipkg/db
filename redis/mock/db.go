package mock

import (
	"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock/v8"
	"github.com/go-redis/redis/v8"

	redisdb "github.com/minipkg/db/redis"
)

// New creates a new mock client
func New() (*redisdb.DB, *redismock.ClientMock, error) {
	mr, err := miniredis.Run()
	if err != nil {
		return nil, nil, err
	}

	client := redis.NewClient(&redis.Options{
		Addr: mr.Addr(),
	})

	mock := redismock.NewNiceMock(client)
	dbobj := &redisdb.DB{Exec: mock}
	return dbobj, mock, nil
}
