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
	"github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
	"github.com/mitchellh/mapstructure"
	"time"
)

type gateConfig struct {
	Port      int
	HeartBeat time.Duration
}

type GateWay struct {
	Net.INetPool
	service.Service
	config            *gateConfig
	addr              string                  // 对外的地址端口
	mapClientSenderId map[int]map[uint64]bool // map[gameServerId]map[clientId]
	mapLoginSrv       map[int]bool            // map[LoginSrvId]bool
}

func (g *GateWay) Init() error {
	if err := mapstructure.Decode(g.ServiceCfg.Other, &g.config); err != nil {
		Log.ErrorLog("Failed to parse gateway Config, err = %v", err)
	}
	g.INetPool = Net.NewTcpServer(g.config.Port, g.connect, g.disConnect, g.parseProtocol, g.NewTcpMsg, Net.WithPoolId(g.GetId()), Net.WithHeartbeat(g.heartBeat, g.config.HeartBeat))
	g.mapClientSenderId = make(map[int]map[uint64]bool, 2)
	g.mapLoginSrv = make(map[int]bool, 2)
	return nil
}

func (g *GateWay) Start() {
	g.StartServer()
	g.addr = fmt.Sprintf("%v:%v", Net.GetOutBoundIP(), g.config.Port)
	g.RegeditDiscoverService(common.LoginSrv, g.discoverLoginSrv)
	g.RegeditLoseService(common.LoginSrv, g.loseLoginSrv)
	g.uploadStatus(0)
	ITimer.SetRepeatTimer(g.GetName(), 1000, g.uploadStatus)
}

func (g *GateWay) InitMsg() {
	g.RegeditHandleTcpMsg(g.handleClientMsg)

	g.RegeditServiceMsg(common.GameSrvToGateClientMsg, g.revMsgFromGameServer)
	g.RegeditServiceMsg(common.GameSrvPlayerOnLine, g.playerOnLine)
	g.RegeditServiceMsg(common.GameSrvPlayerOffLine, g.playerOffLine)
}

func (g *GateWay) sendMsgToGameServer(serverId int, clientId uint64, cmdId int, data []byte) {
	sendMsg := &common.GateWayToGameSrv{
		ClientId: clientId,
		CmdId:    cmdId,
		Data:     data,
	}
	g.PublicEventById(serverId, common.GateToGameSrvClientMsg, sendMsg)
}

func (g *GateWay) sendMsgToLoginSrv(loginSrvId, GameSrvId int, clientId uint64, cmdId int, data []byte) {
	sendMsg := &common.GateWayToLoginSrv{
		GameSrvId: GameSrvId,
		ClientId:  clientId,
		CmdId:     cmdId,
		Data:      data,
	}
	g.PublicEventById(loginSrvId, common.GateToLoginSrvClientMsg, sendMsg)
}

func (g *GateWay) revMsgFromGameServer(gameServerId int, data []byte) {
	revMsg := common.GameSrvToGateWay{}
	if err := json.Unmarshal(data, &revMsg); err == nil {
		if revMsg.ClientId == nil {
			if game, ok := g.mapClientSenderId[gameServerId]; ok && game != nil {
				for clientId := range game {
					g.SendToClient(clientId, common.EncodeSendMsg(int16(gameServerId), 2, int16(revMsg.CmdId), revMsg.Data))
				}
			}
		} else {
			for _, clientId := range revMsg.ClientId {
				g.SendToClient(clientId, common.EncodeSendMsg(0, 2, int16(revMsg.CmdId), revMsg.Data))
			}
		}
	} else {
		Log.ErrorLog("Failed to Unmarshal S2G, data = %v", data)
	}
}

func (g *GateWay) connect(client Net.IClient) {
	Log.Log("new client connect, addr = %v, clientId = %v, have connect = %v", client.GetIp(), client.GetId(), g.GetClientCount())
	g.PublicEventByName("", common.GateWayClientConnect, client.GetId())
}

func (g *GateWay) disConnect(client Net.IClient) {
	Log.Log("client disconnect, addr = %v, clientId = %v, have connect = %v", client.GetIp(), client.GetId(), g.GetClientCount())
	g.PublicEventByName("", common.GateWayClientDisconnect, client.GetId())
}

func (g *GateWay) heartBeat(_ Net.IClient) bool {
	return true
}

func (g *GateWay) parseProtocol(data []byte) (rdata []byte, offset int) {
	if len(data) < 12 {
		return nil, 0
	}
	endPos := int(Util.Uint16(data[3:5])) + 12
	if len(data) >= endPos {
		return data[5 : endPos-1], endPos
	}
	return nil, 0
}

func (g *GateWay) handleClientMsg(clientId uint64, data []byte) {
	cmdId := Util.Uint16(data[0:2])
	serverId := int(Util.Int16(data[4:6]))
	switch cmdId {
	case cmdId_C2S:
		g.sendMsgToGameServer(serverId, clientId, int(Util.Int16(data[2:4])), data[6:])
	case cmdId_C2L:
		for loginSrvId := range g.mapLoginSrv {
			g.sendMsgToLoginSrv(loginSrvId, serverId, clientId, int(Util.Int16(data[2:4])), data[6:])
			break
		}
	}
}

func (g *GateWay) uploadStatus(_ Timer.TimerID, _ ...interface{}) {
	data := &common.GateInfo{
		Addr:         g.addr,
		CpuRate:      System.GetCpuPercent(),
		NetRate:      10,
		ConnectCount: g.GetClientCount(),
	}
	g.PublicEventByName(common.Dispatcher, common.GateToDispatcherStatus, data)
}

func (g *GateWay) playerOnLine(gameServerId int, data []byte) {
	var clientId uint64
	if err := json.Unmarshal(data, &clientId); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on playerOnLine, err = %v, data = %v", err, data)
		return
	}
	if game, ok := g.mapClientSenderId[gameServerId]; ok {
		game[clientId] = true
	} else {
		game = make(map[uint64]bool, 64)
		game[clientId] = true
		g.mapClientSenderId[gameServerId] = game
	}
}

func (g *GateWay) playerOffLine(gameServerId int, data []byte) {
	var clientId uint64
	if err := json.Unmarshal(data, &clientId); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on playerOffLine, err = %v, data = %v", err, data)
		return
	}
	if client := g.GetClientByID(clientId); client != nil {
		client.Close()
	}
	if game, ok := g.mapClientSenderId[gameServerId]; ok {
		delete(game, clientId)
	}
}

func (g *GateWay) loseLoginSrv(serverId int) {
	delete(g.mapLoginSrv, serverId)
	Log.Log("登录服务器断开连接, serverId = %v, 剩余登录服务器 = %v", serverId, len(g.mapLoginSrv))
}

func (g *GateWay) discoverLoginSrv(serverId int) {
	g.mapLoginSrv[serverId] = true
	Log.Log("新登录服务器连接, serverId = %v, 总计登录服务器 = %v", serverId, len(g.mapLoginSrv))
}
