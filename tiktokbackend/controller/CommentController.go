package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"tiktokbackend/models"
	"tiktokbackend/service"
	"tiktokbackend/utils"
	"time"
)

// CommentListResponse
// 评论列表返回参数
type CommentListResponse struct {
	StatusCode  int32                 `json:"status_code"`
	StatusMsg   string                `json:"status_msg,omitempty"`
	CommentList []service.CommentInfo `json:"comment_list,omitempty"`
}

// CommentActionResponse
// 发表评论返回参数
type CommentActionResponse struct {
	StatusCode int32               `json:"status_code"`
	StatusMsg  string              `json:"status_msg,omitempty"`
	Comment    service.CommentInfo `json:"comment"`
}

// CommentAction /comment/action/ - 评论操作
// 登录用户对视频进行评论。
func CommentAction(c *gin.Context) {
	// userId
	// 中间件jwt 设置的userID
	id, e := c.Get("userId")
	if !e {
		fmt.Println("获取userId发生错误：", e)
	}

	newId := id.(string)
	userId, err := strconv.ParseInt(newId, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "userId error",
		})
		return
	}
	// videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "videoId error",
		})
		return
	}
	// actionType
	actionType, err := strconv.ParseInt(c.Query("action_type"), 10, 64)
	log.Println("对评论的操作是：", actionType)

	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "actionType error",
		})
		return
	}

	serviceImpl := new(service.CommentServiceImpl)
	// 判断是为发布评论还是删除评论
	if actionType == 1 { // 发布评论
		text := c.Query("comment_text")
		// 敏感词过滤
		text = utils.Filter.Replace(text, '*')
		// 组装评论
		var SendComment models.Comment
		SendComment.UserId = userId
		SendComment.VideoId = videoId
		SendComment.CommentText = text
		SendComment.CreateDate = time.Now()

		comment, err := serviceImpl.SendComment(SendComment)
		if err != nil {
			c.JSON(http.StatusOK, Response{
				StatusCode: 1,
				StatusMsg:  "SendComment error",
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "send success",
			Comment:    comment,
		})
	} else if actionType == 2 { // 删除评论
		log.Println("controller层接收到的是删除评论命令")
		commentId, err := strconv.ParseInt(c.Query("comment_id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "delete commentId invalid",
			})
			return
		}
		err = serviceImpl.DeleteComment(commentId)
		if err != nil {
			c.JSON(http.StatusOK, CommentActionResponse{
				StatusCode: 1,
				StatusMsg:  "delete fail",
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			StatusCode: 0,
			StatusMsg:  "delete success",
		})
		return
	}
}

// CommentList /comment/list/ - 视频评论列表
// 查看视频的所有评论，按发布时间倒序。
func CommentList(c *gin.Context) {
	log.Println("Controller层CommentList")
	id, _ := c.Get("userId")
	userid, _ := id.(string)
	userId, err := strconv.ParseInt(userid, 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "userId error",
		})
		return
	}
	// videoId
	videoId, err := strconv.ParseInt(c.Query("video_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "videoId error",
		})
		return
	}
	serviceIpml := new(service.CommentServiceImpl)
	list, err := serviceIpml.GetCommentsList(videoId, userId)
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "GetCommentsList error",
		})
		return
	}
	c.JSON(http.StatusOK, CommentListResponse{
		StatusCode:  0,
		StatusMsg:   "GetCommentsList success",
		CommentList: list,
	})
}
