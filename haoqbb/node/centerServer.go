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

var nodeInfoByConn sync.Map // map[Net.IClient]*NodeInfo

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
	sendMsg := protocol.NodeList{}
	nodeInfoByConn.Range(func(key, value interface{}) bool {
		if info, ok := value.(*protocol.NodeInfo); ok {
			sendMsg.NodeList = append(sendMsg.NodeList, info)
			return true
		}
		return false
	})
	client.SendMsg(encodeMsg(&sendMsg))
}

func disConnectCenterServer(client Net.IClient) {
	nodeInfoByConn.Delete(client)
	Log.Log("节点断开, id = %v", client.CustomData())
}

// 新节点上报节点信息
func msgHandleCenterServer(client Net.IClient, data []byte) {
	msg := protocol.NodeInfo{}
	if err := proto.Unmarshal(data, &msg); err != nil {
		Log.ErrorLog("Failed to parse NodeInfo, err = %v", err)
		return
	}

	bExist := false
	nodeInfoByConn.Range(func(key, value interface{}) bool {
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
	client.SetCustomData(msg.NodeId)
	Log.Log("发现新节点, id = %v, name = %v, addr = %v", msg.NodeId, msg.NodeName, msg.Addr)

	// 告知现有节点, 有新节点上线
	sendMsg := protocol.NodeList{
		NodeList: []*protocol.NodeInfo{&msg},
	}
	sendData := encodeMsg(&sendMsg)
	nodeInfoByConn.Range(func(key, value interface{}) bool {
		if c, ok := key.(Net.IClient); ok {
			c.SendMsg(sendData)
			return true
		}
		return false
	})

	// 新节点添加到列表
	nodeInfoByConn.Store(client, &msg)
}
