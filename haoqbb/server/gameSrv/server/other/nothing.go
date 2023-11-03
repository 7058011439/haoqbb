package other

import (
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
)

func NetNothingWithBack(msg *msgHandle.ClientMsg) {
	net.SendMsgToClient(msg.ClientId, protocol.SCmd_S2C_Nothing_WithReply, &protocol.C2S_Test_RT{})
}

func NetNothingWithOutBack(msg *msgHandle.ClientMsg) {

}
