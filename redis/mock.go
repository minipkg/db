package redis

import (
	//"github.com/alicebob/miniredis"
	"github.com/elliotchance/redismock/v8"
	//gocache "github.com/go-redis/cache/v8"
	//"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
)

// New creates a new mock client.
func NewMock() (*DB, *redismock.ClientMock, error) {
	//mr, err := miniredis.Run()
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//client := redis.NewClient(&redis.Options{
	//	Addr: mr.Addr(),
	//})
	//
	//mock := redismock.NewNiceMock(client)
	//dbobj := &DB{
	//	Exec: mock,
	//	cache: gocache.New(&gocache.Options{
	//		Redis: mock,
	//	}),
	//}
	//return dbobj, mock, nil
	return nil, nil, errors.New("This code under construction.")
}
