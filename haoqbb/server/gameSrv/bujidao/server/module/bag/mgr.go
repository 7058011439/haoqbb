package bag

import (
	iBag "Core/haoqbb/server/gameSrv/bujidao/server/interface/bag"
	"Core/haoqbb/server/gameSrv/common"
	cBag "Core/haoqbb/server/gameSrv/server/module/bag"
	"Core/haoqbb/server/gameSrv/server/multiAccess"
)

var agent *Mgr

func newObject(id int) multiAccess.IDBData {
	return &bag{
		Bag: cBag.NewBag(id),
	}
}

func Init() {
	agent = &Mgr{
		Mgr: cBag.NewMgr(newObject, common.TabNameBag),
	}
	iBag.SetAgent(agent)
}

type Mgr struct {
	*cBag.Mgr
}

func (m *Mgr) UseItem(userId int, itemId int, itemCount int) bool {
	return true
}
