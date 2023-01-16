package controller

import "github.com/gin-gonic/gin"

// Login POST douyin/user/login/ 用户登录
func Login(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	str := username + password
	c.JSON(200, str)
}
