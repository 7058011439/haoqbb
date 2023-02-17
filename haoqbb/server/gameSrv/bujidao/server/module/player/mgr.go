package player

import (
	iPlayer "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/server/interface/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common"
	cPlayer "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/module/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/multiAccess"
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
