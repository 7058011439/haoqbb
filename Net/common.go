package Net

var poolId int32

type Options func(config *tcpConnPool)
type ConnectHandle func(client IClient)
type HeartBeatHandle func(client IClient) bool
type ParseProtocol func(data []byte) (rdata []byte, offset int)
type MsgHandle func(client IClient, data []byte)
type CompareCustomData func(dataA interface{}, dataB interface{}) bool

func defaultParseProtocol(data []byte) ([]byte, int) {
	return data, len(data)
}
