package bag

import (
	"Core/haoqbb/server/gameSrv/bujidao/protocol"
	"Core/haoqbb/server/gameSrv/client/player"
	"Core/haoqbb/server/gameSrv/common/msgHandle"
	"math/rand"
)

func C2SGiveItem(player player.IPlayer) bool {
	sendMsg := protocol.C2S_Anything_Add{}
	for i := 0; i < 10; i++ {
		sendMsg.Data = append(sendMsg.Data, &protocol.XAnything{
			Id:    rand.Int31n(10) + 1,
			Count: rand.Int31n(10) + 1,
		})
	}
	player.SendMsgToServer(protocol.SCmd_C2S_Anything_Add, &sendMsg)
	//Log.Debug("C2SGiveItem = %v", sendMsg.data)
	return true
}

func C2STakeItem(player player.IPlayer) bool {
	sendMsg := protocol.C2S_Anything_Sub{}
	for i := 0; i < 10; i++ {
		sendMsg.Data = append(sendMsg.Data, &protocol.XAnything{
			Id:    rand.Int31n(10) + 1,
			Count: rand.Int31n(10) + 1,
		})
	}
	player.SendMsgToServer(protocol.SCmd_C2S_Anything_Sub, &sendMsg)
	return true
}

func S2CGiveItem(msg *msgHandle.ClientMsg) {

}
