package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

// RelationAction /relation/action/ - 关系操作
// 登录用户对其他用户进行关注或取消关注。
func RelationAction(c *gin.Context) {
	log.Println("Controller层RelationAction")
}

// GetFollowingList /relation/follow/list/ - 用户关注列表
// 登录用户关注的所有用户列表。
func GetFollowingList(c *gin.Context) {
	log.Println("Controller层GetFollowingList")
}

// GetFollowerList /relation/follower/list/ - 用户粉丝列表
// 所有关注登录用户的粉丝列表。
func GetFollowerList(c *gin.Context) {
	log.Println("Controller层GetFollowerList")
}

// GetFriendList /relation/friend/list/ - 用户好友列表
// 所有和用户互关的粉丝列表
func GetFriendList(c *gin.Context) {
	log.Println("Controller层GetFriendList")
}
