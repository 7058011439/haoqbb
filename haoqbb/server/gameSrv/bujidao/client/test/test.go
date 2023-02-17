package test

import (
	"Core/haoqbb/server/gameSrv/bujidao/client/bag"
	"Core/haoqbb/server/gameSrv/bujidao/client/home"
	"Core/haoqbb/server/gameSrv/client/player"
	"Core/haoqbb/server/gameSrv/client/test"
)

const (
	testIdHomeUp = iota
	testIdGiveItem
	testIdTakeItem
	testDoNothing
)

func init() {
	test.InsertTestModule(testIdHomeUp, true, home.C2SHomeUpgrade, 0, map[int]int{testIdHomeUp: 1})
	test.InsertTestModule(testIdGiveItem, true, bag.C2SGiveItem, 0, map[int]int{testIdGiveItem: 1, testIdTakeItem: 2})
	test.InsertTestModule(testIdTakeItem, true, bag.C2STakeItem, 0, map[int]int{testIdGiveItem: 1})
	test.InsertTestModule(testDoNothing, true, doNothing, 10, map[int]int{testDoNothing: 1})
	test.OnInitOver()
}

func doNothing(_ player.IPlayer) bool {
	return true
}
