package controller

import (
	"bytes"
	"context"
	"day1/database"
	"day1/model"
	"day1/response"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/spf13/viper"
)

// @summary 上传软件
// @description 软件描述
// @Accept json
// @Produce	json
// @Success 200
// @Router /software/add [post]
func SoftwareAdd(ctx *gin.Context) {
	var software model.Software
	err := ctx.ShouldBindJSON(&software)
	if err != nil {
		response.Fail(ctx, 400, "传参错误", "传参错误")
	}
	//数据验证
	// 软件已经存在
	DB := database.GetDB()
	var check model.Software
	DB.Table("softwares").Where("name = ?", software.Name).First(&check)
	if check.ID != 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "422",
			"data": "软件已经存在",
			"msg":  "软件已经存在",
		})
		return
	}

	DB.Create(&software)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "新增成功",
	})
}

// @summary 修改软件
// @description 修改信息
// @Accept json
// @Produce	json
// @Success 200
// @Router /software/update [post]
func SoftwareUpdate(ctx *gin.Context) {
	var software model.Software
	ctx.ShouldBindJSON(&software)
	//数据验证
	//软件已经存在
	DB := database.GetDB()
	var check model.Software
	DB.Table("softwares").Where("name = ?", software.Name).First(&check)
	if check.ID == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"code": "422",
			"data": "软件不存在",
			"msg":  "软件不存在",
		})
		return
	}

	DB.Table("softwares").Where("id = ?", check.ID).Update("description", software.Description)
	DB.Table("softwares").Where("id = ?", check.ID).Update("category", software.Category)
	DB.Table("softwares").Where("id = ?", check.ID).Update("tag", software.Tag)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "更新成功",
	})
}

// @summary 列出所有软件
// @description 列出软件
// @Accept json
// @Produce	json
// @Success 200
// @Router /software/get [get]
func SoftwareGet(ctx *gin.Context) {
	id := ctx.Query("id")
	DB := database.GetDB()

	// 声明software
	var software model.Software
	// 声明子集
	var oldImage model.Image
	var oldFile model.File
	var oldArticle model.Article
	// 拿到software
	DB.Table("softwares").Where("id = ?", id).First(&software)
	// 拿到子集
	DB.Table("images").Where("id = ?", software.ImageID).First(&oldImage)
	DB.Table("softwares").Where("id = ?", software.FileID).First(&oldFile)
	DB.Table("articles").Where("id = ?", software.ArticleID).First(&oldArticle)

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": gin.H{
			"software": software,
			"image":    oldImage,
			"file":     oldFile,
			"article":  oldArticle,
		},
		"msg": "获取成功",
	})

}

// @summary 删除软件
// @description 删除软件
// @Accept json
// @Produce	json
// @Success 200
// @Param id body string true "软件id"
// @Router /software/del [post]
func SoftwareDel(ctx *gin.Context) {
	id := ctx.PostForm("id")
	DB := database.GetDB()
	// 声明software
	var software model.Software
	// 声明子集
	var oldImage model.Image
	var oldFile model.File
	var oldArticle model.Article
	// 拿到software
	DB.Table("softwares").Where("id = ?", id).First(&software)
	// 拿到子集
	DB.Table("images").Where("id = ?", software.ImageID).First(&oldImage)
	DB.Table("softwares").Where("id = ?", software.FileID).First(&oldFile)
	DB.Table("softwares").Where("id = ?", software.ArticleID).First(&oldArticle)
	// 删除子集
	DB.Delete(&oldImage)
	DB.Delete(&oldFile)
	DB.Delete(&oldArticle)
	// 删除software
	DB.Delete(&software)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "删除成功",
	})
}

// @summary 列出软件
// @description 列出软件
// @Accept json
// @Produce	json
// @Success 200
// @Router /software/list [get]
func SoftwareList(ctx *gin.Context) {
	var softwares []model.Software
	DB := database.GetDB()
	DB.Table("softwares").Find(&softwares)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": softwares,
		"msg":  "查询成功",
	})
}

