package other

import (
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
)

func C2SNothingWithReply(player player.IPlayer) bool {
	for i := 0; i < 20; i++ {
		player.SendMsgToServer(protocol.SCmd_C2S_Nothing_WithReply, &protocol.C2S_Test_RT{})
	}
	return true
}

func C2SNothingWithOutReply(player player.IPlayer) bool {
	// todo
	for i := 0; i < 100; i++ {
		player.SendMsgToServer(protocol.SCmd_C2S_Nothing_WithOutReply, &protocol.C2S_Test_RT{})
	}
	return true
}

func S2CNothingWithReply(msg *msgHandle.ClientMsg) {

}
