package login

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Util"
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

func Login(_ int, data []byte) {
	ret := &common.LoginSrvToGameSrv{}
	if err := json.Unmarshal(data, ret); err != nil {
		Log.ErrorLog("处理登录结果错误, err = %v", err)
		return
	}
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
		net.PublicEventByName(common.GateWay, common.SrvPlayerOnLine, clientId)
	} else {
		net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, clientId)
	}
}
