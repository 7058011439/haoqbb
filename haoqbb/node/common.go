package node

import (
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
)

func SendMsgById(srcServiceId int, destServiceId int, msgType int, data []byte) {
	// 本地节点跳过网络，直接处理数据
	revMsg(srcServiceId, destServiceId, msgType, data)
	// 发送数据到远程节点
	sendMsg(srcServiceId, destServiceId, msgType, data)
}

func SendMsgByName(srcServiceId int, serviceName string, msgType int, data []byte) {
	if serviceName == "" {
		SendMsgById(srcServiceId, 0, msgType, data)
	} else {
		// 本地节点跳过网络，直接处理消息
		revMsg(srcServiceId, getLocalServiceId(serviceName), msgType, data)
		// 发送数据到远程节点
		if list, ok := remoteServiceList[serviceName]; ok {
			for serviceId := range list {
				SendMsgById(srcServiceId, serviceId, msgType, data)
			}
		}
	}
}

func HaveServerId(serverId int) bool {
	if _, ok := remoteServiceConn[serverId]; ok {
		return true
	}
	if _, ok := localNodeService[serverId]; ok {
		return true
	}
	return false
}

func parseProtocol(data []byte) (rdata []byte, offset int) {
	allLen := len(data)
	if allLen < 2 {
		return nil, offset
	}
	msgLen := int(Util.Int16(data[0:2]))
	if allLen >= msgLen+2 {
		return data[2 : 2+msgLen], 2 + msgLen
	}
	return nil, 0
}

func encodeMsg(data INodeMsg) []byte {
	return encodeMsgOrigin(data.Marshal())
}

func encodeMsgOrigin(data []byte) []byte {
	sendBuff := Stl.NewBuffer(2 + len(data))
	sendBuff.Write(Util.Int16ToBytes(int16(len(data))))
	sendBuff.Write(data)
	return sendBuff.Bytes()
}
