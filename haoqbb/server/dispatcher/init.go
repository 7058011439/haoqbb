package dispatcher

import (
	"github.com/7058011439/haoqbb/haoqbb/node"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
)

func init() {
	dispatcher := new(Dispatcher)
	dispatcher.SetName(common.Dispatcher)

	node.Setup(dispatcher)
}
