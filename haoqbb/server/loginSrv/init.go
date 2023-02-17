package gateWay

import (
	"Core/haoqbb/node"
	"Core/haoqbb/server/common"
)

func init() {
	loginSrv := new(LoginSrv)
	loginSrv.SetName(common.LoginSrv)

	node.Setup(loginSrv)
}