// @Summary 上传封面
// @Description 上传封面
// @Accept json
// @Produce json
// @Success 200
// @Router /software/image/add [post]
func SoftwareImageAdd(ctx *gin.Context) {
	bucketName := viper.GetString("minio.bucket.image")
	name := ctx.PostForm("name")
	id := ctx.PostForm("id")
	// 这里是从form表单里面拿到key为file的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("upload fail")
		return
	}
	// 临时保存
	ctx.SaveUploadedFile(file, "./upload/"+name)
	MC := database.GetMinioClient()
	bucket, err := MC.BucketExists(context.Background(), bucketName)
	// mc 调用buckets这个函数能不能正常执行
	if err != nil {
		fmt.Println("mc client fail")
		return
	}
	// bucket 存在
	if !bucket {
		fmt.Println("bucket not exist")
		return
	}
	uploadInfo, err := MC.FPutObject(context.Background(), bucketName, name, "./upload/"+name, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("mc upload fail")
		return
	}
	// 数据库
	DB := database.GetDB()
	newFile := &model.Image{
		Image:  name,
		Bucket: bucketName,
	}
	// 数据库打一条文件上传的记录
	DB.Table("iamges").Create(&newFile)
	// 从库里面重新拿到newFile，为了得到id
	DB.Table("images").Where("image = ?", name).First(&newFile)
	DB.Table("softwares").Where("id = ?", id).Update("image_id", newFile.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": uploadInfo,
		"msg":  "上传成功",
	})
}

// @Summary 上传文件
// @Description 上传文件
// @Accept json
// @Produce json
// @Success 200
// @Router /software/file/add [post]
func SoftwareFileAdd(ctx *gin.Context) {
	bucketName := viper.GetString("minio.bucket.file")
	name := ctx.PostForm("name")
	id := ctx.PostForm("id")
	// 这里是从form表单里面拿到key为file的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("upload fail")
		return
	}
	// 临时保存
	ctx.SaveUploadedFile(file, "./upload/"+name)
	MC := database.GetMinioClient()
	bucket, err := MC.BucketExists(context.Background(), bucketName)
	// mc 调用buckets这个函数能不能正常执行
	if err != nil {
		fmt.Println("mc client fail")
		return
	}
	// bucket 存在
	if !bucket {
		fmt.Println("bucket not exist")
		return
	}
	uploadInfo, err := MC.FPutObject(context.Background(), bucketName, name, "./upload/"+name, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("mc upload fail")
		return
	}
	// 数据库
	DB := database.GetDB()
	newFile := &model.File{
		FileName: name,
		Bucket:   bucketName,
	}
	// 数据库打一条文件上传的记录
	DB.Table("files").Create(&newFile)
	// 从库里面重新拿到newFile，为了得到id
	DB.Table("files").Where("file = ?", name).First(&newFile)
	DB.Table("softwares").Where("id = ?", id).Update("file_id", newFile.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": uploadInfo,
		"msg":  "上传成功",
	})
}

// @Summary 更新封面
// @Description 上传封面
// @Accept json
// @Produce json
// @Success 200
// @Router /software/image/update [post]
func SoftwareUploadImage(ctx *gin.Context) {
	bucketName := viper.GetString("minio.bucket.image")
	name := ctx.PostForm("name")
	id := ctx.PostForm("id")
	// 这里是从form表单里面拿到key为file的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("upload fail")
		return
	}
	//临时保存
	ctx.SaveUploadedFile(file, "./upload/"+name)
	MC := database.GetMinioClient()
	bucket, err := MC.BucketExists(context.Background(), bucketName)
	// mc调用buckets这个函数能不能正常执行
	if err != nil {
		fmt.Println("mc client fail ")
		return
	}
	// bucket 存在
	if !bucket {
		fmt.Println("bucket is not exist")
	}

	uploadInfo, err := MC.FPutObject(context.Background(), bucketName, name, "./upload/"+name, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("mc upload fail")
		return
	}
	//数据库
	DB := database.GetDB()

	var oldImage model.Image
	var software model.Software
	DB.Table("softwares").Where("id = ?", id).First(&software)
	DB.Table("images").Where("id = ?", software.ImageID).First(&oldImage)
	DB.Delete(&oldImage)

	newFile := &model.Image{
		Image:  name,
		Bucket: bucketName,
	}
	// 数据库打上一条文件上传记录
	DB.Table("images").Create(&newFile)
	// 从库里面重新拿到newFile， 为了得到id
	DB.Table("images").Where("image = ?", name).First(&newFile)
	DB.Table("softwares").Where("id = ?", id).Update("image_id", newFile.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": uploadInfo,
		"msg":  "上传成功",
	})
}

