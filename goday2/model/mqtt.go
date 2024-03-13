package model

type MqttMsg struct {
	Mac   string      `json:"mac"`
	Cmd   string      `json:"cmd"`
	Param string      `json:"param"`
	Data  interface{} `json:"data"`
}
