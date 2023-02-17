package player

import (
	"Core/Log"
	"Core/Stl"
	"Core/haoqbb/server/common"
	common2 "Core/haoqbb/server/gameSrv/common"
	"Core/haoqbb/server/gameSrv/common/event"
	"Core/haoqbb/server/gameSrv/server/interface/net"
	"Core/haoqbb/server/gameSrv/server/interface/player"
	"Core/haoqbb/server/gameSrv/server/multiAccess"
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
	m.ShareDataMgr.GetDataAndDo(userId, userId, func(data interface{}, args ...interface{}) {
		clientId := args[0].(uint64)
		m.DoubleMap.Add(clientId, userId)
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
	net.PublicEventByName(common.GateWay, common.GameSrvPlayerOffLine, clientId)
	m.ShareDataMgr.LogOut(userId)
}
