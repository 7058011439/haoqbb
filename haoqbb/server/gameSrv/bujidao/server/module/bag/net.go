package bag

import (
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/bujidao/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/msgHandle"
	iCBag "github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/bag"
)

func NetGiveAnything(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.C2S_Anything_Add)
	for _, item := range data.Data {
		iCBag.GiveItem(msg.UserId, int(item.Id), int(item.Count))
	}
}
