package task

import (
	"day1/database"
	"day1/model"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

// 初始化任务
func InitTask() {
	MQTTClient := database.GetMQTTClient()
	topic := viper.GetString("mqtt.topic")
	if token := MQTTClient.Subscribe(topic, 0, listen); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
}

// 监听
func listen(client mqtt.Client, msg mqtt.Message) {
	topic := viper.GetString("mqtt.topic")
	if topic == msg.Topic() {
		data := model.MqttMsg{}
		json.Unmarshal(msg.Payload(), &data)
		cmd_router(data)
	}
}

// 命令分发
func cmd_router(data model.MqttMsg) {
	if data.Cmd == "register" {
		device_register(data)
	}
	if data.Cmd == "report" {
		device_report(data)
	}
}

func device_register(data model.MqttMsg) {
	DB := database.GetDB()
	device := model.Device{}
	DB.Table("devices").Where("mac = ?", data.Mac).First(&device)
	if device.ID == 0 {
		newDevice := model.Device{
			Mac: data.Mac,
		}
		DB.Create(&newDevice)
	}
}

func device_report(data model.MqttMsg) {
	DB := database.GetDB()
	device := model.Device{}
	DB.Table("devices").Where("mac = ?", data.Mac).First(&device)
	if device.ID == 0 {
		newDevice := model.Device{
			Mac: data.Mac,
		}
		DB.Create(&newDevice)
	}
	DB.Table("devices").Where("mac = ?", data.Mac).First(&device)
	str := fmt.Sprint(data.Data)
	// fmt.Println(str, "---------")
	// fmt.Println(data.Data)
	newDeviceData := model.DeviceData{
		MacID: uint64(device.ID),
		Sorts: data.Param,
		Data:  str,
	}
	DB.Create(&newDeviceData)
}
