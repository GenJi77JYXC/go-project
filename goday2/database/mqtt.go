package database

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/spf13/viper"
)

var MQTTClient mqtt.Client

func InitMQTT(name string) mqtt.Client {
	//拉取配置
	host := viper.GetString("mqtt.host")
	port := viper.GetString("mqtt.port")
	username := viper.GetString("mqtt.username")
	password := viper.GetString("mqtt.password")
	// 设置服务器地址以及clientID
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + host + ":" + port).SetClientID(name)
	// 设用用户名盒密码
	opts.SetUsername(username)
	opts.SetPassword(password)
	// 心跳时间
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	//设置默认发布后调用的函数
	opts.SetDefaultPublishHandler(f)
	//进行连接
	c := mqtt.NewClient(opts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	MQTTClient = c
	return c
}

func GetMQTTClient() mqtt.Client {
	return MQTTClient
}

// 输出主题盒消息
var f mqtt.MessageHandler = func(c mqtt.Client, m mqtt.Message) {
	fmt.Printf("TOPIC:%s\n", m.Topic())
	fmt.Printf("MSG:%s\n", m.Payload())
}
