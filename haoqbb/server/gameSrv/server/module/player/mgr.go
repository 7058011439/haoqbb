package player

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	common2 "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/event"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/multiAccess"
)

var agent *Mgr

func Init() {
	agent = &Mgr{
		DoubleMap:    Stl.NewDoubleMap(),
		ShareDataMgr: multiAccess.NewShareDataMgr(newObject, common2.TabNameUser),
	}
	player.SetAgent(agent)
}

func newObject(id int) multiAccess.IDBData {
	return NewPlayer(id)
}

func NewPlayer(id int) *Player {
	return &Player{
		UserId: id,
	}
}

func NewMgr(fun multiAccess.FunNewObj, collectName string) *Mgr {
	agent = &Mgr{
		DoubleMap:    Stl.NewDoubleMap(),
		ShareDataMgr: multiAccess.NewShareDataMgr(fun, collectName),
	}
	player.SetAgent(agent)
	return agent
}

type Mgr struct {
	*Stl.DoubleMap
	*multiAccess.ShareDataMgr
}

func (m *Mgr) Kick(userId int) {
	if userId == 0 {
		return
	}
	clientId := m.GetClientId(userId)
	m.LogOut(clientId, userId)
}

func (m *Mgr) GetClientId(userId int) uint64 {
	if data := m.DoubleMap.GetKey(userId); data != nil {
		return data.(uint64)
	} else {
		return 0
	}
}

func (m *Mgr) GetUserId(clientId uint64) int {
	if data := m.DoubleMap.GetValue(clientId); data != nil {
		return data.(int)
	} else {
		return 0
	}
}

func (m *Mgr) Login(clientId uint64, userId int) {
	// todo
	m.DoubleMap.Add(clientId, userId)
	m.ShareDataMgr.GetDataAndDo(userId, userId, func(data interface{}, args ...interface{}) {
		clientId := args[0].(uint64)
		Log.Log("Login success userId = %v, clientId = %v, total player = %v", userId, clientId, m.DoubleMap.Len())
		event.PublicGameEvent(event.GameServerLogin, userId)
	}, clientId)
}

func (m *Mgr) LogOut(clientId uint64, userId int) {
	if clientId != 0 {
		m.DoubleMap.RemoveByKey(clientId)
	}
	if userId != 0 {
		m.DoubleMap.RemoveByValue(userId)
	}
	Log.Log("player offline, userId = %v, clientId = %v, total player = %v", userId, clientId, m.DoubleMap.Len())
	net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, clientId)
	m.ShareDataMgr.LogOut(userId)
}
