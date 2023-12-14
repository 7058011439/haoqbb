package db

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
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

type SDBTool struct{}

func (s *SDBTool) NewCondition(userId int64) map[string]interface{} {
	return map[string]interface{}{
		"user.userid": userId,
	}
}

type Mgr struct {
	*Stl.DoubleMap
	*ShareDataMgrSync
}

func (m *Mgr) Kick(userId int64) {
	if userId == 0 {
		return
	}
	clientId := m.GetClientId(userId)
	m.LogOut(clientId, userId)
}

func (m *Mgr) GetClientId(userId int64) uint64 {
	if data := m.DoubleMap.GetKey(userId); data != nil {
		return data.(uint64)
	} else {
		return 0
	}
}

func (m *Mgr) GetUserId(clientId uint64) int64 {
	if data := m.DoubleMap.GetValue(clientId); data != nil {
		return data.(int64)
	} else {
		return 0
	}
}

func (m *Mgr) Login(clientId uint64, userId int64) {
	m.DoubleMap.Add(clientId, userId)
	m.GetData(userId, userId)
	Log.Log("Login success userId = %v, clientId = %v, total player = %v", userId, clientId, m.DoubleMap.Len())
}

func (m *Mgr) LogOut(clientId uint64, userId int64) {
	if clientId != 0 {
		m.DoubleMap.RemoveByKey(clientId)
	}
	if userId != 0 {
		m.DoubleMap.RemoveByValue(userId)
	}
	Log.Log("player offline, userId = %v, clientId = %v, total player = %v", userId, clientId, m.DoubleMap.Len())
	net.PublicEventByName(common.GateWay, common.SrvPlayerOffLine, &common.Uint64{Data: clientId})
	m.ShareDataMgrSync.LogOut(userId)
}
