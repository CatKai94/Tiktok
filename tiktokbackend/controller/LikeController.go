package controller

import (
	"github.com/gin-gonic/gin"
	"log"
)

// FavoriteAction /favorite/action/ - 赞操作
func FavoriteAction(c *gin.Context) {
	log.Println("Controller层FavoriteAction")
}

// GetFavouriteList /favorite/list/ - 喜欢列表
// 登录用户的所有点赞视频。
func GetFavouriteList(c *gin.Context) {
	log.Println("Controller层GetFavouriteList")
}
