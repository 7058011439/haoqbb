package main

import (
	"github.com/7058011439/haoqbb/haoqbb/node"
	_ "github.com/7058011439/haoqbb/haoqbb/server/dispatcher"
	_ "github.com/7058011439/haoqbb/haoqbb/server/gateWay"
	_ "github.com/7058011439/haoqbb/haoqbb/server/loginSrv"
)

func main() {
	node.Start()
}
