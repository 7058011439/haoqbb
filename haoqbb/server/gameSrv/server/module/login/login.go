package login

import (
	"encoding/json"
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Util"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/service"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/redis"
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
	if ret.Ret == 0 {
		userId := Util.StrToInt(IRedis.GetRedisSync(service.GetServiceName(), redisOpenUserIdKey, fmt.Sprintf("%v", ret.OpenId)))
		if userId == 0 {
			userId = int(IRedis.IncRedisSyn(service.GetServiceName(), redisGlobalVarKey, redisGlobalVarFieldUserId, 1))
			IRedis.SetRedisSync(service.GetServiceName(), redisOpenUserIdKey, fmt.Sprintf("%v", ret.OpenId), userId)
		}
		player.Login(ret.ClientId, userId)
	}
	sendLoginRet(ret.ClientId, ret.Msg, ret.Ret)
}

// 发送登录结果
func sendLoginRet(clientId uint64, err string, code int) {
	sendMsg := &protocol.S2C_GameLoginResult{
		Success:       code == 0,
		Err:           err,
		Code:          int32(code),
		ServerTimeNow: uint64(time.Now().Second()),
	}
	net.SendMsgToClient(clientId, protocol.SCmd_S2C_Login, sendMsg)
	if code != 0 {
		net.PublicEventByName(common.GateWay, common.GameSrvPlayerOffLine, clientId)
	} else {
		net.PublicEventByName(common.GateWay, common.GameSrvPlayerOnLine, clientId)
	}
}
