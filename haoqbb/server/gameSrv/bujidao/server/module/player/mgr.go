package player

import (
	iPlayer "Core/haoqbb/server/gameSrv/bujidao/server/interface/player"
	"Core/haoqbb/server/gameSrv/common"
	cPlayer "Core/haoqbb/server/gameSrv/server/module/player"
	"Core/haoqbb/server/gameSrv/server/multiAccess"
)

var agent *Mgr

func newObject(id int) multiAccess.IDBData {
	return &player{
		Player: cPlayer.NewPlayer(id),
	}
}

func Init() {
	agent = &Mgr{
		Mgr: cPlayer.NewMgr(newObject, common.TabNameUser),
	}
	iPlayer.SetAgent(agent)
}

type Mgr struct {
	*cPlayer.Mgr
}
