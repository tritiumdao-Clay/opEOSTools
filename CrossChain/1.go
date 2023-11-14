package main

import (
	"github.com/gin-gonic/gin"
	// 导入cors包
	"github.com/gin-contrib/cors"
)

func main() {
	r := gin.Default()

	// 配置跨域中间件
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:10003"}                   // 允许什么域名访问，支持多个域名
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}  // 允许的 HTTP 方法
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type"} // 允许的 HTTP 头
	// 设置cors中间件
	r.Use(cors.New(config))

	// 测试api
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run() // 启动服务
}
