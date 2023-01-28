package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"tiktokbackend/router"
	"tiktokbackend/utils"
)

func main() {
	// 解析配置文件
	cfg, err := utils.ParseConfig("./config/engineConfig.json")
	if err != nil {
		panic(err.Error())
	}

	engine := gin.Default()
	router.InitRouter(engine)
	engine.Use(Cors()) // 设置跨域
	engine.Run(cfg.EngineHost + ":" + cfg.EnginePort)
}

// 跨域访问
func Cors() gin.HandlerFunc {
	return func(context *gin.Context) {
		method := context.Request.Method
		origin := context.Request.Header.Get("Origin")

		var headerKeys []string
		for key, _ := range context.Request.Header {
			headerKeys = append(headerKeys, key)
		}

		headerStr := strings.Join(headerKeys, ",")
		if headerStr != "" {
			headerStr = fmt.Sprintf("access-control-allow-origin, access-control-allow-headers, %s", headerStr)
		} else {
			headerStr = "access-control-allow-origin, access-control-allow-headers"
		}

		if origin != "" {
			context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			context.Writer.Header().Set("Access-Control-Allow-Methods", "*")
			context.Writer.Header().Set("Access-Control-Allow-Headers", "*")
			context.Writer.Header().Set("Access-Control-Max-Age", "3600")
			context.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			//context.Set("content-type", "application/json") // 设置返回格式是json
		}

		if method == "OPTIONS" {
			context.JSON(http.StatusOK, "Options Request!")
		} else {
			// 处理请求
			context.Next()
		}

	}
}
