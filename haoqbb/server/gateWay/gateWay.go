package gateWay

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/System"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/service"
	ITimer "github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
	"github.com/mitchellh/mapstructure"
	"time"
)

type gateConfig struct {
	Port      int
	HeartBeat time.Duration
}

type GateWay struct {
	service.Service
	Config     *gateConfig             // 网关配置
	ClientList map[int]map[uint64]bool // map[gameServerId]map[clientId]
	addr       string                  // 对外的地址端口
}

func (g *GateWay) Init() error {
	if err := mapstructure.Decode(g.ServiceCfg.Other, &g.Config); err != nil {
		Log.ErrorLog("Failed to parse gateway Config, err = %v", err)
	}
	g.ClientList = make(map[int]map[uint64]bool, 2)
	return nil
}

func (g *GateWay) Start() {
	g.addr = fmt.Sprintf("%v:%v", Net.GetOutBoundIP(), g.Config.Port)
	g.uploadStatus(0)
	ITimer.SetRepeatTimer(g.GetName(), 1000, g.uploadStatus)
}

func (g *GateWay) InitMsg() {
	g.RegeditServiceMsg(common.GwForwardSrvToCl, g.RecvMsgFromSrv)

	g.RegeditServiceMsg(common.SrvPlayerOnLine, g.PlayerOnLine)
	g.RegeditServiceMsg(common.SrvPlayerOffLine, g.PlayerOffLine)
}

func (g *GateWay) RecvMsgFromSrv(serverId int, data []byte) {
	// 这个地方有点绕, 如果其他服有指定发送给具体的客户端，那就发送给指定客户端，如果没指定，那就是区服广播
	revMsg := common.GwForwardSrvToClTag{}
	if err := json.Unmarshal(data, &revMsg); err == nil {
		if revMsg.ClientId == nil {
			if game, ok := g.ClientList[serverId]; ok && game != nil {
				for clientId := range game {
					g.SendToClient(clientId, common.EncodeSendMsg(int16(serverId), 0, int16(revMsg.CmdId), revMsg.Data))
				}
			}
		} else {
			for _, clientId := range revMsg.ClientId {
				g.SendToClient(clientId, common.EncodeSendMsg(0, 0, int16(revMsg.CmdId), revMsg.Data))
			}
		}
	} else {
		Log.ErrorLog("Failed to Unmarshal GwForwardSrvToClTag, data = %v", string(data))
	}
}

func (g *GateWay) OnConnect(client Net.IClient) {
	Log.Log("new client connect, addr = %v, clientId = %v, have connect = %v", client.GetIp(), client.GetId(), g.GetClientCount())
	g.PublicEventByName("", common.GwClConnect, client.GetId())
}

func (g *GateWay) OnDisConnect(client Net.IClient) {
	Log.Log("client disconnect, addr = %v, clientId = %v, have connect = %v", client.GetIp(), client.GetId(), g.GetClientCount())
	g.PublicEventByName("", common.GwClDisconnect, client.GetId())
}

// ParseProtocol 解析数据流, 请配合HandleClientMsg 使用
/* 这是第一个奇葩的协议, 分为 协议头 + 数据 + 协议尾
协议头11个字节，分别为:
0: 固定 0xFE
1-2: 预留
3-4: 数据长度
5-6: 主命令号
7-8: 子命令号
9-10: 预留(目前用作服务id)

数据:
根据协议头的数据长度推算

协议尾1个字节:
0: 固定 0xEE
这个地方就要了 协议头部分(主命令号开始) + 数据 */
func (g *GateWay) ParseProtocol(data []byte) (rdata []byte, offset int) {
	if len(data) < 12 {
		return nil, 0
	}
	endPos := int(Util.Uint16(data[3:5])) + 12
	if len(data) >= endPos {
		return data[5 : endPos-1], endPos
	}
	return nil, 0
}

// HandleClientMsg 处理客户端消息, 请配合ParseProtocol 使用
/* 这是第一个奇葩的数据, 分为 (部分)协议头 + 数据
部分协议头6个字节，分别为:
0-1: 主命令(暂时无用)
2-4: 子命令
5-6: 服务id

数据:
*/
func (g *GateWay) HandleClientMsg(clientId uint64, data []byte) {
	if len(data) < 6 {
		Log.ErrorLog("failed to HandleClientMsg, data too shoot, data = %v", data)
		return
	}
	g.ForwardClMsgToSrv(int(Util.Int16(data[4:6])), clientId, int(Util.Int16(data[2:4])), data[6:])
}

func (g *GateWay) ForwardClMsgToSrv(serverId int, clientId uint64, cmdId int, data []byte) {
	g.PublicEventById(serverId, common.GwForwardClToSrv, &common.GwForwardClToSrvTag{
		ClientId: clientId,
		CmdId:    cmdId,
		Data:     data,
	})
}

func (g *GateWay) uploadStatus(_ Timer.TimerID, _ ...interface{}) {
	// 原则上该框架不应该拉起其他协程，但是该操作因为要读取硬件信息，极为耗时，会阻塞主协程，所以特别go了一下
	go func() {
		data := &common.GsInfoTag{
			Addr:         g.addr,
			MemRate:      System.GetMemPercent(),
			CpuRate:      System.GetCpuPercent(),
			NetRate:      System.GetNetRate(),
			ConnectCount: g.GetClientCount(),
		}
		g.PublicEventByName(common.Dispatcher, common.GwToDsStatus, data)
	}()
}

func (g *GateWay) PlayerOnLine(gameServerId int, data []byte) {
	var clientId uint64
	if err := json.Unmarshal(data, &clientId); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on PlayerOnLine, err = %v, data = %v", err, data)
		return
	}
	if _, ok := g.ClientList[gameServerId]; !ok {
		g.ClientList[gameServerId] = make(map[uint64]bool, 64)
	}
	g.ClientList[gameServerId][clientId] = true
}

func (g *GateWay) PlayerOffLine(gameServerId int, data []byte) {
	var clientId uint64
	if err := json.Unmarshal(data, &clientId); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on PlayerOffLine, err = %v, data = %v", err, data)
		return
	}
	if client := g.GetClientByID(clientId); client != nil {
		client.Close()
	}
	if game, ok := g.ClientList[gameServerId]; ok {
		delete(game, clientId)
	}
}
