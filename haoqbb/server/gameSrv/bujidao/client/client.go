package client

import (
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/client/bag"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/client/home"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/client"
)

type GameClient struct {
	client.GameClient
}

func (g *GameClient) InitMsg() {
	g.GameClient.InitMsg()
	g.RegeditMsgHandle(protocol.SCmd_S2C_Anything_Add, &protocol.S2C_Anything_Add{}, bag.S2CGiveItem)
	g.RegeditMsgHandle(protocol.SCmd_S2C_HomeUp, &protocol.OperationResult{}, home.S2CHomeUpgrade)
}
