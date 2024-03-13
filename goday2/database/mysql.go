package database

import (
	"day1/model"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMysql() *gorm.DB {
	// 获取配置
	host := viper.GetString("mysql.host")
	port := viper.GetString("mysql.port")
	database := viper.GetString("mysql.database")
	username := viper.GetString("mysql.username")
	password := viper.GetString("mysql.password")
	// 格式化
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&loc=Local",
		username, password, host, port, database)
	// mysql连接
	db, err := gorm.Open(mysql.Open(args), &gorm.Config{})
	if err != nil {
		panic("failed to connect mysql database,err:" + err.Error())
	}
	// 自动建表
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.Software{})
	db.AutoMigrate(&model.File{})
	db.AutoMigrate(&model.Image{})
	db.AutoMigrate(&model.Article{})
	db.AutoMigrate(&model.Device{})
	db.AutoMigrate(&model.DeviceData{})
	// 全局变量
	DB = db
	return db
}

// 返回mysql连接
func GetDB() *gorm.DB {
	return DB
}
