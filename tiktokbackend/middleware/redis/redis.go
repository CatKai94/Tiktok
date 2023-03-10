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

// RdbVideoToCommentId 视频ID对应多个评论ID
var RdbVideoToCommentId *redis.Client

// RdbCommentToVideoId 评论ID对应一个视频ID
var RdbCommentToVideoId *redis.Client

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
	
	// 评论
	RdbVideoToCommentId = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       5,
	})
	RdbCommentToVideoId = redis.NewClient(&redis.Options{
		Addr:     config.RedisUrl + ":6379",
		Password: "",
		DB:       6,
	})

}
