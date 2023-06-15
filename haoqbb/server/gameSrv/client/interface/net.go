package Interface

import (
	"github.com/7058011439/haoqbb/Net"
	"github.com/7058011439/haoqbb/haoqbb/server/common"
	"github.com/golang/protobuf/proto"
	"net"
)

var netConn Net.INetPool

func NewConnManager(connect Net.ConnectHandle, disconnect Net.ConnectHandle, parse Net.ParseProtocol, msg Net.MsgHandle) {
	netConn = Net.NewTcpClient(connect, disconnect, parse, msg)
}

func SendMsgToServer(clientId uint64, cmdId int16, msg proto.Message) {
	data, _ := proto.Marshal(msg)
	if cmdId == 30003 {
		netConn.SendToClient(clientId, common.EncodeSendMsg(10000, 1, cmdId, data))
	} else {
		netConn.SendToClient(clientId, common.EncodeSendMsg(10000, 0, cmdId, data))
	}
}

func NewClient(conn net.Conn) {
	netConn.NewConnect(conn, nil)
}

func GetPlayerCount() int {
	return netConn.GetClientCount()
}
