package dispatcher

import (
	"Core/haoqbb/node"
	"Core/haoqbb/server/common"
)

func init() {
	dispatcher := new(Dispatcher)
	dispatcher.SetName(common.Dispatcher)

	node.Setup(dispatcher)
}
