package node

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"github.com/7058011439/haoqbb/haoqbb/server/protocol"
	"github.com/golang/protobuf/proto"
	"net"
	"sync"
	"time"
)

var nodeConnPool = Net.NewTcpClient(newConnectClient, disConnectClient, parseProtocol, msgHandleClient, Net.WithCustomData(compareData), Net.WithPackageMaxSize(65535*64))
var mutex sync.Mutex
var remoteServiceConn = map[int]Net.IClient{} // 远程服务连接 map[serviceId]connect 可能多个不同serviceId对应同一个connect
var remoteServiceList = map[string][]int{}    // 远程服务列表 map[serviceName][]serviceId

func StartClient() {
	tick := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-tick.C
			connServer()
		}
	}()
}

func compareData(dataA interface{}, dataB interface{}) bool {
	return dataA != nil && dataB != nil && dataA.(int) == dataB.(int)
}

func connServer() {
	currNodeId := config.GetCurrNodeId()
	nodeConfig := config.GetAllNodeConfig()
	for _, node := range nodeConfig {
		if node.NodeId != currNodeId && nodeConnPool.GetClientByData(node.NodeId) == nil && node.ListenAddr != "" {
			if conn, err := net.DialTimeout("tcp", node.ListenAddr, time.Second*5); err == nil {
				nodeConnPool.NewConnect(conn, node.NodeId)
			}
		}
	}
}

func newConnectClient(_ Net.IClient) {
}

// 和其他节点断开连接
func disConnectClient(client Net.IClient) {
	defer mutex.Unlock()
	mutex.Lock()
	for serviceId, c := range remoteServiceConn {
		if c == client {
			for serviceName, serviceIdList := range remoteServiceList {
				for index, id := range serviceIdList {
					if serviceId == id {
						serviceIdList = append(serviceIdList[:index], serviceIdList[index+1:]...)
						Log.Log("disconnect from other service, serviceName = %v, serviceId = %v", serviceName, serviceId)
						for _, service := range localNodeService {
							service.LoseService(serviceName, serviceId)
						}
						break
					}
				}
				if len(serviceIdList) == 0 {
					delete(remoteServiceList, serviceName)
				}
			}
			delete(remoteServiceConn, serviceId)
		}
	}
}

// 连接到其他节点后，其他节点报告服务信息
func msgHandleClient(client Net.IClient, data []byte) {
	defer mutex.Unlock()
	mutex.Lock()
	msg := protocol.N2NRegedit{}
	if err := proto.Unmarshal(data, &msg); err != nil {
		Log.ErrorLog("Failed to parse N2NRegedit, err = %v", err)
		return
	}
	for _, info := range msg.ServiceList {
		remoteServiceConn[int(info.ServiceId)] = client
		remoteServiceList[info.ServiceName] = append(remoteServiceList[info.ServiceName], int(info.ServiceId))
		Log.Log("connect to other service, serviceName = %v, serviceId = %v", info.ServiceName, info.ServiceId)
		for _, service := range localNodeService {
			service.DiscoverService(info.ServiceName, int(info.ServiceId))
		}
	}
}

func sendMsg(srcServiceId, destServiceId int, msgType int, data []byte) {
	msg := protocol.N2NMsg{
		SrcServerId:   int32(srcServiceId),
		DestServiceId: int32(destServiceId),
		MsgType:       int32(msgType),
		Data:          data,
	}
	sendData, _ := proto.Marshal(&msg)
	sendBuff := Stl.NewBuffer(2 + len(sendData))
	sendBuff.Write(Util.Int16ToBytes(int16(len(sendData))))
	sendBuff.Write(sendData)
	if destServiceId == 0 {
		nodeConnPool.Range(func(client Net.IClient) {
			client.SendMsg(sendBuff.Bytes())
		})
	} else {
		if conn := remoteServiceConn[destServiceId]; conn != nil {
			conn.SendMsg(sendBuff.Bytes())
		}
	}
}
