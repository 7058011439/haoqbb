package node

import (
	"fmt"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"net/http"
	_ "net/http/pprof"
)

func StartServer() {
	port := config.GetNodeID() + 1000
	tcpServer := Net.NewTcpServer(port, newConnectServer, disConnectServer, parseProtocol, msgHandleServer, Net.WithRecvPackageSize(1024*8))
	tcpServer.StartServer()
	// todo
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:1%v", port), nil)
}

func newConnectServer(client Net.IClient) {
	sendMsg := &N2NRegedit{}
	for serviceId, service := range localNodeService {
		sendMsg.ServiceList = append(sendMsg.ServiceList, &ServiceInfo{
			ServiceName: service.GetName(),
			ServiceId:   serviceId,
		})
	}
	client.SendMsg(encodeMsg(sendMsg))
}

func disConnectServer(client Net.IClient) {
	//Log.WarningLog("node disconnect, ip = %v", client.GetIp())
}

func msgHandleServer(client Net.IClient, data []byte) {
	ret := &N2NMsg{}
	ret.Unmarshal(data)
	revMsg(ret.SrcServerId, ret.DestServiceId, ret.MsgType, ret.Data)
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
