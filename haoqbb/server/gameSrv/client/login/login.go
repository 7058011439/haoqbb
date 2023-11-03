package login

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/interface"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/test"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
)

const (
	getToken = "http://api-chummy.qianchengxing.cn/api/login/game/mobile"
	getMSM   = "http://api-chummy.qianchengxing.cn/api/login/game/send"
)

var phone = 13996434474
var offset = 0

// C2SLogin 匿名登录，只需要currPhone
func C2SLogin(clientId uint64) {
	offset += 1
	currPhone := fmt.Sprintf("%v", phone+offset)
	sendMsg := protocol.C2S_LoginWithToken{
		MachineId: currPhone,
		Token:     "",
		Phone:     currPhone,
		Channel:   0,
	}
	Interface.SendMsgToServer(clientId, protocol.SCmd_C2S_Login, &sendMsg)
}

/*
func C2SLogin(clientId uint64) {
	offset += 1
	currPhone := fmt.Sprintf("%v", phone+offset)
	body := map[string]interface{}{
		"code":     "000000",
		"ditchId":  1,
		"mac":      fmt.Sprintf("%v", Timer.GetOsTimeSecond()+int32(offset)),
		"mobile":   currPhone,
		"serverId": 2,
	}
	Http.PostHttpSync(getMSM, Http.NewHead(nil), Http.NewBody(nil).Add("mobile", currPhone))
	Http.PostHttpAsync(getToken, Http.NewHead(nil), Http.NewBody(body), LoginWithToken, clientId, currPhone)
}

func LoginWithToken(data map[string]interface{}, _ error, callBack ...interface{}) {
	clientId := callBack[0].(uint64)
	currPhone := callBack[1].(string)
	if data["code"].(float64) != 200 {
		Log.ErrorLog("Failed to get token, ret = %v", data)
		return
	} else {
		sendMsg := protocol.C2S_LoginWithToken{
			MachineId: "123456",
			Token:     data["data"].(map[string]interface{})["gameToken"].(string),
			Phone:     currPhone,
			Channel:   1,
		}
		Interface.SendMsgToServer(clientId, protocol.SCmd_C2S_Login, &sendMsg)
	}
}
*/

func S2CLogin(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.S2C_GameLoginResult)
	if data.Success == true {
		p := player.GetPlayerByClientId(msg.ClientId)
		p.UpdateData(player.User{})
		p.SetTestModule(test.GetRandomModule())
		Log.Log("Login success total player = %v", player.Count())
	} else {
		Log.ErrorLog("Login Failed, err = %v", data.Err)
	}
}

func SetStartID(id int) {
	offset = id
}
