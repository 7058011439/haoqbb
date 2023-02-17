package home

import (
	"Core/haoqbb/server/gameSrv/bujidao/protocol"
	iHome "Core/haoqbb/server/gameSrv/bujidao/server/interface/home"
	"Core/haoqbb/server/gameSrv/common"
	"Core/haoqbb/server/gameSrv/server/interface/net"
	"Core/haoqbb/server/gameSrv/server/multiAccess"
)

func OnLogin(userId int) {

}

func login(guestId int, hostId int, args ...interface{}) {

}

var agent *Mgr

func newObject(id int) multiAccess.IDBData {
	return &home{
		UserId: id,
	}
}

func Init() {
	agent = &Mgr{
		ShareDataMgr: multiAccess.NewShareDataMgr(newObject, common.TabNameHome),
	}
	iHome.SetAgent(agent)
}

type Mgr struct {
	*multiAccess.ShareDataMgr
}

func (m *Mgr) upgradeLevel(userId int) {
	sendMsg := &protocol.OperationResult{
		IsSuccess: true,
	}
	agent.GetDataAndDo(userId, userId, func(data interface{}, args ...interface{}) {
		h := data.(*home)
		h.upgradeLevel()
		net.SendMsgToUser(userId, protocol.SCmd_S2C_HomeUp, sendMsg)
	})
}
