package capability

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"time"
)

var index int64
var count int64
var allCost int64
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
	cost := time.Now().Sub(beginTime[data.Index]).Milliseconds()
	allCost += cost
	// Log.Debug("Recv rt index = %v, time = %v", data.Index, time.Now().UnixNano() / int64(time.Millisecond))
	if cost > 30 {
		Log.WarningLog("rt index = %v, argCost = %vms, currCost = %vms", data.Index, allCost/count, cost)
	} else {
		Log.Debug("rt index = %v, argCost = %vms, currCost = %vms", data.Index, allCost/count, cost)
	}
	delete(beginTime, data.Index)
}
