package main

import (
	"day1/config"
	"day1/database"
	"day1/routers"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// 初始化配置
	config.InitConfig()
	// 初始化数据库
	database.InitMysql()
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
