package server

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	cProtocol "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/capability"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	iPlayer "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
	iService "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/service"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/module/bag"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/module/login"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/module/player"
	"github.com/7058011439/haoqbb/haoqbb/service"
)

type GameSrv struct {
	service.Service
	mapClientSenderId map[uint64]int // map[clientId]GateWayId
}

func (g *GameSrv) Init() error {
	g.mapClientSenderId = make(map[uint64]int, 2048) // map[clientId]GateWayId
	net.SetNetAgent(g)
	iService.SetServiceAgent(g)
	return nil
}

func (g *GameSrv) Start() {
	g.RegeditLoseService(common.GateWay, g.lostGateWay)
	bag.Init()
	player.Init()
}

func (g *GameSrv) InitMsg() {
	g.RegeditServiceMsg(common.GwForwardClToSrv, g.revMsgFromGateWay)
	g.RegeditServiceMsg(common.GwClConnect, g.clientConnect)
	g.RegeditServiceMsg(common.GwClDisconnect, g.clientDisconnect)
	g.RegeditServiceMsg(common.EventLoginSrvLogin, login.Login)

	g.IDispatcher = msgHandle.NewPBDispatcher()
	g.RegeditMsgHandle(cProtocol.SCmd_C2S_RT, &cProtocol.C2S_Test_RT{}, capability.NetC2SRT)
}

func (g *GameSrv) revMsgFromGateWay(_ int, data []byte) {
	msg := &common.GwForwardClToSrvTag{}
	if err := json.Unmarshal(data, msg); err != nil {
		Log.ErrorLog("Failed to Unmarshal S2S, data = %v", data)
	} else {
		g.DispatchMsg(msg.ClientId, iPlayer.GetUserId(msg.ClientId), int32(msg.CmdId), msg.Data)
		if userId := iPlayer.GetUserId(msg.ClientId); userId != 0 {
			g.DispatchMsg(msg.ClientId, userId, int32(msg.CmdId), msg.Data)
		} else {
			Log.WarningLog("没有找到对应的userId, clientId = %v", msg.ClientId)
			net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, msg.ClientId)
		}
	}
}

func (g *GameSrv) SendMsgToClient(clientId uint64, cmdId int32, data []byte) {
	g.BroadCastMsgToClient([]uint64{clientId}, cmdId, data)
}

func (g *GameSrv) SendMsgToUser(userId int, cmdId int32, data []byte) {
	g.SendMsgToClient(iPlayer.GetClientId(userId), cmdId, data)
}

func (g *GameSrv) BroadCastMsgToClient(clientIds []uint64, cmdId int32, data []byte) {
	if clientIds == nil {
		sendMsg := &common.GwForwardSrvToClTag{
			ClientId: clientIds,
			CmdId:    int(cmdId),
			Data:     data,
		}
		g.PublicEventById(0, common.GwForwardSrvToCl, sendMsg)
	}
	serverIds := map[int][]uint64{}
	for _, clientId := range clientIds {
		if serverId, ok := g.mapClientSenderId[clientId]; ok {
			serverIds[serverId] = append(serverIds[serverId], clientId)
		} else {
			Log.ErrorLog("Failed to SendMsgToClient, not find serverId, clientId = %v", clientId)
		}
	}
	for serverId, clientIds := range serverIds {
		sendMsg := &common.GwForwardSrvToClTag{
			ClientId: clientIds,
			CmdId:    int(cmdId),
			Data:     data,
		}
		g.PublicEventById(serverId, common.GwForwardSrvToCl, sendMsg)
	}
}

func (g *GameSrv) BroadCastMsgToUser(userIds []int, cmdId int32, data []byte) {
	if userIds == nil {
		g.BroadCastMsgToClient(nil, cmdId, data)
	}
	var clientIds []uint64
	for _, userId := range userIds {
		clientIds = append(clientIds, iPlayer.GetClientId(userId))
	}
	g.BroadCastMsgToClient(clientIds, cmdId, data)
}

func (g *GameSrv) clientConnect(srcServiceId int, data []byte) {
	var clientId uint64
	if err := json.Unmarshal(data, &clientId); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on clientConnect, err = %v, data = %v", err, data)
		return
	}
	g.mapClientSenderId[clientId] = srcServiceId
}

func (g *GameSrv) clientDisconnect(_ int, data []byte) {
	var clientId uint64
	if err := json.Unmarshal(data, &clientId); err != nil {
		Log.ErrorLog("Failed to json.Unmarshal on clientDisconnect, err = %v, data = %v", err, data)
		return
	}
	delete(g.mapClientSenderId, clientId)
	userId := iPlayer.GetUserId(clientId)
	iPlayer.Kick(userId)
}

func (g *GameSrv) lostGateWay(serviceId int) {
	for clientId, id := range g.mapClientSenderId {
		if id == serviceId {
			delete(g.mapClientSenderId, clientId)
			userId := iPlayer.GetUserId(clientId)
			iPlayer.Kick(userId)
		}
	}
}
