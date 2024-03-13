package middleware

import (
	"day1/database"
	"day1/model"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头拿到Authorization
		tokenString := ctx.GetHeader("Authorization")
		// 看前缀Bearer xxxxxxx
		if tokenString == "" || !strings.HasPrefix(tokenString, "Bearer") {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": "401",
				"data": "token验证失败",
				"msg":  "token验证失败",
			})
			ctx.Abort()
			return
		}
		// 切掉前七个 Bearer空格 6+1
		tokenString = tokenString[7:]
		// 解token
		token, claims, err := ParseToken(tokenString)
		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code": "401",
				"data": "token失效",
				"msg":  "token失效",
			})
			ctx.Abort()
			return
		}
		// 从claims中拿到userID
		userID := claims.UserID
		// 拿到数据库
		DB := database.GetDB()
		var user model.User

		DB.Table("users").Where("id = ?", userID).First(&user)
		if user.ID == 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"code": "401",
				"data": "用户不存在",
				"msg":  "用户不存在",
			})
			ctx.Abort()
			return
		}
		// 把某些东西放到上下文
		ctx.Set("user", user)
		ctx.Set("userid", user.ID)
		ctx.Next()
	}
}
