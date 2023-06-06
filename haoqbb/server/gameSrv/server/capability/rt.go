package capability

import (
	"github.com/7058011439/haoqbb/haoqbb/msgHandle"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/common/protocol"
	"github.com/7058011439/haoqbb/haoqbb/server/gameSrv/server/interface/net"
)

func NetC2SRT(msg *msgHandle.ClientMsg) {
	data := msg.Data.(*protocol.C2S_Test_RT)
	sendMsg := &protocol.C2S_Test_RT{
		Index: data.Index,
	}
	net.SendMsgToClient(msg.ClientId, protocol.SCmd_S2C_RT, sendMsg)
	//Log.Debug("Handle rt index = %v, time = %v", data.Index, time.Now().UnixNano() / int64(time.Millisecond))
}
