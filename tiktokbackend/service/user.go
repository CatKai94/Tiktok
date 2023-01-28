package service

import (
	"github.com/gin-gonic/gin"
	"tiktokbackend/models"
)

func SaveUser(c *gin.Context) {
	user := &models.User{
		Username: "zhangsan",
		Password: "123456",
	}
	models.SaveUser(user)
	c.JSON(200, user)
}
