package routers

import (
	"day1/controller"
	"day1/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	swggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func CollectRouter(r *gin.Engine) *gin.Engine {
	r.GET("/ip", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, ctx.ClientIP())
	})

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(200, "ok")
	})
	// 不需要认证的
	r.POST("/register", controller.Register)
	r.POST("/login", controller.Login)
	// 需要认证的
	r.GET("/userinfo", middleware.AuthMiddleWare(), controller.Userinfo)

	device := r.Group("/device")
	{
		device.GET("/devices", controller.ListDevice)
		device.GET("/device_datas", controller.GetDeviceData)
	}

	// // 文件上传下载
	// r.POST("/upload", controller.UploadFile)
	// r.GET("/download", controller.DownloadFile)
	// // mqtt
	// r.POST("/publish", controller.PublicMsg)
	// r.POST("/listen", controller.Listen)

	software := r.Group("/software", middleware.AuthMiddleWare())
	{
		software.POST("/add", controller.SoftwareAdd)
		software.POST("/update", controller.SoftwareUpdate)
		software.GET("/get", controller.SoftwareGet)
		software.GET("/list", controller.SoftwareList)
		software.POST("/del", controller.SoftwareDel)
		softwareImage := software.Group("/image")
		{
			softwareImage.POST("/add", controller.SoftwareImageAdd)
			softwareImage.POST("/update", controller.SoftwareUploadImage)
		}
		softwareFile := software.Group("/file")
		{
			softwareFile.POST("/add", controller.SoftwareFileAdd)
			softwareFile.POST("/update", controller.SoftwareUploadFile)
		}
		softwareArticle := software.Group("/article")
		{
			softwareArticle.POST("/add", controller.ArticleAdd)
			softwareArticle.POST("/update", controller.ArticleUpdate)
		}
	}

	r.GET("swagger/*any", ginSwagger.WrapHandler(swggerfiles.Handler))
	return r
}
