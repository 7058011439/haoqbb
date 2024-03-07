package bag

import (
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/bag"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/multiAccess"
)

func OnLogin(userId int) {
}

func login(guestId int, hostId int, args ...interface{}) {
	//itemId    := rand.Int31n(10000) + 1
	//itemCount := rand.Int31n(10) + 1
	//GiveItem(guestId, itemId, itemCount)
	//Log.Debug("Bag login, giveItem hostId = %v, itemId = %v, itemCount = %v", hostId, itemId, itemCount)
}

var agent *Mgr

func Init() {
	agent = &Mgr{
		ShareDataMgr: multiAccess.NewShareDataMgr(newObject, common.TabNameBag),
	}
	bag.SetAgent(agent)
}

func newObject(id int) multiAccess.IDBData {
	return NewBag(id)
}

func NewBag(id int) *Bag {
	return &Bag{
		UserId:   id,
		ItemList: make(map[int]int),
	}
}

func NewMgr(fun multiAccess.FunNewObj, collectName string) *Mgr {
	agent = &Mgr{
		ShareDataMgr: multiAccess.NewShareDataMgr(fun, collectName),
	}
	bag.SetAgent(agent)
	return agent
}

type Mgr struct {
	*multiAccess.ShareDataMgr
}

func (m *Mgr) GiveItem(userId int, itemId int, itemCount int) {
	m.ShareDataMgr.GetDataAndDo(userId, userId, func(data interface{}, args ...interface{}) {
		bag := data.(*Bag)
		bag.GiveItem(itemId, itemCount)
	}, itemId, itemCount)
}

func (m *Mgr) TakeItem(userId int, itemId int, itemCount int) {
	m.ShareDataMgr.GetDataAndDo(userId, userId, func(data interface{}, args ...interface{}) {
		b := data.(*Bag)
		b.TakeItem(itemId, itemCount)
	}, itemId, itemCount)
}

func (m *Mgr) CheckItem(userId int, itemId int, itemCount int) bool {
	bag := m.ShareDataMgr.GetData(userId, userId).(*Bag)
	if bag != nil {
		return bag.CheckItem(itemId, itemCount)
	} else {
		return false
	}
}
