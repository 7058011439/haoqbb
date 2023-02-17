package home

import (
	"Core/haoqbb/server/gameSrv/bujidao/protocol"
	"Core/haoqbb/server/gameSrv/client/player"
	"Core/haoqbb/server/gameSrv/common/msgHandle"
)

func C2SHomeUpgrade(player player.IPlayer) bool {
	sendMsg := protocol.C2S_HomeUp{}
	player.SendMsgToServer(protocol.SCmd_C2S_HomeUp, &sendMsg)
	return true
}

func S2CHomeUpgrade(msg *msgHandle.ClientMsg) {
	//Log.Debug("S2CHomeUpgrade = %v", msg.data)
}
