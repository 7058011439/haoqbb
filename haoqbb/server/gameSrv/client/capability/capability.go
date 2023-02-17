package capability

import (
	"Core/Timer"
	"Core/haoqbb/server/gameSrv/client/player"
)

func Main(_ Timer.TimerID, args ...interface{}) {
	clientId := args[0].(uint64)
	player := player.GetPlayerByClientId(clientId)
	if player != nil {
		C2SRT(player)
	}
}
