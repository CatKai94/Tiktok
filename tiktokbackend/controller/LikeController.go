package controller

import (
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/service"
)

// LikeActionResponse 点赞或取消点赞操作的返回体
type LikeActionResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

type LikeListResponse struct {
	StatusCode int32              `json:"status_code"`
	StatusMsg  string             `json:"status_msg"`
	VideoList  []service.FmtVideo `json:"video_list"`
}

// LikeAction  /favorite/action/ - 赞操作
func LikeAction(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	videoId, _ := strconv.ParseInt(c.Query("video_id"), 10, 64)
	actionType, _ := strconv.ParseInt(c.Query("action_type"), 10, 32)

	likeService := new(service.LikeServiceImpl)
	err := likeService.LikeAction(userId, videoId, int32(actionType))
	if err != nil {
		log.Println("service层方法LikeAction失败", err)
		c.JSON(http.StatusOK, LikeActionResponse{
			StatusCode: 1,
			StatusMsg:  "点赞或取消赞失败",
		})
	}
	log.Println("service层方法LikeAction成功")
	c.JSON(http.StatusOK, LikeActionResponse{
		StatusCode: 0,
		StatusMsg:  "点赞或取消赞成功",
	})
}

// GetLikeVideoList /favorite/list/ - 喜欢列表
// 登录用户的所有点赞视频。
func GetLikeVideoList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Println("当前用户tokenID:  ", curId, "!!!!!!!!!!!")
	if curId == 0 {
		c.JSON(http.StatusOK, LikeListResponse{
			StatusCode: 1,
			StatusMsg:  "用户未登录",
		})
	}

	likeService := new(service.LikeServiceImpl)
	videoList, err := likeService.GetLikeVideoList(userId, curId)
	if err != nil {
		log.Println("service层方法GetLikeVideoList失败", err)
		c.JSON(http.StatusOK, LikeListResponse{
			StatusCode: 1,
			StatusMsg:  "获取喜欢视频列表失败",
		})
	}
	log.Println("service层方法GetLikeVideoList成功")
	c.JSON(http.StatusOK, LikeListResponse{
		StatusCode: 0,
		StatusMsg:  "获取喜欢视频列表成功",
		VideoList:  videoList,
	})

}
