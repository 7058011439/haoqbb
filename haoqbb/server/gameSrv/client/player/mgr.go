package player

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client/interface"
	"github.com/7058011439/haoqbb/haoqbb/service/interface/timer"
	"sync"
)

var mutex sync.RWMutex
var allPlayer = make(map[uint64]IPlayer, 1024)

func Count() int {
	mutex.RLock()
	defer mutex.RUnlock()
	count := 0
	for _, p := range allPlayer {
		if p.IsLogin() {
			count++
		}
	}
	return count
}

func Range(fun func(player IPlayer)) {
	mutex.RLock()
	defer mutex.RUnlock()
	for _, player := range allPlayer {
		fun(player)
	}
}

func NewPlayer(clientId uint64) IPlayer {
	mutex.Lock()
	defer mutex.Unlock()
	player := &Player{clientId: clientId}
	allPlayer[clientId] = player
	return player
}

func RemovePlayer(clientId uint64) {
	mutex.Lock()
	defer mutex.Unlock()
	if player := allPlayer[clientId]; player != nil {
		ITimer.SetOffTimer(Interface.GetServiceName(), player.TimerId())
		Log.Log("player offline, total player = %v", len(allPlayer)-1)
	}
	delete(allPlayer, clientId)
}

func GetPlayerByClientId(clientId uint64) IPlayer {
	mutex.RLock()
	defer mutex.RUnlock()
	return allPlayer[clientId]
}
