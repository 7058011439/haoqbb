package node

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"github.com/7058011439/haoqbb/haoqbb/protocol"
	"github.com/golang/protobuf/proto"
	"net"
	"sync"
	"time"
)

var nodeConnPool = Net.NewTcpClient(func(client Net.IClient) {}, disConnectClient, parseProtocol, msgHandleClient, Net.WithCustomData(compareData), Net.WithPackageMaxSize(65535*64))

var nodeConnPoolCenterClient = Net.NewTcpClient(connectCenterNode, func(client Net.IClient) {}, parseProtocol, msgHandleCenterClient)

var mutex sync.Mutex
var remoteServiceConn = map[int]Net.IClient{}     // 远程服务连接 map[serviceId]connect 可能多个不同serviceId对应同一个connect
var remoteServiceList = map[string]map[int]bool{} // 远程服务列表 map[serviceName][serviceId]bool

func compareData(dataA interface{}, dataB interface{}) bool {
	return dataA != nil && dataB != nil && dataA.(int) == dataB.(int)
}

func StartClient() {
	if conn, err := net.DialTimeout("tcp", config.GetCenterAddr(), time.Second*5); err == nil {
		nodeConnPoolCenterClient.NewConnect(conn, nil)
	} else {
		Log.ErrorLog("连接到中心节点错误, err = %v", err)
	}
}

func connectCenterNode(client Net.IClient) {
	Log.Log("成功连接到中心节点")
	client.SendMsg(encodeMsgOrigin([]byte(config.GetSign())))

	// 上报自己节点信息
	nodeConfig := config.GetNodeConfig()
	client.SendMsg(encodeMsg(&protocol.NodeInfo{
		NodeId:      int32(nodeConfig.NodeId),
		NodeName:    nodeConfig.NodeName,
		Addr:        nodeConfig.ListenAddr,
		ServiceList: nodeConfig.ServiceList,
		NeedService: nodeConfig.NeedService,
	}))
}

// 和其他节点断开连接
func disConnectClient(client Net.IClient) {
	defer mutex.Unlock()
	mutex.Lock()
	for serviceId, c := range remoteServiceConn {
		if c == client {
			for serviceName, serviceIdList := range remoteServiceList {
				if _, ok := serviceIdList[serviceId]; ok {
					Log.Log("disconnect from other service, serviceName = %v, serviceId = %v", serviceName, serviceId)
					for _, service := range localNodeService {
						service.LoseService(serviceName, serviceId)
					}
					delete(serviceIdList, serviceId)
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
		if remoteServiceList[info.ServiceName] == nil {
			remoteServiceList[info.ServiceName] = map[int]bool{}
		}
		remoteServiceList[info.ServiceName][int(info.ServiceId)] = true
		Log.Log("connect to other node, nodeId = %v, serviceName = %v, serviceId = %v", client.CustomData(), info.ServiceName, info.ServiceId)
		for _, service := range localNodeService {
			service.DiscoverService(info.ServiceName, int(info.ServiceId))
		}
	}
}

// 连接到中心节点后，中心节点报其他告节点信息
func msgHandleCenterClient(client Net.IClient, data []byte) {
	defer mutex.Unlock()
	mutex.Lock()
	msg := protocol.NodeList{}
	if err := proto.Unmarshal(data, &msg); err != nil {
		Log.ErrorLog("Failed to parse NodeList, err = %v", err)
		return
	}
	for _, info := range msg.NodeList {
		Log.Log("发现新节点, id = %v, name = %v, addr = %v", info.NodeId, info.NodeName, info.Addr)
		if conn, err := net.DialTimeout("tcp", info.Addr, time.Second*5); err == nil {
			nodeConnPool.NewConnect(conn, info.NodeId)
		} else {
			Log.ErrorLog("连接到新节点失败, err = %v", err)
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
	sendData := encodeMsg(&msg)
	if destServiceId == 0 {
		nodeConnPool.Range(func(client Net.IClient) {
			client.SendMsg(sendData)
		})
	} else {
		if conn := remoteServiceConn[destServiceId]; conn != nil {
			conn.SendMsg(sendData)
		}
	}
}
