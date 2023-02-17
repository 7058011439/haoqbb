package node

import (
	"Core/Log"
	"Core/Net"
	"Core/Stl"
	"Core/Util"
	"Core/haoqbb/config"
	"Core/haoqbb/server/protocol"
	"fmt"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
)

func StartServer() {
	tcpAddr := config.GetCurrNodeListenAddr()
	if tcpAddr != "" {
		params := strings.Split(tcpAddr, ":")
		if len(params) != 2 {
			Log.ErrorLog("Failed to Start node server, ListenAddr error, try ***.***.***.***:****")
			return
		}
		tcpServer := Net.NewTcpServer(Util.StrToInt(params[1]), newConnectServer, disConnectServer, parseProtocol, msgHandleServer, Net.WithPackageMaxSize(socketCacheSize))
		tcpServer.StartServer()
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:1%v", params[1]), nil)
	}
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
	sendData, _ := proto.Marshal(&sendMsg)
	sendBuff := Stl.NewBuffer(2 + len(sendData))
	sendBuff.Write(Util.Int16ToBytes(int16(len(sendData))))
	sendBuff.Write(sendData)
	client.SendMsg(sendBuff.Bytes())
}

func disConnectServer(client Net.IClient) {
	//Log.WarningLog("node disconnect, ip = %v", client.GetIp())
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
