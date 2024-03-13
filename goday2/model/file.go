package model

import "gorm.io/gorm"

type Software struct {
	gorm.Model
	Name        string `json:"name"`        // 软件名
	Description string `json:"description"` // 描述

	Category  string `json:"category"` // 软件种类
	Tag       string `json:"tag"`      //  标签
	ImageID   uint   // 封面id
	FileID    uint   //文件id
	ArticleID uint   // 图文教程id
}

type File struct {
	gorm.Model
	FileName string `json:"filename"` // minio的名字
	Bucket   string `json:"bucket"`   // 存储桶
}

type Image struct {
	gorm.Model
	Image  string `json:"imagename"` // 封面
	Bucket string `json:"bucket"`
}

type Article struct {
	gorm.Model
	SoftwareID uint   // 用于反查
	Content    string `gorm:"type:text"`
}
