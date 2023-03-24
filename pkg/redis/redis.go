package redis

import (
	"github.com/go-redis/redis/v8"
)

func NewRedisClient(address string, password string, dbNum int) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       dbNum,
	})
	return rdb

}
