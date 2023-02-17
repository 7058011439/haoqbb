package Interface

type INetWorkGateWay interface {
	SendMsgToClient(clientId uint64, data []byte)
	PublicEvent(eventType int16, data interface{})
}

type NetWork struct {
	INetWorkGateWay
}

var net = NetWork{}

func PublicEvent(eventType int16, data interface{}) {
	net.PublicEvent(eventType, data)
}
