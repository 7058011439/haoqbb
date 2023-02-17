package home

import (
	"Core/haoqbb/server/gameSrv/common/msgHandle"
)

func NetUpdateHome(msg *msgHandle.ClientMsg) {
	agent.upgradeLevel(msg.UserId)
}
