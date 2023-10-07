package Net

import (
	"github.com/7058011439/haoqbb/Log"
	"github.com/7058011439/haoqbb/Stl"
	"github.com/7058011439/haoqbb/String"
	"github.com/7058011439/haoqbb/Timer"
	"github.com/7058011439/haoqbb/Util"
	"net"
	"sync"
	"testing"
	"time"
)

/*
1个链接 + 单包大小 50 字节 = 970W/s  500M/s
10个链接 + 单包大小 50 字节 = 1000W/s  520M/s
100个链接 + 单包大小 50 字节 = 1000W/s  520M/s
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

var mutex sync.Mutex
var totalCount uint64
var currCount uint64
var connectCount uint64
var dataLen uint64
var cost = Timer.NewTiming(Timer.Millisecond)

var data = []byte(String.RandStr(50))

func msgHandle(client IClient, data []byte) {
	defer mutex.Unlock()
	mutex.Lock()
	totalCount++
	currCount++
	dataLen += uint64(len(data))
	if cost.GetCost() > 1000 {
		Log.Debug("%v, currCount = %v, totalCount = %v, dataLen = %.3f m, tcpConn = %v", cost, currCount, totalCount, float64(dataLen)/1048576, connectCount)
		currCount = 0
		cost.ReStart()
	}
	//fmt.Println(string(data))
	//client.SendMsg(encodeMsg(data))
}

func TestClient_SendMsg(t *testing.T) {
	tcpServer := NewTcpServer(6666, nil, nil, parseProtocol, msgHandle)
	tcpServer.StartServer()

	tcpClient := NewTcpClient(func(client IClient) {
		connectCount++
	}, nil, parseProtocol, msgHandle, WithSendPackageSize(1024*1024))
	for i := 0; i < 10; i++ {
		if conn, err := net.DialTimeout("tcp", "127.0.0.1:6666", time.Second*5); err == nil {
			client := tcpClient.NewConnect(conn, nil)
			go func() {
				for {
					//time.Sleep(time.Nanosecond * 10)
					client.SendMsg(encodeMsg(data))
				}
			}()
			time.Sleep(time.Millisecond)
		} else {
			Log.ErrorLog("连接到新节点失败, err = %v", err)
			break
		}
	}
	select {}
}
