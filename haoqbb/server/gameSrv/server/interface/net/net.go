package net

import (
	"github.com/golang/protobuf/proto"
)

type INetGameSrv interface {
	SendMsgToUser(userId int, cmdId int32, data []byte)
	SendMsgToClient(clientId uint64, cmdId int32, data []byte)
	PublicEventByName(serviceName string, eventType int, data interface{})
	BroadCastMsgToClient(clientIds []uint64, cmdId int32, data []byte)
	BroadCastMsgToUser(userIds []int, cmdId int32, data []byte)
}

type Net struct {
	INetGameSrv
}

var net = Net{}

func SetNetAgent(s INetGameSrv) {
	net.INetGameSrv = s
}

func SendMsgToClient(clientId uint64, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.SendMsgToClient(clientId, cmdId, data)
}

func SendMsgToUser(userId int, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.SendMsgToUser(userId, cmdId, data)
}

func BroadCastMsgToClient(clientIds []uint64, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.BroadCastMsgToClient(clientIds, cmdId, data)
}

func BroadCastMsgToUser(userIds []int, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.BroadCastMsgToUser(userIds, cmdId, data)
}

func PublicEventByName(serviceName string, eventType int, data interface{}) {
	net.PublicEventByName(serviceName, eventType, data)
}
