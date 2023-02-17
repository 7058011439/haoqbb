package gateWay

import (
	"github.com/7058011439/haoqbb/haoqbb/node"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
)

func init() {
	loginSrv := new(LoginSrv)
	loginSrv.SetName(common.LoginSrv)

	node.Setup(loginSrv)
}
