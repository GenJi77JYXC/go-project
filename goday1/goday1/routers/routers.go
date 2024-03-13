package routers

import (
	"day1/controller"
	"day1/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
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
	return r
}
