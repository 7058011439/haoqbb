package server

import (
	"Core/haoqbb/server/gameSrv/bujidao/protocol"
	"Core/haoqbb/server/gameSrv/bujidao/server/module/bag"
	"Core/haoqbb/server/gameSrv/bujidao/server/module/home"
	"Core/haoqbb/server/gameSrv/bujidao/server/module/player"
	"Core/haoqbb/server/gameSrv/server"
)

type BuJiDaoSrv struct {
	server.GameSrv
}

func (g *BuJiDaoSrv) Start() {
	g.GameSrv.Start()
	bag.Init()
	home.Init()
	player.Init()
}

func (g *BuJiDaoSrv) InitMsg() {
	g.GameSrv.InitMsg()

	g.RegeditMsgHandle(protocol.SCmd_C2S_Anything_Add, &protocol.C2S_Anything_Add{}, bag.NetGiveAnything)
	g.RegeditMsgHandle(protocol.SCmd_C2S_HomeUp, &protocol.C2S_HomeUp{}, home.NetUpdateHome)
}
