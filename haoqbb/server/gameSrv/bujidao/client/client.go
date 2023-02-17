package client

import (
	"Core/haoqbb/server/gameSrv/bujidao/client/bag"
	"Core/haoqbb/server/gameSrv/bujidao/client/home"
	"Core/haoqbb/server/gameSrv/bujidao/protocol"
	"Core/haoqbb/server/gameSrv/client"
)

type GameClient struct {
	client.GameClient
}

func (g *GameClient) InitMsg() {
	g.GameClient.InitMsg()
	g.RegeditMsgHandle(protocol.SCmd_S2C_Anything_Add, &protocol.S2C_Anything_Add{}, bag.S2CGiveItem)
	g.RegeditMsgHandle(protocol.SCmd_S2C_HomeUp, &protocol.OperationResult{}, home.S2CHomeUpgrade)
}
