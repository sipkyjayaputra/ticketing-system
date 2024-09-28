package configuration

import (
	"log"

	"github.com/go-redis/redis"
)

func ConnectRedis() (*redis.Client, error) {
	cache := redis.NewClient(&redis.Options{
		Addr: CONFIG["REDIS_HOST"],
	})

	log.Println("Redis connection success")

	return cache, nil
}
