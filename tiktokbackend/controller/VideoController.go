package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/service"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []service.FmtVideo `json:"video_list"`
	NextTime  int64              `json:"next_time"`
}

// Feed /feed/
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	log.Printf("传入的时间: " + inputTime)
	var lastTime time.Time
	if inputTime != "0" {
		//me, _ := strconv.ParseInt(inputTime, 10, 64)
		//lastTime = time.Unix(me, 0)
		log.Println("inputTime != 0")
		// 调试用
		lastTime = time.Now()
	} else {
		log.Println("inputTime == 0")
		lastTime = time.Now()
	}
	log.Printf("获取到时间戳%v", lastTime)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	videoService := GetVideo()

	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	log.Printf("方法videoService.Feed(lastTime, userId) 成功")

	// 打印结果
	fmt.Println("FmtVideo：", feed)
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

// GetVideo 拼装videoService
func GetVideo() service.VideoServiceImpl {
	var userService service.UserServiceImpl
	//var followService service.FollowServiceImp
	var videoService service.VideoServiceImpl
	//var likeService service.LikeServiceImpl
	//var commentService service.CommentServiceImpl
	//userService.FollowService = &followService
	//userService.LikeService = &likeService
	//followService.UserService = &userService
	//likeService.VideoService = &videoService
	//commentService.UserService = &userService
	//videoService.CommentService = &commentService
	//videoService.LikeService = &likeService
	videoService.UserService = &userService
	return videoService
}
