package server

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	cProtocol "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/capability"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	iPlayer "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
	iService "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/service"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/module/login"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/other"
	"github.com/7058011439/haoqbb/haoqbb/service"
)

type GameSrv struct {
	service.Service
	mapClientSenderId map[uint64]int // map[clientId]GateWayId
	loginSrv          map[int]bool   // map[loginSrvId]bool
}

func (g *GameSrv) Init() error {
	g.mapClientSenderId = make(map[uint64]int, 2048) // map[clientId]GateWayId
	g.loginSrv = make(map[int]bool, 2)               // map[loginSrvId]bool
	net.SetNetAgent(g)
	iService.SetServiceAgent(g)
	return nil
}

func (g *GameSrv) InitMsg() {
	g.RegeditLoseService(common.GateWay, g.lostGateWay)
	g.RegeditDiscoverService(common.LoginSrv, g.discoverLoginSrv)
	g.RegeditLoseService(common.LoginSrv, g.loseLoginSrv)

	g.RegeditServiceMsg(common.GwForwardClToSrv, g.revMsgFromGateWay)
	g.RegeditServiceMsg(common.GwClConnect, g.clientConnect)
	g.RegeditServiceMsg(common.GwClDisconnect, g.clientDisconnect)
	g.RegeditServiceMsg(common.EventLoginSrvLogin, login.Ret)

	g.IDispatcher = msgHandle.NewPBDispatcher()
	g.RegeditMsgHandle(cProtocol.SCmd_C2S_Login, &cProtocol.C2S_LoginWithToken{}, login.WithToken)
	g.RegeditMsgHandle(cProtocol.SCmd_C2S_RT, &cProtocol.C2S_Test_RT{}, capability.NetC2SRT)
	g.RegeditMsgHandle(cProtocol.SCmd_C2S_Nothing_WithReply, &cProtocol.C2S_Test_Nothing_WithReply{}, other.NetNothingWithBack)
	g.RegeditMsgHandle(cProtocol.SCmd_C2S_Nothing_WithOutReply, &cProtocol.C2S_Test_Nothing_WithOutReply{}, other.NetNothingWithOutBack)
}

func (g *GameSrv) Start() {

}

func (g *GameSrv) revMsgFromGateWay(_ int, data []byte) {
	msg := &common.GwForwardClToSrvTag{}
	msg.Unmarshal(data)
	if userId := iPlayer.GetUserId(msg.ClientId); userId != 0 || msg.CmdId == cProtocol.SCmd_C2S_Login {
		g.DispatchMsg(msg.ClientId, userId, int32(msg.CmdId), msg.Data)
	} else {
		Log.WarningLog("没有找到对应的userId, clientId = %v", msg.ClientId)
		net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, &common.Uint64{Data: msg.ClientId})
	}
}

func (g *GameSrv) SendMsgToClient(clientId uint64, cmdId int32, data []byte) {
	g.BroadCastMsgToClient([]uint64{clientId}, cmdId, data)
}

func (g *GameSrv) SendMsgToUser(userId int64, cmdId int32, data []byte) {
	g.SendMsgToClient(iPlayer.GetClientId(userId), cmdId, data)
}

func (g *GameSrv) BroadCastMsgToClient(clientIds []uint64, cmdId int32, data []byte) {
	if clientIds == nil {
		sendMsg := &common.GwForwardSrvToClTag{
			ClientId: clientIds,
			CmdId:    int(cmdId),
			Data:     data,
		}
		g.SendMsgToServiceById(0, common.GwForwardSrvToCl, sendMsg)
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
		g.SendMsgToServiceById(serverId, common.GwForwardSrvToCl, sendMsg)
	}
}

func (g *GameSrv) BroadCastMsgToUser(userIds []int64, cmdId int32, data []byte) {
	if userIds == nil {
		g.BroadCastMsgToClient(nil, cmdId, data)
	}
	var clientIds []uint64
	for _, userId := range userIds {
		clientIds = append(clientIds, iPlayer.GetClientId(userId))
	}
	g.BroadCastMsgToClient(clientIds, cmdId, data)
}

func (g *GameSrv) GetLoginSrvId() int {
	for id := range g.loginSrv {
		return id
	}
	return 0
}

func (g *GameSrv) clientConnect(srcServiceId int, data []byte) {
	clientId := &common.Uint64{}
	clientId.Unmarshal(data)
	g.mapClientSenderId[clientId.Data] = srcServiceId
}

func (g *GameSrv) clientDisconnect(_ int, data []byte) {
	clientId := &common.Uint64{}
	clientId.Unmarshal(data)
	delete(g.mapClientSenderId, clientId.Data)
	userId := iPlayer.GetUserId(clientId.Data)
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

func (g *GameSrv) loseLoginSrv(serverId int) {
	delete(g.loginSrv, serverId)
	Log.Log("登录服务器断开连接, serverId = %v, 剩余登录服务器 = %v", serverId, len(g.loginSrv))
}

func (g *GameSrv) discoverLoginSrv(serverId int) {
	g.loginSrv[serverId] = true
	Log.Log("新登录服务器连接, serverId = %v, 总计登录服务器 = %v", serverId, len(g.loginSrv))
}
