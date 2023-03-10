package router

import (
	"github.com/gin-gonic/gin"
	"tiktokbackend/controller"
	"tiktokbackend/utils"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	/*
		基础接口
	*/
	//视频流接口
	apiRouter.GET("/feed/", utils.AuthWithoutLogin(), controller.Feed)
	//用户注册
	apiRouter.POST("/user/register/", controller.Register)
	//用户登录
	apiRouter.POST("/user/login/", controller.Login)
	//用户信息
	apiRouter.GET("/user/", controller.UserInfo)
	//视频投稿
	apiRouter.POST("/publish/action/", utils.AuthPostValue(), controller.PublishAction)
	//发布列表
	apiRouter.GET("/publish/list/", utils.Auth(), controller.PublishList)

	/*
		互动接口
	*/
	//赞操作
	apiRouter.POST("/favorite/action/", utils.Auth(), controller.LikeAction)
	//喜欢列表
	apiRouter.GET("/favorite/list/", utils.Auth(), controller.GetLikeVideoList)
	//评论操作
	apiRouter.POST("/comment/action/", utils.Auth(), controller.CommentAction)
	//视频评论列表
	apiRouter.GET("/comment/list/", utils.AuthWithoutLogin(), controller.CommentList)

	/*
		社交接口
	*/
	//关系操作
	apiRouter.POST("/relation/action/", utils.Auth(), controller.RelationAction)
	//用户关注列表
	apiRouter.GET("/relation/follow/list/", utils.Auth(), controller.GetFollowingList)
	//用户粉丝列表
	apiRouter.GET("/relation/follower/list/", utils.Auth(), controller.GetFollowerList)
	//用户好友列表
	apiRouter.GET("relation/friend/list/", utils.Auth(), controller.GetFriendList)
	//消息操作
	apiRouter.GET("/message/chat/", utils.Auth(), controller.MessageChat)
	//聊天记录
	apiRouter.POST("/message/action/", utils.Auth(), controller.MessageAction)
}
