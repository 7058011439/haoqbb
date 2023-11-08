package login

import (
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/service"
	IRedis "github.com/7058011439/haoqbb/haoqbb/service/interface/redis"
	"time"
)

const (
	redisOpenUserIdKey        = "OpenUserId"
	redisGlobalVarKey         = "GlobalVar"
	redisGlobalVarFieldUserId = "UserId"
)

func WithToken(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.C2S_LoginWithToken)
	if userId := player.GetUserId(msg.ClientId); userId != 0 {
		sendLoginRet(msg.ClientId, "重复登录", false)
		return
	}
	if loginSrvId := service.GetLoginSrvId(); loginSrvId != 0 {
		sendData := &common.GameSrvToLoginSrv{
			ClientId: msg.ClientId,
			Data:     data,
		}
		net.SendMsgToServiceById(loginSrvId, common.EventGameSrvLogin, sendData)
	} else {
		sendLoginRet(msg.ClientId, "服务器准备中...", false)
		return
	}
}

func Ret(_ int, data []byte) {
	ret := &common.LoginSrvToGameSrv{}
	ret.Unmarshal(data)
	userId := 0
	if ret.OpenId != "" {
		userId = Util.StrToInt(IRedis.GetRedisSync(service.GetServiceName(), redisOpenUserIdKey, ret.OpenId))
		if userId == 0 {
			userId = int(IRedis.IncRedisSyn(service.GetServiceName(), redisGlobalVarKey, redisGlobalVarFieldUserId, 1))
			IRedis.SetRedisSync(service.GetServiceName(), redisOpenUserIdKey, ret.OpenId, userId)
		}

		player.Login(ret.ClientId, userId)
	}
	sendLoginRet(ret.ClientId, ret.Msg, userId != 0)
}

// 发送登录结果
func sendLoginRet(clientId uint64, err string, ret bool) {
	sendMsg := &protocol.S2C_GameLoginResult{
		Success:       ret,
		Err:           err,
		ServerTimeNow: uint64(time.Now().Second()),
	}
	net.SendMsgToClient(clientId, protocol.SCmd_S2C_Login, sendMsg)
	if ret {
		net.PublicEventByName(common.GateWay, common.SrvPlayerOnLine, &common.Uint64{Data: clientId})
	} else {
		net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, &common.Uint64{Data: clientId})
	}
}
