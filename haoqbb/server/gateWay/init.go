package gateWay

import (
	"github.com/7058011439/haoqbb/haoqbb/node"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
)

func init() {
	gateWay := new(GateWay)
	gateWay.SetName(common.GateWay)

	node.Setup(gateWay)
}
