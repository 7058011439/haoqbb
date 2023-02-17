package event

import (
	"github.com/7058011439/haoqbb/EventBus"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/server/module/home"
	cBag "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/module/bag"
)

const (
	ServerGameServer = "GameSrv"
	GameServerLogin  = "login"
)

func init() {
	EventBus.Subscribe(ServerGameServer+GameServerLogin, home.OnLogin)
	EventBus.Subscribe(ServerGameServer+GameServerLogin, cBag.OnLogin)
}

func PublicGameEvent(topic string, args ...interface{}) {
	EventBus.Publish(ServerGameServer+topic, args...)
}
