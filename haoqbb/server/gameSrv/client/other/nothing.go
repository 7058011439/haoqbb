package other

import (
	"github.com/7058011439/haoqbb/String"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"math/rand"
)

func C2SNothingWithReply(player player.IPlayer) bool {
	sendData := &protocol.C2S_Test_Nothing_WithReply{
		Index: rand.Int63(),
		Msg:   String.RandStr(32),
	}
	for i := 0; i < 50; i++ {
		player.SendMsgToServer(protocol.SCmd_C2S_Nothing_WithReply, sendData)
	}
	return true
}

func C2SNothingWithOutReply(player player.IPlayer) bool {
	sendData := &protocol.C2S_Test_Nothing_WithOutReply{
		Index: rand.Int63(),
		Msg:   String.RandStr(32),
	}
	for i := 0; i < 50; i++ {
		player.SendMsgToServer(protocol.SCmd_C2S_Nothing_WithOutReply, sendData)
	}
	return true
}

func S2CNothingWithReply(msg *msgHandle.ClientMsg) {

}
