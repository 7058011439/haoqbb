package node

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"github.com/7058011439/haoqbb/haoqbb/protocol"
	"github.com/golang/protobuf/proto"
	"net/http"
	"strings"
	"sync"
)

const (
	signOK     = -1
	errorTimes = 3
)

var nodeInfo sync.Map // 节点信息 map[Net.IClient]*NodeInfo
var signInfo sync.Map // 签名信息 map[Ip]int

func StartCenterServer() {
	tcpAddr := config.GetCenterAddr()
	if tcpAddr != "" {
		params := strings.Split(tcpAddr, ":")
		if len(params) != 2 {
			Log.ErrorLog("Failed to Start node server, ListenAddr error, try ***.***.***.***:****")
			return
		}
		tcpServer := Net.NewTcpServer(Util.StrToInt(params[1]), newConnectCenterServer, disConnectCenterServer, parseProtocol, msgHandleCenterServer, Net.WithPackageMaxSize(socketCacheSize))
		tcpServer.StartServer()
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:1%v", params[1]), nil)
	}
}

// 新节点连接, 告知新节点现有节点信息
func newConnectCenterServer(client Net.IClient) {
	if info, ok := signInfo.Load(client.GetIp()); ok && info.(int) >= errorTimes {
		Log.ErrorLog("黑名单连接, ip = %v", client.GetIp())
		client.Close()
	} else {
		signInfo.Store(client.GetIp(), 0)
	}
}

func disConnectCenterServer(client Net.IClient) {
	nodeInfo.Delete(client)
	if info, ok := signInfo.Load(client.GetIp()); ok {
		if info.(int) == signOK {
			signInfo.Delete(client.GetIp())
		}
	}
	Log.Log("节点断开, id = %v", client.CustomData())
}

func haveCombine(data1 []string, data2 []string) bool {
	for _, dataA := range data1 {
		for _, dataB := range data2 {
			if dataA == dataB {
				return true
			}
		}
	}
	return false
}

// 告诉新节点, 其他节点信息
func sendNodeInfo(client Net.IClient, needService []string) {
	sendMsg := protocol.NodeList{}
	nodeInfo.Range(func(key, value interface{}) bool {
		info, _ := value.(*protocol.NodeInfo)
		if haveCombine(info.ServiceList, needService) {
			sendMsg.NodeList = append(sendMsg.NodeList, info)
		}
		return true
	})
	client.SendMsg(encodeMsg(&sendMsg))
}

// 新节点上报节点信息
func msgHandleCenterServer(client Net.IClient, data []byte) {
	// 没有签名信息
	if info, _ := signInfo.Load(client.GetIp()); info.(int) != signOK {
		if string(data) == config.GetSign() {
			signInfo.Store(client.GetIp(), signOK)
		} else {
			signInfo.Store(client.GetIp(), info.(int)+1)
			Log.ErrorLog("签名失败, ip = %v", client.GetIp())
			client.Close()
		}
		return
	}

	msg := protocol.NodeInfo{}
	if err := proto.Unmarshal(data, &msg); err != nil {
		Log.ErrorLog("Failed to parse NodeInfo, err = %v", err)
		return
	}

	bExist := false
	nodeInfo.Range(func(key, value interface{}) bool {
		if value.(*protocol.NodeInfo).NodeId == msg.NodeId {
			Log.ErrorLog("发现重复节点, id = %v", msg.NodeId)
			bExist = true
			return false
		}
		return true
	})

	if bExist {
		return
	}
	sendNodeInfo(client, msg.NeedService)
	client.SetCustomData(msg.NodeId)
	Log.Log("发现新节点, id = %v, name = %v, addr = %v", msg.NodeId, msg.NodeName, msg.Addr)

	// 告知现有节点, 有新节点上线
	sendMsg := protocol.NodeList{
		NodeList: []*protocol.NodeInfo{&msg},
	}
	sendData := encodeMsg(&sendMsg)
	nodeInfo.Range(func(key, value interface{}) bool {
		c, _ := key.(Net.IClient)
		if haveCombine(value.(*protocol.NodeInfo).NeedService, msg.GetServiceList()) {
			c.SendMsg(sendData)
		}
		return true
	})

	// 新节点添加到列表
	nodeInfo.Store(client, &msg)
}
