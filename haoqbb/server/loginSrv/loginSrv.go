package gateWay

import (
	"encoding/json"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/service"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/http"
)

const (
	httpCheckToken = "http://api-chummy.qianchengxing.cn/api/game/check/token"
)

type LoginSrv struct {
	service.Service
}

func (l *LoginSrv) InitMsg() {
	l.RegeditServiceMsg(common.GateToLoginSrvClientMsg, l.revMsgFromGateWay)

	l.IDispatcher = msgHandle.NewPBDispatcher()
	l.RegeditMsgHandle(protocol.SCmd_C2S_Login, &protocol.C2S_LoginWithToken{}, l.loginWithToken)
}

func (l *LoginSrv) revMsgFromGateWay(_ int, data []byte) {
	msg := &common.GateWayToLoginSrv{}
	if err := json.Unmarshal(data, msg); err != nil {
		Log.ErrorLog("Failed to Unmarshal S2S, data = %v", data)
	} else {
		// 这个地方处理很不科学，将gamesrvid 的值传递给了userid
		l.DispatchMsg(msg.ClientId, msg.GameSrvId, int32(msg.CmdId), msg.Data)
	}
}

func (l *LoginSrv) loginWithToken(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.C2S_LoginWithToken)
	if data.MachineId == "" || data.Token == "" {
		l.noticeLoginRet(msg.UserId, msg.ClientId, 3, "MachineId or Token is nil", 0)
		return
	}
	IHttp.GetAsync(l.GetName(), httpCheckToken+"/"+data.Token, nil, l.httpVerifyTokenCallBack, msg.ClientId, msg.UserId)
}

func (l *LoginSrv) httpVerifyTokenCallBack(getData map[string]interface{}, backData ...interface{}) {
	clientId := backData[0].(uint64)
	gameSrvId := backData[1].(int)
	if getData["code"].(float64) == 200 {
		openId := int(getData["data"].(float64))
		l.noticeLoginRet(gameSrvId, clientId, 0, getData["msg"].(string), openId)
	} else {
		l.noticeLoginRet(gameSrvId, clientId, 1, getData["msg"].(string), 0)
	}
}

func (l *LoginSrv) noticeLoginRet(gameSrvId int, clientId uint64, ret int, msg string, openId int) {
	data := &common.LoginSrvToGameSrv{
		ClientId: clientId,
		Ret:      ret,
		OpenId:   openId,
		Msg:      msg,
	}
	l.PublicEventById(gameSrvId, common.EventLoginSrvLogin, data)
}
