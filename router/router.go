package router

import (
	"github.com/gin-gonic/gin"
	"tiktokbackend/controller"
)

func InitRouter(r *gin.Engine) {
	// public directory is used to serve static resources
	r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	//apiRouter.GET("/test", controller.Test)

	/*
		基础接口
	*/
	//视频流接口
	apiRouter.GET("/feed/", controller.Feed)
	//用户注册
	apiRouter.POST("/user/register/", controller.Register)
	//用户登录
	apiRouter.POST("/user/login/", controller.Login)
	//用户信息
	apiRouter.GET("/user/", controller.UserInfo)
	//视频投稿
	apiRouter.POST("/publish/action/", controller.Publih)
	//发布列表
	apiRouter.GET("/publish/list/", controller.PublishList)

	/*
		互动接口
	*/
	//赞操作
	apiRouter.POST("/favorite/action/", controller.LikeAction)
	//喜欢列表
	apiRouter.GET("/favorite/list/", controller.GetLikeVideoList)
	//评论操作
	apiRouter.POST("/comment/action/", nil)
	//视频评论列表
	apiRouter.GET("/comment/list/", nil)

	/*
		社交接口
	*/
	//关系操作
	apiRouter.POST("/relation/action/", nil)
	//用户关注列表
	apiRouter.GET("/relation/follow/list/", nil)
	//用户粉丝列表
	apiRouter.GET("/relation/follower/list/", nil)
	//用户好友列表
	apiRouter.GET("relation/friend/list/", nil)
	//消息操作
	apiRouter.POST("/message/chat/", nil)
	//聊天记录
	apiRouter.GET("/message/action/", nil)
}
