package node

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/config"
	"net/http"
	"strings"
	"sync"
)

const (
	signOK     = -1
	errorTimes = 3
)

var nodeInfo sync.Map // 节点信息 map[Net.IClient]*NodeInfo map<节点连接>节点信息
var signInfo sync.Map // 签名信息 map[string]int map<ip>签名失败次数

func StartCenterServer() {
	tcpAddr := config.GetCenterAddr()
	if tcpAddr != "" {
		params := strings.Split(tcpAddr, ":")
		if len(params) != 2 {
			Log.Error("Failed to Start node server, ListenAddr error, try ***.***.***.***:****")
			return
		}
		tcpServer := Net.NewTcpServer(Util.StrToInt(params[1]), newConnectCenterServer, disConnectCenterServer, parseProtocol, msgHandleCenterServer)
		tcpServer.StartServer()
		go http.ListenAndServe(fmt.Sprintf("0.0.0.0:1%v", params[1]), nil)
	}
}

// 新节点连接, 判定是否黑名单ip，如果是黑名单ip，关闭连接
func newConnectCenterServer(client Net.IClient) {
	ip := client.GetIp()
	if info, ok := signInfo.Load(ip); ok {
		// 该节点多次连接中心节点, 而且都认证失败, 认为该节点有问题, 需要关闭连接
		if info.(int) >= errorTimes {
			Log.Error("黑名单连接, ip = %v", ip)
			client.Close()
		}
	} else {
		// 该节点没有被连接过，设置值
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
	sendMsg := &NodeList{}
	nodeInfo.Range(func(key, value interface{}) bool {
		info, _ := value.(*NodeInfo)
		if haveCombine(info.ServiceList, needService) {
			sendMsg.NodeList = append(sendMsg.NodeList, info)
		}
		return true
	})
	client.SendMsg(encodeMsg(sendMsg))
}

// 新节点上报节点信息
func msgHandleCenterServer(client Net.IClient, data []byte) {
	// 没有签名信息
	if info, _ := signInfo.Load(client.GetIp()); info.(int) != signOK {
		if string(data) == config.GetSign() {
			signInfo.Store(client.GetIp(), signOK)
		} else {
			signInfo.Store(client.GetIp(), info.(int)+1)
			Log.Error("签名失败, ip = %v", client.GetIp())
			client.Close()
		}
		return
	}

	nInfo := &NodeInfo{}
	nInfo.Unmarshal(data)

	bExist := false
	nodeInfo.Range(func(key, value interface{}) bool {
		if value.(*NodeInfo).NodeId == nInfo.NodeId {
			Log.Error("发现重复节点, id = %v", nInfo.NodeId)
			bExist = true
			return false
		}
		return true
	})

	if bExist {
		return
	}
	sendNodeInfo(client, nInfo.NeedService)
	client.SetCustomData(nInfo.NodeId)
	Log.Log("发现新节点, id = %v, name = %v, addr = %v", nInfo.NodeId, nInfo.NodeName, nInfo.Addr)

	// 告知现有节点, 有新节点上线
	sendData := encodeMsg(&NodeList{
		NodeList: []*NodeInfo{nInfo},
	})
	nodeInfo.Range(func(key, value interface{}) bool {
		c, _ := key.(Net.IClient)
		if haveCombine(value.(*NodeInfo).NeedService, nInfo.ServiceList) {
			c.SendMsg(sendData)
		}
		return true
	})

	// 新节点添加到列表
	nodeInfo.Store(client, nInfo)
}
