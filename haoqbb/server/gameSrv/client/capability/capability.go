package capability

import (
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/player"
)

func Main(_ Timer.TimerID, args ...interface{}) {
	clientId := args[0].(uint64)
	player := player.GetPlayerByClientId(clientId)
	if player != nil {
		C2SRT(player)
	}
}
