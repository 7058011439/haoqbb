package Net

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/String"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/Util"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

/*
连接数量	单包大小	qps(W/S)	网络(MB/S)
1		50		1750		830
10		50		2800		1350
100		50		2700		1300
*/

func encodeMsg(data []byte) []byte {
	buff := Stl.NewBuffer(len(data) + 4)
	buff.WriteInt(len(data))
	buff.Write(data)
	return buff.Bytes()
}

func parseProtocol(data []byte) (rdata []byte, offset int) {
	if len(data) >= 4 {
		dataLen := Util.Int(data[0:4])
		if len(data) >= dataLen+4 {
			return data[4 : dataLen+4], dataLen + 4
		}
	}
	return nil, 0
}

var lastRecvCount uint64
var totalRecvCount uint64
var totalSendCount uint64
var connectCount uint64
var cost = Timer.NewTiming(Timer.Second)

func msgHandle(client IClient, data []byte) {
	atomic.AddUint64(&totalRecvCount, 1)
}

func printData(_ Timer.TimerID, _ ...interface{}) {
	recvCount := atomic.LoadUint64(&totalRecvCount)
	sendCount := atomic.LoadUint64(&totalSendCount)
	Log.Debug("%v, currCount = %3.f/s, totalRecvCount = %v, totalSendCount = %v, dataLen = %.3f m/s, tcpConn = %v", cost, float64(recvCount-lastRecvCount)/cost.GetCost(), recvCount, sendCount, float64((recvCount-lastRecvCount)*50)/1048576/cost.GetCost(), connectCount)
	lastRecvCount = recvCount
	cost.ReStart()
}

func TestClient_SendMsg(t *testing.T) {
	tcpServer := NewTcpServer(6666, nil, nil, parseProtocol, msgHandle, WithRecvPackageSize(1024*16))
	tcpServer.StartServer()
	Timer.AddRepeatTimer(1000, printData)

	tcpClient := NewTcpClient(func(client IClient) {
		connectCount++
	}, nil, parseProtocol, msgHandle, WithSendPackageSize(1024*1024))
	data := encodeMsg([]byte(String.RandStr(50)))
	for i := 0; i < 10; i++ {
		if conn, err := net.DialTimeout("tcp", "127.0.0.1:6666", time.Second*5); err == nil {
			client := tcpClient.NewConnect(conn, nil)
			go func() {
				for {
					client.SendMsg(data)
					//atomic.AddUint64(&totalSendCount, 1)
				}
			}()
			time.Sleep(time.Millisecond)
		} else {
			Log.Error("连接到新节点失败, err = %v", err)
			break
		}
	}
	select {}
}

var tcpServer INetPool
var tcpClient INetPool

func timerClose(id Timer.TimerID, args ...interface{}) {
	clientId := args[0].(uint64)
	if client := tcpServer.GetClientByID(clientId); client != nil {
		client.Close()
	}
}

func TestClient_Close(t *testing.T) {
	tcpServer = NewTcpServer(6666, func(client IClient) {
		Log.Debug("新连接 = %v, 总计连接 = %v", client.GetId(), tcpServer.GetClientCount())
		Timer.AddOnceTimer(500, timerClose, client.GetId())
	}, func(client IClient) {
		Log.Debug("断开连接 = %v, 总计连接 = %v", client.GetId(), tcpServer.GetClientCount())
	}, parseProtocol, msgHandle)
	tcpServer.StartServer()

	tcpClient = NewTcpClient(func(client IClient) {
		Log.Debug("连接到服务器, 总计连接 = %v", tcpClient.GetClientCount())
	}, nil, parseProtocol, msgHandle, WithSendPackageSize(1024*1024))

	for {
		if tcpClient.GetClientCount() < 10 {
			if conn, err := net.DialTimeout("tcp", "127.0.0.1:6666", time.Second*5); err == nil {
				tcpClient.NewConnect(conn, nil)
			} else {
				Log.Error("连接到新节点失败, err = %v", err)
				//break
			}
		}
		time.Sleep(time.Second)
	}
	select {}
}

var ch chan []byte

func msgHandleOther(client IClient, data []byte) {
	Log.Warn("接受消息 = %v", string(data))
	// msgHandleOtherA(data[5:])
	ch <- data
}

func TestOther(t *testing.T) {
	tcpServer := NewTcpServer(6666, nil, nil, parseProtocol, msgHandleOther, WithRecvPackageSize(40))
	tcpServer.StartServer()
	tcpClient := NewTcpClient(nil, nil, nil, nil)

	ch = make(chan []byte, 50)

	if conn, err := net.DialTimeout("tcp", "127.0.0.1:6666", time.Second*5); err == nil {
		client := tcpClient.NewConnect(conn, nil)
		go func() {
			for i := 0; i < 100; i++ {
				msg := String.RandStr(25)
				Log.Debug("发送消息 = %v", msg)
				client.SendMsg(encodeMsg([]byte(msg)))
			}
		}()
		time.Sleep(time.Millisecond)
	} else {
		Log.Error("连接到新节点失败, err = %v", err)
	}
	time.Sleep(time.Second * 5)
	for {
		select {
		case d := <-ch:
			Log.Log("处理消息 = %v", string(d))
		}
	}
}
