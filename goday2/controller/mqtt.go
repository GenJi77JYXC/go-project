package controller

import (
	"day1/database"
	"fmt"
	"net/http"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gin-gonic/gin"
)

func PublicMsg(ctx *gin.Context) {
	topic := ctx.PostForm("topic")
	data := ctx.PostForm("data")
	MQTTClient := database.GetMQTTClient()
	if token := MQTTClient.Publish(topic, 0, false, data); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "public success",
	})
}

func subscribe(c mqtt.Client, m mqtt.Message) {
	fmt.Printf("TOPIC:%s\n", m.Topic())
	fmt.Printf("MSG:%s\n", m.Payload())
}

func Listen(ctx *gin.Context) {
	topic := ctx.PostForm("topic")
	MQTTClient := database.GetMQTTClient()
	if token := MQTTClient.Subscribe(topic, 0, subscribe); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code": "200",
		"data": "",
		"msg":  "subscribe success",
	})

}
