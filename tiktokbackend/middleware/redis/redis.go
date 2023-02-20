package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"tiktokbackend/config"
)

var Ctx = context.Background()

var RdbUserToVideo *redis.Client
var RdbVideoToUser *redis.Client
var RdbFollowings *redis.Client
var RdbFollowers *redis.Client
var RdbFriends *redis.Client

func InitRedis() {
	RdbUserToVideo = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       0,
	})
	RdbVideoToUser = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       1,
	})

	RdbFollowings = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       2,
	})

	RdbFollowers = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       3,
	})

	RdbFriends = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       4,
	})

}
