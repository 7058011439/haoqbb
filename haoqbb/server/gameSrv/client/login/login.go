package login

import (
	"Core/Http"
	"Core/Log"
	"Core/Timer"
	"Core/haoqbb/server/gameSrv/client/interface"
	"Core/haoqbb/server/gameSrv/client/player"
	"Core/haoqbb/server/gameSrv/common/msgHandle"
	"Core/haoqbb/server/gameSrv/common/protocol"
	"fmt"
)

const (
	getToken = "http://api-chummy.qianchengxing.cn/api/login/game/mobile"
	getMSM   = "http://api-chummy.qianchengxing.cn/api/login/game/send"
)

var phone = 13996434474
var offset = 0

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
	Http.PostHttpSyncNew(getMSM, Http.NewHead(nil), Http.NewBody(nil).Add("mobile", currPhone))
	Http.PostHttpAsyncNew(getToken, Http.NewHead(nil), Http.NewBody(body), LoginWithToken, clientId, currPhone)
}

func LoginWithToken(data map[string]interface{}, callBack ...interface{}) {
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
		}
		Interface.SendMsgToServer(clientId, protocol.SCmd_C2S_Login, &sendMsg)
	}
}

func S2CLogin(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.S2C_GameLoginResult)
	if data.Success == true {
		p := player.GetPlayerByClientId(msg.ClientId)
		p.UpdateData(player.User{})
		Log.Log("Login success total player = %v", player.Count())
	} else {
		Log.ErrorLog("Login Failed, err = %v", data.Err)
	}
}

func SetStartID(id int) {
	offset = id
}
