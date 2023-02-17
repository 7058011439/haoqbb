package gateWay

import (
	"Core/haoqbb/node"
	"Core/haoqbb/server/common"
)

func init() {
	gateWay := new(GateWay)
	gateWay.SetName(common.GateWay)

	node.Setup(gateWay)
}
