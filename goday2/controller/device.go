package controller

import (
	"day1/database"
	"day1/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func ListDevice(ctx *gin.Context) {
	DB := database.GetDB()
	var devices []model.Device
	DB.Table("devices").Find(&devices)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": devices,
		"msg":  "获取设备成功",
	})
}

func GetDeviceData(ctx *gin.Context) {
	id := ctx.Query("id")
	DB := database.GetDB()
	var datas []model.DeviceData
	DB.Table("device_data").Where("mac_id = ?", id).Find(&datas)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": datas,
		"msg":  "获取设备数据成功",
	})
}

func GetDeviceDataByTime(ctx *gin.Context) {
	id := ctx.Query("id")
	dateStart := ctx.Query("start")
	dateEnd := ctx.Query("end")
	start, _ := time.Parse("2006-01-02", dateStart)
	end, _ := time.Parse("2006-01-02", dateEnd)
	DB := database.GetDB()
	var datas []model.DeviceData
	DB.Table("device_data").Where("mac_id = ?", id).Where("created_at >= ?", start).Where("created_at <= ?", end).Find(&datas)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": datas,
		"msg":  "安装日期获取设备数据成功",
	})
}
