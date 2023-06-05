package common

import "github.com/redis/go-redis/v9"

var MyRedis *redis.Client

func RedisInit() *redis.Client {
	rds := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
	MyRedis = rds
	return rds
}
