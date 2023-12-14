package net

import (
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/golang/protobuf/proto"
)

type INetGameSrv interface {
	SendMsgToUser(userId int64, cmdId int32, data []byte)
	SendMsgToClient(clientId uint64, cmdId int32, data []byte)
	SendMsgToServiceByName(serviceName string, eventType int, msg common.ServiceMsg)
	SendMsgToServiceById(serviceId int, msgType int, msg common.ServiceMsg)
	BroadCastMsgToClient(clientIds []uint64, cmdId int32, data []byte)
	BroadCastMsgToUser(userIds []int64, cmdId int32, data []byte)
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

func SendMsgToUser(userId int64, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.SendMsgToUser(userId, cmdId, data)
}

func BroadCastMsgToClient(clientIds []uint64, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.BroadCastMsgToClient(clientIds, cmdId, data)
}

func BroadCastMsgToUser(userIds []int64, cmdId int32, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	net.BroadCastMsgToUser(userIds, cmdId, data)
}

func PublicEventByName(serviceName string, eventType int, msg common.ServiceMsg) {
	net.SendMsgToServiceByName(serviceName, eventType, msg)
}

func SendMsgToServiceById(serviceId int, msgType int, msg common.ServiceMsg) {
	net.SendMsgToServiceById(serviceId, msgType, msg)
}
