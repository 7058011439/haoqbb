package login

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
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
	if ret.UserId != 0 {
		player.Login(ret.ClientId, ret.UserId)
	}
	sendLoginRet(ret.ClientId, ret.Msg, ret.UserId != 0)
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
