package config

import (
	"os"

	"github.com/spf13/viper"
)

func InitConfig() {
	//main 的工作目录
	workDir, _ := os.Getwd()
	// 配置文件名字
	viper.SetConfigName("config")
	// 配置文件类型
	viper.SetConfigType("yml")
	// 配置文件路径
	viper.AddConfigPath(workDir + "/config")
	// 尝试读入配置文件
	err := viper.ReadInConfig()
	// 读取失败报错退出
	if err != nil {
		panic(err)
	}
}
