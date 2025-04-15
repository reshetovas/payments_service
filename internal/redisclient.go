package internal

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client
var Ctx = context.Background()

func InitRedis() error {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	// declare a global variable - connect to Redis
	Rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := Rdb.Ping(Ctx).Result()
	return err
}
