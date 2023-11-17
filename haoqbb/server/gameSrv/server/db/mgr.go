package db

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/event"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/player"
)

var agent *Mgr

func NewMgr(tool IDBTool) *Mgr {
	agent = &Mgr{
		DoubleMap:        Stl.NewDoubleMap(),
		ShareDataMgrSync: NewShareDataMgrSync("user", tool),
	}
	player.SetAgent(agent)
	return agent
}

type IDBData interface {
	IsValid() bool
	Update()
	GetUserId() int64
	Data() interface{}
}

type IDBTool interface {
	NewObj() IDBData
	NewCondition(userId int64) map[string]interface{}
}

type Mgr struct {
	*Stl.DoubleMap
	*ShareDataMgrSync
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
	m.GetData(int64(userId), int64(userId))
	Log.Log("Login success userId = %v, clientId = %v, total player = %v", userId, clientId, m.DoubleMap.Len())
	event.PublicGameEvent(event.GameServerLogin, userId)
}

func (m *Mgr) LogOut(clientId uint64, userId int) {
	if clientId != 0 {
		m.DoubleMap.RemoveByKey(clientId)
	}
	if userId != 0 {
		m.DoubleMap.RemoveByValue(userId)
	}
	Log.Log("player offline, userId = %v, clientId = %v, total player = %v", userId, clientId, m.DoubleMap.Len())
	net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, &common.Uint64{Data: clientId})
	m.ShareDataMgrSync.LogOut(int64(userId))
}
