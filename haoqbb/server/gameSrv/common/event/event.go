package event

import (
	"Core/EventBus"
	"Core/haoqbb/server/gameSrv/bujidao/server/module/home"
	cBag "Core/haoqbb/server/gameSrv/server/module/bag"
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
