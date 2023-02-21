package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/service"
)

// FollowingResp 获取关注列表需要返回的结构。
type FollowingResp struct {
	Response
	UserList []service.FmtUser `json:"user_list,omitempty"`
}

// FollowersResp 获取粉丝列表需要返回的结构。
type FollowersResp struct {
	Response
	UserList []service.FmtUser `json:"user_list,omitempty"`
}

// RelationAction /relation/action/ - 关系操作
// 登录用户对其他用户进行关注或取消关注。
func RelationAction(c *gin.Context) {
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	toUserId, _ := strconv.ParseInt(c.Query("to_user_id"), 10, 64)
	actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 64)

	followService := new(service.FollowServiceImp)
	var err error

	if actionType == 1 {
		err = followService.FollowAction(toUserId, curId)
	} else if actionType == 2 {
		err = followService.UnFollowAction(toUserId, curId)
	}

	if err != nil {
		log.Println("关注或取消关注操作发生错误：", err)
		c.JSON(http.StatusOK, LikeActionResponse{
			StatusCode: 1,
			StatusMsg:  "点赞或取消赞失败",
		})
		return
	}

	log.Println("关注、取关成功。")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "OK",
	})
}

// GetFollowingList /relation/follow/list/ - 用户关注列表
// 登录用户关注的所有用户列表。
func GetFollowingList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	followService := new(service.FollowServiceImp)

	users, err := followService.GetFollowingsList(userId)
	// 获取关注列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表时出错。",
			},
			UserList: nil,
		})
		return
	}
	// 成功获取到关注列表。
	log.Println("获取关注列表成功。")
	c.JSON(http.StatusOK, FollowingResp{
		UserList: users,
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}

// GetFollowerList /relation/follower/list/ - 用户粉丝列表
// 所有关注登录用户的粉丝列表。
func GetFollowerList(c *gin.Context) {
	userId, err := strconv.ParseInt(c.Query("user_id"), 10, 64)

	followService := new(service.FollowServiceImp)
	followers, err := followService.GetFollowersList(userId)
	if err != nil { // 获取粉丝列表时发生错误
		c.JSON(http.StatusOK, FollowersResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表时出错。",
			},
			UserList: nil,
		})
		return
	}
	c.JSON(http.StatusOK, FollowersResp{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		UserList: followers,
	})
}

// GetFriendList /relation/friend/list/ - 用户好友列表
// 所有和用户互关的粉丝列表
func GetFriendList(c *gin.Context) {
	log.Println("Controller层GetFriendList")
}
