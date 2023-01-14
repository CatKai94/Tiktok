package cotroller

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type TestpostController struct {
}

func (testpost *TestpostController) Router(engine *gin.Engine) {
	engine.POST("/testpost", testpost.Testpost)
}

// 解析 /hello
func (testpost *TestpostController) Testpost(context *gin.Context) {
	fmt.Println("post请求正常")

	context.JSON(200, map[string]interface{}{
		"message": "TestPost！ ！ ！",
	})

	fmt.Println("返回post信息给前端")
}
