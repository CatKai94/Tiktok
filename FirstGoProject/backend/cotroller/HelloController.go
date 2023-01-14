package cotroller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type HelloController struct {
}

func Hello(context *gin.Context) {
	fmt.Println("get请求正常")

	context.JSON(200, map[string]interface{}{
		"message": "hello gin！ ！ ！",
	})

	fmt.Println("返回get信息给前端")
}
