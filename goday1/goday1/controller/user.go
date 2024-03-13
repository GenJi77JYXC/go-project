package controller

import (
	"day1/database"
	"day1/middleware"
	"day1/model"
	"day1/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Register(ctx *gin.Context) {
	DB := database.GetDB()

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	//ctx.ShouldBindJSON()
	//验证
	// 用户名正确传入
	if len(username) == 0 {
		// ctx.JSON(http.StatusOK, gin.H{
		// 	"code": "400",
		// 	"data": "用户名错误",
		// 	"msg":  "用户名错误",
		// })
		response.Fail(ctx, 400, "用户名错误", "用户名错误")
		return
	}
	// 密码长度
	if len(password) < 6 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "400",
			"data": "密码需要六位以上",
			"msg":  "密码需要六位以上",
		})
		return
	}
	// 用户有没有被注册过
	if isUserExits(DB, username) {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "400",
			"data": "用户名已经被使用",
			"msg":  "用户名已经被使用",
		})
		return
	}
	// 开始注册流程
	hasedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "500",
			"data": "加密错误",
			"msg":  "加密错误",
		})
		return
	}
	newUser := model.User{
		Username: username,
		Password: string(hasedPassword),
	}
	DB.Create(&newUser)
	//
	var user model.User
	DB.Table("users").Where("username = ?", username).First(&user)
	// 发放token
	token, err := middleware.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "500",
			"data": "token加密错误",
			"msg":  "token加密错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"token": token,
		},
		"msg": "注册成功",
	})
}

func Login(ctx *gin.Context) {
	DB := database.GetDB()

	username := ctx.PostForm("username")
	password := ctx.PostForm("password")
	// 密码长度
	if len(password) < 6 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "400",
			"data": "密码需要六位以上",
			"msg":  "密码需要六位以上",
		})
		return
	}
	var user model.User
	DB.Table("users").Where("username = ?", username).First(&user)
	if user.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "400",
			"data": "用户不存在",
			"msg":  "用户不存在",
		})
		return
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "422",
			"data": "密码错误",
			"msg":  "密码错误",
		})
		return
	}
	// 发放token
	token, err := middleware.ReleaseToken(user)
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "500",
			"data": "token加密错误",
			"msg":  "token加密错误",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"token": token,
		},
		"msg": "登录成功",
	})
}

func Userinfo(ctx *gin.Context) {
	user, _ := ctx.Get("user")
	u := user.(model.User)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"ID":       u.ID,
			"Username": u.Username,
		},
		"msg": "获取成功",
	})
}

func isUserExits(db *gorm.DB, username string) bool {
	var user model.User
	db.Table("users").Where("username = ?", username).First(&user)
	return user.ID != 0
}
