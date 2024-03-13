package routers

import (
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
	return r
}
