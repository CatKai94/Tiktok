package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/service"
	"tiktokbackend/utils"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []service.FmtVideo `json:"video_list"`
	NextTime  int64              `json:"next_time"`
}

type PublishListResponse struct {
	Response
	VideoList []service.FmtVideo `json:"video_list"`
}

// Feed /feed/
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	log.Printf("传入的时间: " + inputTime)
	var lastTime time.Time
	if inputTime != "0" {
		log.Println("inputTime != 0")
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

// Publish /publish/action/
func Publih(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	title := c.PostForm("title")
	file, err := c.FormFile("data")

	if err != nil {
		log.Println("获取视频流失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "获取视频流失败",
		})
		return
	}
	fileName := utils.CreateFileName(userId) //根据userId和时间戳得到唯一的文件名
	videoPath := "./public/videos/" + fileName + ".mp4"

	err = c.SaveUploadedFile(file, videoPath)
	if err != nil {
		log.Println("保存视频时发生了错误: ", err)
		return
	}
	//截取一帧画面作为封面
	err = utils.SaveFaceImage(fileName)
	if err != nil {
		log.Println("截取视频封面失败")
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "截取视频封面失败",
		})
	}

	videoService := GetVideo()
	err = videoService.Publish(fileName, userId, title)

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "后端data入库失败",
		})
		return
	}

	log.Printf("成功发布视频")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "发布视频成功",
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	log.Println("被查看发布列表的用户id是：", userId)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Println("查看发布列表的用户id是： ", curId)

	vedioService := GetVideo()

	// 获取用户所发布视频的列表
	list, err := vedioService.List(userId, curId)
	// 获取用户所发布视频的列表 -- 失败
	if err != nil {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取用户所发布视频的列表失败"},
		})
		return
	}
	// 获取用户所发布视频的列表 -- 成功
	c.JSON(http.StatusOK, PublishListResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "获取用户所发布视频的列表成功"},
		VideoList: list,
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
