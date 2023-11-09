package Net

import (
	"fmt"
	"github.com/7058011439/haoqbb/Log"
	"net"
	"time"
)

func NewTcpServer(port int, connect ConnectHandle, disconnect ConnectHandle, parse ParseProtocol, msg MsgHandle, options ...Options) INetPool {
	return &TcpService{tcpConnPool: newTcpConnPool(connect, disconnect, parse, msg, options...), listenPort: port}
}

func NewTcpClient(connect ConnectHandle, disconnect ConnectHandle, parse ParseProtocol, msg MsgHandle, options ...Options) INetPool {
	return &TcpClient{tcpConnPool: newTcpConnPool(connect, disconnect, parse, msg, options...)}
}

func WithHeartbeat(handle HeartBeatHandle, interval time.Duration) Options {
	return func(pool *tcpConnPool) {
		if handle != nil && interval > 0 {
			pool.heartbeatHandle = handle
			pool.heartbeatInterval = interval * time.Second
		}
	}
}

func WithPoolId(poolId int) Options {
	return func(pool *tcpConnPool) {
		pool.id = poolId
	}
}

func WithCustomData(compare CompareCustomData) Options {
	return func(pool *tcpConnPool) {
		if compare != nil {
			pool.compareData = compare
		}
	}
}

func WithRecvPackageMaxLimit(size int) Options {
	return func(pool *tcpConnPool) {
		pool.recvPackageLimit = size
	}
}

func WithSendPackageSize(size int) Options {
	return func(pool *tcpConnPool) {
		pool.sendPackageSize = size
	}
}

type TcpClient struct {
	*tcpConnPool
}

type TcpService struct {
	*tcpConnPool
	listenPort int
}

func (s *TcpService) StartServer() bool {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%v", s.listenPort))
	if err != nil {
		Log.ErrorLog("listen error : %v", err)
		return false
	}
	Log.Log("listen server, port: %v", s.listenPort)
	go func() {
		for {
			if c, err := l.Accept(); err != nil {
				Log.ErrorLog("accept error : %v", err)
				break
			} else {
				s.NewConnect(c, nil)
			}
		}
	}()

	return true
}

type INetPool interface {
	StartServer() bool
	GetClientByID(Id uint64) IClient
	Close(Id uint64)
	GetClientByData(data interface{}) IClient
	SendToClient(Id uint64, data []byte)
	GetClientCount() int
	NewConnect(conn net.Conn, data interface{}) IClient
	Range(fun func(client IClient))
	onParseProtocol(data []byte) ([]byte, int)
	onHandleMsg(client IClient, msg []byte)
	onDisconnect(client IClient)
	getRecvPackageLimit() int
	getSendPackageSize() int
}

type IClient interface {
	GetId() uint64
	Close()
	SendMsg([]byte)
	GetIp() string
	GetAddr() string
	CustomData() interface{}
	SetCustomData(interface{})
}
