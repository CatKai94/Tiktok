package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()

var RdbUserToVideo *redis.Client
var RdbVideoToUser *redis.Client

func InitRedis() {
	RdbUserToVideo = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	RdbVideoToUser = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       1,
	})
}