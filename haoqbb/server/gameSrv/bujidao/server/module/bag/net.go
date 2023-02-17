package bag

import (
	"Core/haoqbb/server/gameSrv/bujidao/protocol"
	"Core/haoqbb/server/gameSrv/common/msgHandle"
	iCBag "Core/haoqbb/server/gameSrv/server/interface/bag"
)

func NetGiveAnything(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.C2S_Anything_Add)
	for _, item := range data.Data {
		iCBag.GiveItem(msg.UserId, int(item.Id), int(item.Count))
	}
}
