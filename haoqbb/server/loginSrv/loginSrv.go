package gateWay

import (
	"fmt"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/service"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/http"
)

const (
	chanGuest = iota // 匿名登录
	chanSelf         // 自营渠道
)

type LoginSrv struct {
	service.Service
	loginFun map[int]func(data *protocol.C2S_LoginWithToken, clientId uint64)
}

func (l *LoginSrv) Init() error {
	l.loginFun = map[int]func(data *protocol.C2S_LoginWithToken, clientId uint64){
		chanGuest: l.loginGuest,
		chanSelf:  l.loginMy,
	}
	return nil
}

func (l *LoginSrv) InitMsg() {
	l.RegeditServiceMsg(common.EventGameSrvLogin, l.login)
}

func (l *LoginSrv) login(_ int, data []byte) {
	ret := &common.GameSrvToLoginSrv{}
	ret.Unmarshal(data)
	if fun, ok := l.loginFun[int(ret.Data.Channel)]; ok {
		fun(ret.Data, ret.ClientId)
	} else {
		l.noticeLoginRet(ret.Data.Channel, ret.Data.SrvId, ret.ClientId, "未知渠道", "")
	}
}

func (l *LoginSrv) loginGuest(data *protocol.C2S_LoginWithToken, clientId uint64) {
	l.noticeLoginRet(data.Channel, data.SrvId, clientId, "", data.MachineId)
}

func (l *LoginSrv) loginMy(data *protocol.C2S_LoginWithToken, clientId uint64) {
	IHttp.GetAsync(l.GetName(), "http://api-chummy.qianchengxing.cn/api/game/check/token/"+data.Token, nil, func(getData map[string]interface{}, backData ...interface{}) {
		clientId := backData[0].(uint64)
		data := backData[1].(*protocol.C2S_LoginWithToken)
		openId := ""
		if getData["code"].(float64) == 200 {
			openId = fmt.Sprintf("%v", getData["data"])
		}
		l.noticeLoginRet(data.Channel, data.SrvId, clientId, getData["msg"].(string), openId)
	}, clientId, data)
}

func (l *LoginSrv) noticeLoginRet(channel int32, gameSrvId int32, clientId uint64, msg string, openId string) {
	if openId != "" {
		openId = fmt.Sprintf("%v_%v", channel, openId)
	}
	data := &common.LoginSrvToGameSrv{
		ClientId: clientId,
		OpenId:   openId,
		Msg:      msg,
	}
	l.SendMsgToServiceById(int(gameSrvId), common.EventLoginSrvLogin, data)
}
