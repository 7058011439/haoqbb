package home

import (
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/msgHandle"
)

func NetUpdateHome(msg *msgHandle.ClientMsg) {
	agent.upgradeLevel(msg.UserId)
}
