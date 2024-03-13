package main

import (
	"goday3/config"
	"goday3/routers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server celler server.
// @host 47.108.235.139:8080
// @BasePath /api/
func main() {
	// 初始化配置
	config.InitConfig()

	// docs.SwaggerInfo.BasePath("/api")
	// 创建一个默认的gin路由
	r := gin.Default()
	// 路由绑定
	r = routers.CollectRouter(r)
	// 从viper获取到运行端口
	port := viper.GetString("server.port")
	// 如果指定了端口
	if port != "" {
		panic(r.Run(":" + port))
	}
	// 没指定端口就直接运行
	panic(r.Run())
}
