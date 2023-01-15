package service

import (
	"github.com/gin-gonic/gin"
	"tiktokbackend/models"
	"time"
)

func SaveUser(c *gin.Context) {
	user := &models.User{
		Username:   "zhangsan",
		Password:   "123456",
		CreateTime: time.Now().UnixMilli(),
	}
	models.SaveUser(user)
	c.JSON(200, user)
}

func GetUser(context *gin.Context) {
	user := models.GetById(4)
	context.JSON(200, user)
}

func UpdateUser(context *gin.Context) {
	models.UpdateUser(1)
	user := models.GetById(1)
	context.JSON(200, user)
}

func DeleteUser(context *gin.Context) {
	models.DeleteUser(1)
	user := models.GetById(1)
	context.JSON(200, user)
}
