package controller

import (
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
	lastTime = time.Now()

	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	//videoService := GetVideo()
	videoService := service.VideoServiceImpl{}

	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

// Publish /publish/action/
func Publish(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Println("发布视频的用户id为", userId)

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

	videoService := service.VideoServiceImpl{}
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

	videoService := service.VideoServiceImpl{}

	// 获取用户所发布视频的列表
	list, err := videoService.List(userId, curId)
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

