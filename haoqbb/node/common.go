package node

const (
	socketCacheSize = 1024 * 1024 * 64
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
			for _, serviceId := range list {
				SendMsgById(srcServiceId, serviceId, msgType, data)
			}
		}
	}
}
