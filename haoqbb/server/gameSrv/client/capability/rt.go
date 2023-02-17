package capability

import (
	"Core/Log"
	"Core/haoqbb/server/gameSrv/client/player"
	"Core/haoqbb/server/gameSrv/common/msgHandle"
	"Core/haoqbb/server/gameSrv/common/protocol"
	"time"
)

var index int64
var count int
var allCost float64
var beginTime = map[int64]time.Time{}

func C2SRT(player player.IPlayer) bool {
	index++
	sendMsg := protocol.C2S_Test_RT{
		Index: index,
	}
	beginTime[index] = time.Now()
	player.SendMsgToServer(protocol.SCmd_C2S_RT, &sendMsg)
	//Log.Debug("Send rt index = %v, time = %v", index, time.Now().UnixNano() / int64(time.Millisecond))
	return false
}

func S2CRT(msg *msgHandle.ClientMsg) {
	count++
	data := msg.Data.(*protocol.S2C_Test_RT)
	cost := time.Now().Sub(beginTime[data.Index]).Seconds()
	allCost += cost
	//Log.Debug("Recv rt index = %v, time = %v", data.Index, time.Now().UnixNano() / int64(time.Millisecond))
	if cost > 0.1 {
		Log.WarningLog("rt index = %v, argCost = %v, currCost = %v", data.Index, allCost/float64(count), cost)
	} else {
		Log.Debug("rt index = %v, argCost = %v, currCost = %v", data.Index, allCost/float64(count), cost)
	}
}
