package node

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"github.com/7058011439/haoqbb/haoqbb/protocol"
	"github.com/golang/protobuf/proto"
)

func StartServer() {
	port := config.GetNodeID() + 1000
	tcpServer := Net.NewTcpServer(port, newConnectServer, disConnectServer, parseProtocol, msgHandleServer)
	tcpServer.StartServer()
	// todo
	// go http.ListenAndServe(fmt.Sprintf("0.0.0.0:1%v", port), nil)
}

func newConnectServer(client Net.IClient) {
	//Log.Log("new node connect, ip = %v", client.GetIp())
	sendMsg := protocol.N2NRegedit{}
	for serviceId, service := range localNodeService {
		sendMsg.ServiceList = append(sendMsg.ServiceList, &protocol.ServiceInfo{
			ServiceName: service.GetName(),
			ServiceId:   int32(serviceId),
		})
	}
	client.SendMsg(encodeMsg(&sendMsg))
}

func disConnectServer(client Net.IClient) {
	//Log.WarningLog("node disconnect, ip = %v", client.GetIp())
}

func msgHandleServer(client Net.IClient, data []byte) {
	msg := protocol.N2NMsg{}
	if err := proto.Unmarshal(data, &msg); err != nil {
		Log.ErrorLog("Failed to parse N2NMsg, err = %v", err)
		return
	}
	revMsg(int(msg.SrcServerId), int(msg.DestServiceId), int(msg.MsgType), msg.Data)
}

func revMsg(srcServiceId int, destServiceId int, msgType int, data []byte) {
	if destServiceId == 0 {
		for _, service := range localNodeService {
			service.NewServiceMsg(srcServiceId, msgType, data)
		}
	} else {
		if service, _ := localNodeService[destServiceId]; service != nil {
			service.NewServiceMsg(srcServiceId, msgType, data)
		}
	}
}
