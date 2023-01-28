package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type response struct {
	// 必须大写才能序列化
	Message string `json:"responsemsg"`
	Code    int64  `json:"responseCode"`
}

func Test(c *gin.Context) {
	c.JSON(200, response{
		Message: "hello, gin !",
		Code:    200,
	})

	fmt.Println("test connect")
}
