package model

import "gorm.io/gorm"

// 设备表
type Device struct {
	gorm.Model
	Mac string
}

// 设备数据记录表
type DeviceData struct {
	gorm.Model
	MacID uint64
	Sorts string
	Data  string
}
