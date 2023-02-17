package home

import (
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/protocol"
	iHome "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/server/interface/home"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/multiAccess"
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