// @Summary 更新文件
// @Description 上传文件
// @Accept json
// @Produce json
// @Success 200
// @Router /software/file/update [post]
func SoftwareUploadFile(ctx *gin.Context) {
	bucketName := viper.GetString("minio.bucket.file")
	name := ctx.PostForm("name")
	id := ctx.PostForm("id")
	// 这里是从form表单里面拿到key为file的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println("upload fail")
		return
	}
	//临时保存
	ctx.SaveUploadedFile(file, "./upload/"+name)
	MC := database.GetMinioClient()
	bucket, err := MC.BucketExists(context.Background(), bucketName)
	// mc调用buckets这个函数能不能正常执行
	if err != nil {
		fmt.Println("mc client fail ")
		return
	}
	// bucket 存在
	if !bucket {
		fmt.Println("bucket is not exist")
	}

	uploadInfo, err := MC.FPutObject(context.Background(), bucketName, name, "./upload/"+name, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("mc upload fail")
		return
	}
	//数据库
	DB := database.GetDB()

	var oldFile model.File
	var software model.Software
	DB.Table("softwares").Where("id = ?", id).First(&software)
	DB.Table("files").Where("id = ?", software.ImageID).First(&oldFile)
	DB.Delete(&oldFile)

	newFile := &model.File{
		FileName: name,
		Bucket:   bucketName,
	}
	// 数据库打上一条文件上传记录
	DB.Table("files").Create(&newFile)
	// 从库里面重新拿到newFile，为了得到ID
	DB.Table("files").Where("image = ?", name).First(&newFile)
	DB.Table("softwares").Where("id = ?", id).Update("file_id", newFile.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": uploadInfo,
		"msg":  "更新文件成功",
	})
}

func ArticleAdd(ctx *gin.Context) {
	id := ctx.PostForm("id")
	content := ctx.PostForm("content")
	DB := database.GetDB()
	id_int, _ := strconv.Atoi(id)

	newArticle := model.Article{
		SoftwareID: uint(id_int),
		Content:    content,
	}
	DB.Create(&newArticle)
	var article model.Article
	//这里应该用softwareID反查查
	DB.Table("articles").Where("software_id = ?", id).First(&article)
	DB.Table("softwares").Where("id = ?", id).Update("article_id", article.ID)

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "文章上传成功",
	})
}

func ArticleUpdate(ctx *gin.Context) {
	id := ctx.PostForm("id")
	content := ctx.PostForm("content")
	DB := database.GetDB()

	DB.Table("articles").Where("software_id = ?", id).Update("content", content)
	// var software model.Software
	// DB.Table("softwares").Where("id = ?", id).First(&software)
	// DB.Table("articles").Where("id = ?", software.ArticleID).Update("content", content)

	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "文章修改成功",
	})
}

// func UploadFile(ctx *gin.Context) {
// 	bucketName := viper.GetString("minio.bucket")
// 	name := ctx.PostForm("name")
// 	// 这里是从form表单里面拿到key为file的文件
// 	file, err := ctx.FormFile("file")
// 	if err != nil {
// 		fmt.Println("upload fail")
// 		return
// 	}
// 	//临时保存
// 	ctx.SaveUploadedFile(file, "./upload/"+name)
// 	MC := database.GetMinioClient()
// 	bucket, err := MC.BucketExists(context.Background(), bucketName)
// 	// mc调用buckets这个函数能不能正常执行
// 	if err != nil {
// 		fmt.Println("mc client fail ")
// 		return
// 	}
// 	// bucket 存在
// 	if !bucket {
// 		fmt.Println("bucket is not exist")
// 	}
// 	uploadInfo, err := MC.FPutObject(context.Background(), bucketName, name, "./upload/"+name, minio.PutObjectOptions{})
// 	if err != nil {
// 		fmt.Println("mc upload fail")
// 		return
// 	}
// 	//数据库
// 	DB := database.GetDB()
// 	newFile := &model.File{
// 		Name:   name,
// 		Bucket: bucketName,
// 	}
// 	// 数据库打上一条文件上传记录
// 	DB.Table("files").Create(&newFile)
// 	ctx.JSON(http.StatusOK, gin.H{
// 		"code": "200",
// 		"data": uploadInfo,
// 		"msg":  "上传成功",
// 	})
// }

func DownloadFile(ctx *gin.Context) {
	bucketName := viper.GetString("minio.bucket")
	name := ctx.Query("name")
	MC := database.GetMinioClient()
	bucket, err := MC.BucketExists(context.Background(), bucketName)
	// mc调用buckets这个函数能不能正常执行
	if err != nil {
		fmt.Println("mc client fail ")
		return
	}
	// bucket 存在
	if !bucket {
		fmt.Println("bucket is not exist")
	}
	file, err := MC.GetObject(context.Background(), bucketName, name, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println("minio get file fail")
		return
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, file)
	// 给浏览器看的头
	ctx.Writer.Header().Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	ctx.Writer.Header().Add("Content-Type", "application/octet-stream")
	ctx.Writer.Header().Add("Content-Transfer-Encoding", "binary")
	ctx.Data(http.StatusOK, "application/octet-stream", buf.Bytes())
}
