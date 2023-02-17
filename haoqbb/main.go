package main

import (
	"Core/haoqbb/node"
	_ "Core/haoqbb/server/dispatcher"
	_ "Core/haoqbb/server/gateWay"
	_ "Core/haoqbb/server/loginSrv"
)

func main() {
	node.Start()
}
