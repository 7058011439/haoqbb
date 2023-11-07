package Net

import (
	"github.com/7058011439/haoqbb/GoroutinePool"
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

const (
	defaultSendPackageMaxSize = 1024  // 默认发送缓存区大小
	defaultPackageMaxSize     = 65535 // 默认最大包长度(超过该包长度丢弃包，防止恶意包攻击)
	revCacheSize              = 1024
)

func getPoolId() int {
	ret := atomic.AddInt32(&poolId, 1)
	return int(ret)
}

func newTcpConnPool(connect ConnectHandle, disconnect ConnectHandle, parse ParseProtocol, msg MsgHandle, opts ...Options) *tcpConnPool {
	if parse == nil {
		parse = defaultParseProtocol
	}

	ret := &tcpConnPool{
		id:               getPoolId(),
		mapClient:        make(map[uint64]IClient, 1024),
		recvPackageLimit: defaultPackageMaxSize,
		sendPackageSize:  defaultSendPackageMaxSize,
		connectHandle:    connect,
		disconnectHandle: disconnect,
		parseProtocol:    parse,
		msgHandle:        msg,
	}
	for _, op := range opts {
		op(ret)
	}

	go ret.heartbeat()
	return ret
}

type tcpConnPool struct {
	id                int                          // 池Id
	seed              int                          // 连接Id自增字段
	mapClient         map[uint64]IClient           // 客户端信息
	mutexClient       sync.RWMutex                 // 客户端锁
	connectHandle     ConnectHandle                // 新连接处理函数
	disconnectHandle  ConnectHandle                // 断开连接处理函数
	parseProtocol     ParseProtocol                // 消息解析函数
	msgHandle         MsgHandle                    // 消息处理函数
	compareData       CompareCustomData            // 自定义数据比较函数
	heartbeatInterval time.Duration                // 心跳间隔(秒)
	heartbeatHandle   HeartBeatHandle              // 心跳处理函数
	recvPackageLimit  int                          // 接收包最大长度(防止乱报攻击)
	sendPackageSize   int                          // 发送包最大长度
	sendTaskPool      *GoroutinePool.GoRoutinePool // 发送任务协程池
}

func (t *tcpConnPool) getClientId() uint64 {
	t.seed = t.seed + 1
	return uint64(t.id)<<32 | uint64(t.seed)
}

func (t *tcpConnPool) StartServer() bool {
	return true
}

func (t *tcpConnPool) GetClientByID(Id uint64) IClient {
	t.mutexClient.RLock()
	defer t.mutexClient.RUnlock()
	return t.mapClient[Id]
}

func (t *tcpConnPool) GetClientByData(data interface{}) IClient {
	t.mutexClient.RLock()
	defer t.mutexClient.RUnlock()
	if t.compareData != nil {
		for _, client := range t.mapClient {
			if t.compareData(client.CustomData(), data) {
				return client
			}
		}
	}
	return nil
}

func (t *tcpConnPool) SendToClient(Id uint64, data []byte) {
	client := t.GetClientByID(Id)
	if client != nil {
		client.SendMsg(data)
	} else {
		//Log.ErrorLog("Failed to SendToClient, client is nil, id = %v", Id)
	}
}

func (t *tcpConnPool) GetClientCount() int {
	t.mutexClient.RLock()
	defer t.mutexClient.RUnlock()
	return len(t.mapClient)
}

func (t *tcpConnPool) Range(fun func(client IClient)) {
	t.mutexClient.RLock()
	defer t.mutexClient.RUnlock()
	for _, client := range t.mapClient {
		fun(client)
	}
}

func (t *tcpConnPool) onParseProtocol(data []byte) ([]byte, int) {
	return t.parseProtocol(data)
}

func (t *tcpConnPool) onHandleMsg(client IClient, msg []byte) {
	if t.msgHandle != nil {
		t.msgHandle(client, msg)
	}
}

func (t *tcpConnPool) onDisconnect(client IClient) {
	t.disconnect(client)
}

func (t *tcpConnPool) getRecvPackageLimit() int {
	return t.recvPackageLimit
}

func (t *tcpConnPool) getSendPackageSize() int {
	return t.sendPackageSize
}

func (t *tcpConnPool) NewConnect(conn net.Conn, data interface{}) IClient {
	client := &Client{
		id:         t.getClientId(),
		conn:       conn,
		customData: data,
		recvBuff:   Stl.NewBuffer(revCacheSize),
		sendBuff:   Stl.NewBuffer(defaultSendPackageMaxSize),
		chClose:    make(chan struct{}, 1),
		INetPool:   t,
	}
	//client.conn.(*net.TCPConn).SetNoDelay(true)
	if client != nil {
		t.mutexClient.Lock()
		t.mapClient[client.GetId()] = client
		t.mutexClient.Unlock()
		if t.connectHandle != nil {
			t.connectHandle(client)
		}
	}
	go client.revMsg()
	return client
}

func (t *tcpConnPool) disconnect(client IClient) {
	t.mutexClient.Lock()
	delete(t.mapClient, client.GetId())
	t.mutexClient.Unlock()
	if t.disconnectHandle != nil {
		t.disconnectHandle(client)
	}
}

func (t *tcpConnPool) heartbeat() {
	if t.heartbeatHandle == nil || t.heartbeatInterval < 1 {
		return
	}
	timerHeartbeat := time.NewTicker(t.heartbeatInterval)
	for {
		select {
		case <-timerHeartbeat.C:
			t.Range(func(client IClient) {
				if !t.heartbeatHandle(client) {
					Log.Log("client heart timeout : %v", client.GetAddr())
					client.Close()
				}
			})
		}
	}
}
