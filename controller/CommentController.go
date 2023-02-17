package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

// CommentAction /comment/action/ - 评论操作
// 登录用户对视频进行评论。
func CommentAction(c *gin.Context) {
	log.Println("Controller层CommentAction")
}

// CommentList /comment/list/ - 视频评论列表
// 查看视频的所有评论，按发布时间倒序。
func CommentList(c *gin.Context) {
	log.Println("Controller层CommentList")
}
