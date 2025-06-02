package main

import "github.com/redis/go-redis/v9"

func RedisInit() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     LoadConfig().REDIS_HOST,
		Password: LoadConfig().REDIS_PASSWORD,
	})
	return rdb
}
